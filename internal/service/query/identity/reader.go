package identity

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/samber/oops"
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

// GetUser returns the User projection row for the given Zitadel id.
func (r *Reader) GetUser(ctx context.Context, userID string) (*UserView, error) {
	var v UserView
	err := r.db.WithContext(ctx).Raw(`
		SELECT id, user_name, first_name, last_name, display_name, nick_name,
		       email, email_verified, preferred_language, gender,
		       created_at, updated_at
		  FROM projections.users
		 WHERE id = ?`, userID,
	).Row().Scan(
		&v.ID, &v.UserName, &v.FirstName, &v.LastName, &v.DisplayName, &v.NickName,
		&v.Email, &v.EmailVerified, &v.PreferredLanguage, &v.Gender,
		&v.CreatedAt, &v.UpdatedAt,
	)
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

// GetSession returns the Session projection row for the given id.
func (r *Reader) GetSession(ctx context.Context, sessionID string) (*SessionView, error) {
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
