package identity

import (
	"context"
	"encoding/json"
	"time"

	"github.com/samber/oops"
	"gorm.io/gorm"

	sessionsv1 "github.com/medincident/medincident-zitadel-actions/gen/zitadel/sessions/v1"
	usersv1 "github.com/medincident/medincident-zitadel-actions/gen/zitadel/users/v1"
)

// Error codes emitted by Projector methods.
const (
	ErrCodeUserAggregateIDEmpty      = "user_aggregate_id_empty"
	ErrCodeUserProjectionFailed      = "user_projection_failed"
	ErrCodeSessionAggregateIDEmpty   = "session_aggregate_id_empty"
	ErrCodeSessionUserAgentEncodeErr = "session_user_agent_encode_failed"
	ErrCodeSessionProjectionFailed   = "session_projection_failed"
)

// emptyUserAgentJSON is the canonical empty object stored when a
// session arrives without a user_agent payload.
var emptyUserAgentJSON = []byte(`{}`)

// userAgentJSON is the on-disk shape of a zitadel.sessions.v1.UserAgent
// value. We transcode the proto message into a plain map before writing
// so the storage format is independent of protojson quirks and
// survives future proto additions without a DB migration.
type userAgentJSON struct {
	IP            string              `json:"ip"`
	Headers       map[string][]string `json:"headers,omitempty"`
	FingerprintID *string             `json:"fingerprint_id,omitempty"`
	Description   *string             `json:"description,omitempty"`
}

