// Package identity is the gRPC transport for IdentityQueryService.
package identity

import (
	"context"
	"encoding/json"
	"time"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	identityread "github.com/medincident/medincident-backend/internal/service/query/identity"
	identityqueryv1 "github.com/medincident/medincident-backend/pkg/query/identity/v1"
)

// IdentityQueryHandler implements identityqueryv1.IdentityQueryServiceServer.
type IdentityQueryHandler struct {
	identityqueryv1.UnimplementedIdentityQueryServiceServer

	reader *identityread.Reader
}

// NewIdentityQueryHandler wires the handler with a Reader.
func NewIdentityQueryHandler(reader *identityread.Reader) *IdentityQueryHandler {
	return &IdentityQueryHandler{reader: reader}
}

// GetUser returns one user row.
func (h *IdentityQueryHandler) GetUser(
	ctx context.Context,
	req *identityqueryv1.GetUserRequest,
) (*identityqueryv1.GetUserResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	v, err := h.reader.GetUser(ctx, caller, req.GetId())
	if err != nil {
		return nil, err
	}
	return &identityqueryv1.GetUserResponse{User: userToProto(v)}, nil
}

// GetSession returns one session row.
func (h *IdentityQueryHandler) GetSession(
	ctx context.Context,
	req *identityqueryv1.GetSessionRequest,
) (*identityqueryv1.GetSessionResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}
	v, err := h.reader.GetSession(ctx, caller, req.GetId())
	if err != nil {
		return nil, err
	}
	return &identityqueryv1.GetSessionResponse{Session: sessionToProto(v)}, nil
}

// userAgentWire mirrors the on-disk JSON shape produced by the
// projector. Keeping a separate type for decoding keeps the reader
// surface free of json tags.
type userAgentWire struct {
	IP            string              `json:"ip"`
	Headers       map[string][]string `json:"headers,omitempty"`
	FingerprintID *string             `json:"fingerprint_id,omitempty"`
	Description   *string             `json:"description,omitempty"`
}

// userToProto adapts one UserView into the proto User message.
func userToProto(v *identityread.UserView) *identityqueryv1.User {
	return &identityqueryv1.User{
		Id:                v.ID,
		UserName:          v.UserName,
		FirstName:         v.FirstName,
		LastName:          v.LastName,
		DisplayName:       v.DisplayName,
		NickName:          v.NickName,
		Email:             v.Email,
		EmailVerified:     v.EmailVerified,
		PreferredLanguage: v.PreferredLanguage,
		Gender:            v.Gender,
		CreatedAt:         v.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:         v.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

// sessionToProto adapts one SessionView into the proto Session message.
func sessionToProto(v *identityread.SessionView) *identityqueryv1.Session {
	out := &identityqueryv1.Session{
		Id:                v.ID,
		UserId:            v.UserID,
		UserResourceOwner: v.UserResourceOwner,
		PreferredLanguage: v.PreferredLanguage,
		UserAgent:         userAgentToProto(v.UserAgent),
		CreatedAt:         v.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:         v.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
	if v.CheckedAt != nil {
		s := v.CheckedAt.UTC().Format(time.RFC3339Nano)
		out.CheckedAt = &s
	}
	return out
}

// userAgentToProto decodes the JSON blob and returns a proto
// UserAgent. On decode failure it returns an empty UserAgent so the
// caller still gets a sensible response.
func userAgentToProto(raw []byte) *identityqueryv1.UserAgent {
	out := &identityqueryv1.UserAgent{}
	if len(raw) == 0 {
		return out
	}
	var wire userAgentWire
	if err := json.Unmarshal(raw, &wire); err != nil {
		return out
	}
	out.Ip = wire.IP
	out.FingerprintId = wire.FingerprintID
	out.Description = wire.Description
	if len(wire.Headers) > 0 {
		out.Headers = make(map[string]*identityqueryv1.UserAgentHeaderValues, len(wire.Headers))
		for k, values := range wire.Headers {
			out.Headers[k] = &identityqueryv1.UserAgentHeaderValues{
				Values: append([]string(nil), values...),
			}
		}
	}
	return out
}
