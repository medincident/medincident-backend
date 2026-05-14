package orgstructure

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// ActivateOrganizationPayload identifies the organization to activate.
type ActivateOrganizationPayload struct {
	ID string `validate:"required,uuid"`
}

// ActivateOrganizationCommand = caller + payload.
type ActivateOrganizationCommand struct {
	Caller  authz.Caller
	Payload ActivateOrganizationPayload
}

// Activate marks a deactivated organization as active. No cascade —
// child entities keep their individual is_active state.
//
// See: docs/services/OrgStructure.md
func (s *OrganizationService) Activate(
	ctx context.Context,
	cmd ActivateOrganizationCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	orgID := uuid.MustParse(cmd.Payload.ID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Organization(orgID)); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var org model.Organization
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&org, "id = ?", orgID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.organization").
					Code(ErrCodeOrganizationNotFound).
					Public("Organization not found.").
					With("organization_id", orgID).
					Wrap(err)
			}
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationLoadFailed).
				With("organization_id", orgID).
				Wrap(err)
		}

		if org.IsActive {
			return nil
		}

		var updatedAt time.Time
		if err := tx.Raw(
			`UPDATE domain.organizations SET is_active = TRUE, updated_at = now() WHERE id = ? RETURNING updated_at`,
			orgID,
		).Row().Scan(&updatedAt); err != nil {
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationSaveFailed).
				With("organization_id", orgID).
				Wrap(err)
		}
		return appendOrganizationActivatedEvent(tx, orgID, updatedAt)
	})
}

func appendOrganizationActivatedEvent(tx *gorm.DB, orgID uuid.UUID, now time.Time) error {
	msg := &orgv1.OrganizationActivated{OrganizationId: orgID.String(), UpdatedAt: timestamppb.New(now)}
	payload, err := anypb.New(msg)
	if err != nil {
		return oops.In("services.orgstructure.organization").Code(ErrCodeOrganizationSaveFailed).Wrap(err)
	}
	env := &eventv1.Envelope{
		OccurredAt:    timestamppb.New(now),
		AggregateType: "organization",
		AggregateId:   orgID.String(),
		Payload:       payload,
	}
	return outbox.Append(tx, "medincident.event.organization.v1.activated", env)
}
