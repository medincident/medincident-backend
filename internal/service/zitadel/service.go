// Package zitadel wraps the official zitadel-go v3 client and provides
// user-existence verification used by the membership service. Uses the
// SDK directly; no abstraction layer.
package zitadel

import (
	"context"
	"errors"
	"net"
	"net/url"
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/client"
	user_v2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/user/v2"
	zitadelcfg "github.com/zitadel/zitadel-go/v3/pkg/zitadel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

const (
	// ErrCodeZitadelClientBuildFailed is used when the zitadel-go client
	// cannot be constructed.
	ErrCodeZitadelClientBuildFailed = "zitadel_client_build_failed"
	// ErrCodeZitadelVerifyFailed is used when Verify encounters an
	// unexpected error from the Zitadel API.
	ErrCodeZitadelVerifyFailed = "zitadel_verify_failed"
)

// keepalive parameters for the persistent gRPC connection to Zitadel.
// Proxies and load balancers (nginx, AWS ALB, k8s ingress) typically close
// idle TCP connections after 60–600 s. Sending a ping every 30 s keeps the
// connection alive without triggering Zitadel's server-side GOAWAY (which
// enforces a minimum ping interval of 10 s by default).
const (
	zitadelKeepaliveTime    = 30 * time.Second
	zitadelKeepaliveTimeout = 10 * time.Second
)

// zitadelDialOpts returns gRPC dial options shared by all Zitadel client
// constructors. PermitWithoutStream must be true so pings are sent even
// when no RPC is in flight — exactly the idle-connection case we guard
// against.
func zitadelDialOpts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                zitadelKeepaliveTime,
			Timeout:             zitadelKeepaliveTimeout,
			PermitWithoutStream: true,
		}),
	}
}

// ErrUserNotFound is returned by Service.Verify and Service.GetUser when
// Zitadel reports NotFound for a given user ID. Any other error is wrapped
// with ErrCodeZitadelVerifyFailed.
var ErrUserNotFound = errors.New("zitadel: user not found")

// UserProfile holds the human-user profile fields fetched from Zitadel.
type UserProfile struct {
	UserName          string
	FirstName         string
	LastName          string
	DisplayName       string
	Email             string
	PreferredLanguage string
}

// Service wraps a configured zitadel-go v3 client. Concrete type —
// do NOT introduce an interface for it. Safe for concurrent use.
type Service struct {
	cl *client.Client
}

// NewServiceFromKeyFile builds a Service authenticated via JWT
// Profile using the service user key at keyPath. This is the
// production constructor. domain must be an https URL such as
// https://auth.example.com, with TLS.
//
// Does not perform an explicit health-check, but the underlying
// zitadel-go client.New implicitly hits the OIDC discovery endpoint
// while wiring the JWT Profile token source — callers must pass a
// bounded ctx so DI bootstrap cannot hang on an unreachable Zitadel.
func NewServiceFromKeyFile(ctx context.Context, logger *zerolog.Logger, domain, keyPath string) (*Service, error) {
	hostname, port, tls, tlsDefaulted, err := ParseDomain(domain)
	if err != nil {
		return nil, oops.In("services.zitadel").
			Code(ErrCodeZitadelClientBuildFailed).
			With("domain", domain).
			Wrap(err)
	}
	warnIfTLSDefaulted(logger, domain, tlsDefaulted)
	opts := ZitadelOptsFromParsed(port, tls)
	cl, err := client.New(
		ctx,
		zitadelcfg.New(hostname, opts...),
		client.WithAuth(client.DefaultServiceUserAuthentication(
			keyPath,
			oidc.ScopeOpenID,
			client.ScopeZitadelAPI(),
		)),
		client.WithGRPCDialOptions(zitadelDialOpts()...),
	)
	if err != nil {
		return nil, oops.In("services.zitadel").
			Code(ErrCodeZitadelClientBuildFailed).
			With("domain", domain).
			With("key_path", keyPath).
			Wrap(err)
	}
	return &Service{cl: cl}, nil
}

// NewServiceFromPAT builds a Service authenticated with a PAT. Used
// only by integration tests — production uses NewServiceFromKeyFile.
// domain may be an http(s) URL or a bare host:port string.
func NewServiceFromPAT(ctx context.Context, logger *zerolog.Logger, domain, pat string) (*Service, error) {
	hostname, port, tls, tlsDefaulted, err := ParseDomain(domain)
	if err != nil {
		return nil, oops.In("services.zitadel").
			Code(ErrCodeZitadelClientBuildFailed).
			With("domain", domain).
			Wrap(err)
	}
	warnIfTLSDefaulted(logger, domain, tlsDefaulted)
	opts := ZitadelOptsFromParsed(port, tls)
	cl, err := client.New(
		ctx,
		zitadelcfg.New(hostname, opts...),
		client.WithAuth(client.PAT(pat)),
		client.WithGRPCDialOptions(zitadelDialOpts()...),
	)
	if err != nil {
		return nil, oops.In("services.zitadel").
			Code(ErrCodeZitadelClientBuildFailed).
			With("domain", domain).
			Wrap(err)
	}
	return &Service{cl: cl}, nil
}

