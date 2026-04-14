// Package zitadel wraps the official zitadel-go v3 client and provides
// user-existence verification used by the membership service. Uses the
// SDK directly; no abstraction layer.
package zitadel

import (
	"context"
	"errors"
	"net"
	"net/url"

	"github.com/samber/oops"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/client"
	user_v2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/user/v2"
	zitadelcfg "github.com/zitadel/zitadel-go/v3/pkg/zitadel"
	"google.golang.org/grpc/codes"
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

// ErrUserNotFound is returned by Service.Verify when Zitadel reports
// NotFound for a given user ID. Any other error is wrapped with
// ErrCodeZitadelVerifyFailed.
var ErrUserNotFound = errors.New("zitadel: user not found")

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
// Does NOT ping Zitadel — per project convention DI factories don't
// health-check.
func NewServiceFromKeyFile(ctx context.Context, domain, keyPath string) (*Service, error) {
	hostname, port, tls, err := parseDomain(domain)
	if err != nil {
		return nil, oops.In("zitadel").
			Code(ErrCodeZitadelClientBuildFailed).
			With("domain", domain).
			Wrap(err)
	}
	opts := zitadelOptsFromParsed(hostname, port, tls)
	cl, err := client.New(
		ctx,
		zitadelcfg.New(hostname, opts...),
		client.WithAuth(client.DefaultServiceUserAuthentication(
			keyPath,
			oidc.ScopeOpenID,
			client.ScopeZitadelAPI(),
		)),
	)
	if err != nil {
		return nil, oops.In("zitadel").
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
func NewServiceFromPAT(ctx context.Context, domain, pat string) (*Service, error) {
	hostname, port, tls, err := parseDomain(domain)
	if err != nil {
		return nil, oops.In("zitadel").
			Code(ErrCodeZitadelClientBuildFailed).
			With("domain", domain).
			Wrap(err)
	}
	opts := zitadelOptsFromParsed(hostname, port, tls)
	cl, err := client.New(
		ctx,
		zitadelcfg.New(hostname, opts...),
		client.WithAuth(client.PAT(pat)),
	)
	if err != nil {
		return nil, oops.In("zitadel").
			Code(ErrCodeZitadelClientBuildFailed).
			With("domain", domain).
			Wrap(err)
	}
	return &Service{cl: cl}, nil
}

// parseDomain parses a domain string that may be a full URL
// (http://localhost:8080) or a bare host:port (localhost:8080).
// Returns (hostname, port, isTLS, error). Port is the string form
// (e.g. "8080"); "" means use the scheme default.
func parseDomain(domain string) (hostname, port string, tls bool, err error) {
	u, parseErr := url.Parse(domain)
	if parseErr != nil {
		return "", "", false, parseErr
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
			// No colon at all — plain hostname, use 443 TLS default.
			hostname = domain
			port = ""
			tls = true
		} else {
			hostname = h
			port = p
		}
	}
	return hostname, port, tls, nil
}

// zitadelOptsFromParsed returns the zitadelcfg options that match the
// parsed hostname/port/tls combination.
func zitadelOptsFromParsed(_, port string, tls bool) []zitadelcfg.Option {
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
	return oops.In("zitadel").
		Code(ErrCodeZitadelVerifyFailed).
		With("zitadel_user_id", zitadelUserID).
		Wrap(err)
}
