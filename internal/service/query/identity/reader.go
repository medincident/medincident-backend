package identity

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/service/authz"
)

// Error codes emitted by Reader methods.
const (
	ErrCodeUserNotFound      = "user_not_found"
	ErrCodeUserLoadFailed    = "user_load_failed"
	ErrCodeSessionNotFound   = "session_not_found"
	ErrCodeSessionLoadFailed = "session_load_failed"
)

// UserView mirrors projections.users.
type UserView struct {
	ID                string
	UserName          string
	FirstName         string
	LastName          string
	DisplayName       string
	NickName          *string
	Email             string
	EmailVerified     bool
	PreferredLanguage string
	Gender            int32
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// SessionView mirrors projections.sessions. UserAgent is the raw JSONB
// bytes; callers decode to a structured type at handler/reader
// boundaries.
type SessionView struct {
	ID                string
	UserID            *string
	UserResourceOwner *string
	PreferredLanguage *string
	CheckedAt         *time.Time
	UserAgent         []byte
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// userSelectColumns lists every projections.users column read by
// GetUser and GetUserByEmail in the exact order scanUser expects.
const userSelectColumns = `id, user_name, first_name, last_name, display_name, nick_name,
		       email, email_verified, preferred_language, gender,
		       created_at, updated_at`

// scanUser populates v from row. Callers share the same SELECT list
// (see userSelectColumns) so both lookup variants scan in lockstep.
func scanUser(row interface{ Scan(...any) error }, v *UserView) error {
	return row.Scan(
		&v.ID, &v.UserName, &v.FirstName, &v.LastName, &v.DisplayName, &v.NickName,
		&v.Email, &v.EmailVerified, &v.PreferredLanguage, &v.Gender,
		&v.CreatedAt, &v.UpdatedAt,
	)
}

// GetUser returns the User projection row for the given Zitadel id.
// Authorization: the caller must be the target user themselves, or a
// system admin. The comparison runs in Go because both identifiers
// are opaque Zitadel strings, so no SQL round trip is needed to gate
// the "self" case.
func (r *Reader) GetUser(
	ctx context.Context,
	caller authz.Caller,
	userID string,
) (*UserView, error) {
	if caller.ZitadelUserID != userID {
		if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.SystemAdmin); err != nil {
			return nil, err
		}
	}
	var v UserView
	err := scanUser(r.db.WithContext(ctx).Raw(`
		SELECT `+userSelectColumns+`
		  FROM projections.users
		 WHERE id = ?`, userID,
	).Row(), &v)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In("reader.identity.user").
				Code(ErrCodeUserNotFound).
				Public("User not found.").
				With("user_id", userID).
				Errorf("not found")
		}
		return nil, oops.In("reader.identity.user").
			Code(ErrCodeUserLoadFailed).
			With("user_id", userID).
			Wrap(err)
	}
	return &v, nil
}

// GetUserByEmail returns the User projection row whose email matches
// the given address case-insensitively. Authorization: SystemAdmin
// only — a "find user by email" probe from a non-admin would leak
// membership of the identity projection, so no self fast-path exists.
//
// Zitadel enforces email uniqueness at the source, so only one row is
// expected. If multiple somehow come back, the query returns the
// first by insertion order; any duplicate implies a projection drift
// worth investigating rather than silently deduping.
func (r *Reader) GetUserByEmail(
	ctx context.Context,
	caller authz.Caller,
	email string,
) (*UserView, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.SystemAdmin); err != nil {
		return nil, err
	}
	var v UserView
	err := scanUser(r.db.WithContext(ctx).Raw(`
		SELECT `+userSelectColumns+`
		  FROM projections.users
		 WHERE LOWER(email) = LOWER(?)
		 LIMIT 1`, email,
	).Row(), &v)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In("reader.identity.user").
				Code(ErrCodeUserNotFound).
				Public("User not found.").
				With("email", email).
				Errorf("not found")
		}
		return nil, oops.In("reader.identity.user").
			Code(ErrCodeUserLoadFailed).
			With("email", email).
			Wrap(err)
	}
	return &v, nil
}

// GetSession returns the Session projection row for the given id.
// Authorization: AnyOf(SystemAdmin, SelfSession(id)) — a non-owner
// who is not a system admin receives permission_denied with no
// distinction from "session does not exist".
func (r *Reader) GetSession(
	ctx context.Context,
	caller authz.Caller,
	sessionID string,
) (*SessionView, error) {
	policy := authz.AnyOf(authz.SystemAdmin, authz.SelfSession(sessionID))
	if err := r.authz.Require(ctx, caller.ZitadelUserID, policy); err != nil {
		return nil, err
	}
	var v SessionView
	err := r.db.WithContext(ctx).Raw(`
		SELECT id, user_id, user_resource_owner, preferred_language, checked_at,
		       user_agent, created_at, updated_at
		  FROM projections.sessions
		 WHERE id = ?`, sessionID,
	).Row().Scan(
		&v.ID, &v.UserID, &v.UserResourceOwner, &v.PreferredLanguage, &v.CheckedAt,
		&v.UserAgent, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In("reader.identity.session").
				Code(ErrCodeSessionNotFound).
				Public("Session not found.").
				With("session_id", sessionID).
				Errorf("not found")
		}
		return nil, oops.In("reader.identity.session").
			Code(ErrCodeSessionLoadFailed).
			With("session_id", sessionID).
			Wrap(err)
	}
	return &v, nil
}