// ParseDomain parses a domain string that may be a full URL
// (http://localhost:8080) or a bare host:port (localhost:8080).
// Returns (hostname, port, isTLS, tlsDefaulted, error). Port is the
// string form (e.g. "8080"); "" means use the scheme default.
// tlsDefaulted is true only on the bare-hostname branch where we had
// no explicit scheme and no colon — see warnIfTLSDefaulted for the
// matching advisory log.
func ParseDomain(domain string) (hostname, port string, tls, tlsDefaulted bool, err error) {
	u, parseErr := url.Parse(domain)
	if parseErr != nil {
		return "", "", false, false, parseErr
	}

	switch u.Scheme {
	case "https":
		tls = true
		hostname = u.Hostname()
		port = u.Port() // "" if default 443
	case "http":
		tls = false
		hostname = u.Hostname()
		port = u.Port() // "" if default 80
	default:
		// No scheme — treat the whole string as host[:port].
		tls = false
		h, p, splitErr := net.SplitHostPort(domain)
		if splitErr != nil {
			// No colon at all — plain hostname, assume TLS on 443 and
			// surface the assumption via tlsDefaulted so the caller can
			// warn about it. Misconfigured dev setups that meant
			// http://host would fail loudly rather than silently hang.
			hostname = domain
			port = ""
			tls = true
			tlsDefaulted = true
		} else {
			hostname = h
			port = p
		}
	}
	return hostname, port, tls, tlsDefaulted, nil
}

// warnIfTLSDefaulted logs a warning whenever ParseDomain had to fall
// back to TLS-on-443 because the caller gave it a bare hostname with
// no scheme. Silent in the happy path; no-ops if the logger is nil.
func warnIfTLSDefaulted(logger *zerolog.Logger, domain string, tlsDefaulted bool) {
	if !tlsDefaulted || logger == nil {
		return
	}
	logger.Warn().
		Str("domain", domain).
		Msg("zitadel: no scheme in domain, assuming https://<host>:443 — add explicit http:// or https:// to silence this")
}

// ZitadelOptsFromParsed returns the zitadelcfg options that match the
// parsed hostname/port/tls combination.
func ZitadelOptsFromParsed(port string, tls bool) []zitadelcfg.Option {
	if !tls {
		// WithInsecure sets both port and disables TLS.
		p := port
		if p == "" {
			p = "80"
		}
		return []zitadelcfg.Option{zitadelcfg.WithInsecure(p)}
	}
	if port != "" && port != "443" {
		// TLS on a non-standard port.
		var n uint16
		for _, c := range port {
			n = n*10 + uint16(c-'0') //nolint:gosec // port string is already digits
		}
		return []zitadelcfg.Option{zitadelcfg.WithPort(n)}
	}
	return nil
}

// Verify returns nil if the user exists, ErrUserNotFound if Zitadel
// reports NotFound, or a wrapped error otherwise.
func (s *Service) Verify(ctx context.Context, zitadelUserID string) error {
	_, err := s.cl.UserServiceV2().GetUserByID(ctx, &user_v2.GetUserByIDRequest{
		UserId: zitadelUserID,
	})
	if err == nil {
		return nil
	}
	if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
		return ErrUserNotFound
	}
	return oops.In("services.zitadel").
		Code(ErrCodeZitadelVerifyFailed).
		With("zitadel_user_id", zitadelUserID).
		Wrap(err)
}

// GetUser fetches the human user profile from Zitadel. Returns
// ErrUserNotFound when Zitadel reports NotFound, or a wrapped error
// otherwise.
func (s *Service) GetUser(ctx context.Context, zitadelUserID string) (UserProfile, error) {
	resp, err := s.cl.UserServiceV2().GetUserByID(ctx, &user_v2.GetUserByIDRequest{
		UserId: zitadelUserID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return UserProfile{}, ErrUserNotFound
		}
		return UserProfile{}, oops.In("services.zitadel").
			Code(ErrCodeZitadelVerifyFailed).
			With("zitadel_user_id", zitadelUserID).
			Wrap(err)
	}
	u := resp.GetUser()
	p := UserProfile{UserName: u.GetUsername()}
	if h := u.GetHuman(); h != nil {
		if prof := h.GetProfile(); prof != nil {
			p.FirstName = prof.GetGivenName()
			p.LastName = prof.GetFamilyName()
			p.DisplayName = prof.GetDisplayName()
			p.PreferredLanguage = prof.GetPreferredLanguage()
		}
		if em := h.GetEmail(); em != nil {
			p.Email = em.GetEmail()
		}
	}
	return p, nil
}