// ApplyUserHumanAdded upserts a user projection and back-fills any
// employee_cards already created for this Zitadel user id.
func (p *Projector) ApplyUserHumanAdded(
	ctx context.Context,
	userID string,
	occurredAt time.Time,
	event *usersv1.UserHumanAdded,
) error {
	if userID == "" {
		return oops.In("projector.identity.user").
			Code(ErrCodeUserAggregateIDEmpty).
			Errorf("user aggregate id empty")
	}

	var nickName *string
	if event.NickName != nil {
		v := event.GetNickName()
		nickName = &v
	}

	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
			INSERT INTO projections.users
			    (id, user_name, first_name, last_name, display_name, nick_name,
			     email, email_verified, preferred_language, gender,
			     created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, FALSE, ?, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING`,
			userID, event.GetUserName(), event.GetFirstName(), event.GetLastName(),
			event.GetDisplayName(), nickName,
			event.GetEmail(), event.GetPreferredLanguage(), int32(event.GetGender()),
			occurredAt.UTC(), occurredAt.UTC(),
		).Error; err != nil {
			return oops.In("projector.identity.user").
				Code(ErrCodeUserProjectionFailed).
				With("user_id", userID).
				Wrap(err)
		}
		if err := tx.Exec(`
			UPDATE projections.employee_cards
			   SET first_name = ?, last_name = ?, display_name = ?, email = ?,
			       updated_at = now()
			 WHERE zitadel_user_id = ?`,
			event.GetFirstName(), event.GetLastName(), event.GetDisplayName(), event.GetEmail(),
			userID,
		).Error; err != nil {
			return oops.In("projector.identity.user").
				Code(ErrCodeUserProjectionFailed).
				With("user_id", userID).
				Wrap(err)
		}
		return nil
	})
}

// ApplyUserHumanProfileChanged applies a partial update driven by the
// FieldMask carried in the event. Paths not listed in the mask are
// left untouched. Missing user row is a soft-miss (out-of-order).
func (p *Projector) ApplyUserHumanProfileChanged(
	ctx context.Context,
	userID string,
	occurredAt time.Time,
	event *usersv1.UserHumanProfileChanged,
) error {
	if userID == "" {
		return oops.In("projector.identity.user").
			Code(ErrCodeUserAggregateIDEmpty).
			Errorf("user aggregate id empty")
	}
	mask := event.GetUpdatedFields()
	if mask == nil || len(mask.GetPaths()) == 0 {
		return nil
	}

	type setExpr struct {
		col string
		val any
	}
	sets := []setExpr{{col: "updated_at", val: occurredAt.UTC()}}
	touchesCard := false
	for _, path := range mask.GetPaths() {
		switch path {
		case "first_name":
			sets = append(sets, setExpr{col: "first_name", val: event.GetFirstName()})
			touchesCard = true
		case "last_name":
			sets = append(sets, setExpr{col: "last_name", val: event.GetLastName()})
			touchesCard = true
		case "display_name":
			sets = append(sets, setExpr{col: "display_name", val: event.GetDisplayName()})
			touchesCard = true
		case "preferred_language":
			sets = append(sets, setExpr{col: "preferred_language", val: event.GetPreferredLanguage()})
		case "gender":
			sets = append(sets, setExpr{col: "gender", val: int32(event.GetGender())})
		case "nick_name":
			var v *string
			if event.NickName != nil {
				s := event.GetNickName()
				v = &s
			}
			sets = append(sets, setExpr{col: "nick_name", val: v})
		}
	}
	if len(sets) == 1 {
		// Only updated_at would change — skip.
		return nil
	}

	query := "UPDATE projections.users SET "
	args := make([]any, 0, len(sets)+1)
	for i, s := range sets {
		if i > 0 {
			query += ", "
		}
		query += s.col + " = ?"
		args = append(args, s.val)
	}
	query += " WHERE id = ?"
	args = append(args, userID)

	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Exec(query, args...)
		if res.Error != nil {
			return oops.In("projector.identity.user").
				Code(ErrCodeUserProjectionFailed).
				With("user_id", userID).
				Wrap(res.Error)
		}
		if res.RowsAffected == 0 {
			p.logger.Debug().Str("user_id", userID).Msg("user profile_changed before added; skipped")
			return nil
		}
		if !touchesCard {
			return nil
		}
		// Re-read authoritative name set from the row and propagate to
		// employee_cards so unchanged mask paths stay preserved.
		var first, last, display, email string
		if err := tx.Raw(
			`SELECT first_name, last_name, display_name, email
			   FROM projections.users WHERE id = ?`, userID,
		).Row().Scan(&first, &last, &display, &email); err != nil {
			return oops.In("projector.identity.user").
				Code(ErrCodeUserProjectionFailed).
				With("user_id", userID).
				Wrap(err)
		}
		if err := tx.Exec(`
			UPDATE projections.employee_cards
			   SET first_name = ?, last_name = ?, display_name = ?, email = ?,
			       updated_at = now()
			 WHERE zitadel_user_id = ?`,
			first, last, display, email, userID,
		).Error; err != nil {
			return oops.In("projector.identity.user").
				Code(ErrCodeUserProjectionFailed).
				With("user_id", userID).
				Wrap(err)
		}
		return nil
	})
}

// ApplyUserHumanEmailChanged rewrites the email column and resets
// email_verified to false (matching Zitadel's own semantics).
func (p *Projector) ApplyUserHumanEmailChanged(
	ctx context.Context,
	userID string,
	occurredAt time.Time,
	event *usersv1.UserHumanEmailChanged,
) error {
	if userID == "" {
		return oops.In("projector.identity.user").
			Code(ErrCodeUserAggregateIDEmpty).
			Errorf("user aggregate id empty")
	}
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Exec(`
			UPDATE projections.users
			   SET email = ?, email_verified = FALSE, updated_at = ?
			 WHERE id = ?`, event.GetEmail(), occurredAt.UTC(), userID)
		if res.Error != nil {
			return oops.In("projector.identity.user").
				Code(ErrCodeUserProjectionFailed).
				With("user_id", userID).
				Wrap(res.Error)
		}
		if res.RowsAffected == 0 {
			p.logger.Debug().Str("user_id", userID).Msg("user email_changed before added; skipped")
			return nil
		}
		if err := tx.Exec(`
			UPDATE projections.employee_cards
			   SET email = ?, updated_at = now()
			 WHERE zitadel_user_id = ?`, event.GetEmail(), userID,
		).Error; err != nil {
			return oops.In("projector.identity.user").
				Code(ErrCodeUserProjectionFailed).
				With("user_id", userID).
				Wrap(err)
		}
		return nil
	})
}

// ApplyUserHumanEmailVerified flips email_verified to true. Marker
// event with no payload fields consulted.
func (p *Projector) ApplyUserHumanEmailVerified(
	ctx context.Context,
	userID string,
	occurredAt time.Time,
) error {
	if userID == "" {
		return oops.In("projector.identity.user").
			Code(ErrCodeUserAggregateIDEmpty).
			Errorf("user aggregate id empty")
	}
	res := p.db.WithContext(ctx).Exec(`
		UPDATE projections.users
		   SET email_verified = TRUE, updated_at = ?
		 WHERE id = ?`, occurredAt.UTC(), userID)
	if res.Error != nil {
		return oops.In("projector.identity.user").
			Code(ErrCodeUserProjectionFailed).
			With("user_id", userID).
			Wrap(res.Error)
	}
	if res.RowsAffected == 0 {
		p.logger.Debug().Str("user_id", userID).Msg("user email_verified before added; skipped")
	}
	return nil
}

// ApplySessionAdded upserts a fresh session row.
func (p *Projector) ApplySessionAdded(
	ctx context.Context,
	sessionID string,
	occurredAt time.Time,
	event *sessionsv1.SessionAdded,
) error {
	if sessionID == "" {
		return oops.In("projector.identity.session").
			Code(ErrCodeSessionAggregateIDEmpty).
			Errorf("session aggregate id empty")
	}
	agentBytes, err := encodeUserAgent(event.GetUserAgent())
	if err != nil {
		return oops.In("projector.identity.session").
			Code(ErrCodeSessionUserAgentEncodeErr).
			With("session_id", sessionID).
			Wrap(err)
	}
	if err := p.db.WithContext(ctx).Exec(`
		INSERT INTO projections.sessions
		    (id, user_agent, created_at, updated_at)
		VALUES (?, ?::jsonb, ?, ?)
		ON CONFLICT (id) DO NOTHING`,
		sessionID, string(agentBytes), occurredAt.UTC(), occurredAt.UTC(),
	).Error; err != nil {
		return oops.In("projector.identity.session").
			Code(ErrCodeSessionProjectionFailed).
			With("session_id", sessionID).
			Wrap(err)
	}
	return nil
}

// ApplySessionUserChecked updates user_id / user_resource_owner /
// preferred_language / checked_at on an existing session row. No-op
// if the row is absent.
func (p *Projector) ApplySessionUserChecked(
	ctx context.Context,
	sessionID string,
	occurredAt time.Time,
	event *sessionsv1.SessionUserChecked,
) error {
	if sessionID == "" {
		return oops.In("projector.identity.session").
			Code(ErrCodeSessionAggregateIDEmpty).
			Errorf("session aggregate id empty")
	}
	var preferred *string
	if pl := event.GetPreferredLanguage(); pl != "" {
		v := pl
		preferred = &v
	}
	var checkedAt *time.Time
	if ts := event.GetCheckedAt(); ts != nil && ts.IsValid() {
		v := ts.AsTime().UTC()
		checkedAt = &v
	}
	res := p.db.WithContext(ctx).Exec(`
		UPDATE projections.sessions
		   SET user_id = ?, user_resource_owner = ?, preferred_language = ?, checked_at = ?,
		       updated_at = ?
		 WHERE id = ?`,
		event.GetUserId(), event.GetUserResourceOwner(), preferred, checkedAt,
		occurredAt.UTC(), sessionID,
	)
	if res.Error != nil {
		return oops.In("projector.identity.session").
			Code(ErrCodeSessionProjectionFailed).
			With("session_id", sessionID).
			Wrap(res.Error)
	}
	if res.RowsAffected == 0 {
		// SessionAdded hasn't been projected yet. Unlike user_human_added
		// (which upserts and therefore retro-fills on late arrival),
		// SessionAdded is a straight UPDATE — if we lose this
		// SessionUserChecked payload the session row never gets its
		// user_id/checked_at columns. Log at WARN so this is visible in
		// prod without forcing a retry loop here.
		p.logger.Warn().
			Str("session_id", sessionID).
			Str("user_id", event.GetUserId()).
			Msg("session user_checked arrived before session_added; fields will be missing until a follow-up event")
	}
	return nil
}

// encodeUserAgent marshals a proto UserAgent into the on-disk JSON
// format used by projections.sessions.user_agent.
func encodeUserAgent(ua *sessionsv1.UserAgent) ([]byte, error) {
	if ua == nil {
		return emptyUserAgentJSON, nil
	}
	wire := userAgentJSON{IP: ua.GetIp()}
	if fid := ua.GetFingerprintId(); fid != "" {
		v := fid
		wire.FingerprintID = &v
	}
	if d := ua.GetDescription(); d != "" {
		v := d
		wire.Description = &v
	}
	if headers := ua.GetHeaders(); len(headers) > 0 {
		wire.Headers = make(map[string][]string, len(headers))
		for name, values := range headers {
			wire.Headers[name] = append([]string(nil), values.GetValues()...)
		}
	}
	return json.Marshal(wire)
}
