package identity

import (
	"context"
	"strings"
	"time"

	"github.com/samber/oops"
	"gorm.io/gorm"

	usersv1 "github.com/medincident/medincident-zitadel-actions/gen/zitadel/users/v1"
)

// Error codes emitted by Projector methods.
const (
	ErrCodeUserAggregateIDEmpty = "user_aggregate_id_empty"
	ErrCodeUserProjectionFailed = "user_projection_failed"
)

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

	var sb strings.Builder
	sb.WriteString("UPDATE projections.users SET ")
	args := make([]any, 0, len(sets)+1)
	for i, s := range sets {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(s.col)
		sb.WriteString(" = ?")
		args = append(args, s.val)
	}
	sb.WriteString(" WHERE id = ?")
	args = append(args, userID)
	query := sb.String()

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
