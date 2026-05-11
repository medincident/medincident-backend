package orgstructure

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	anypb "google.golang.org/protobuf/types/known/anypb"
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

// UpdateOrganizationLegalAddressPayload carries the new legal address.
type UpdateOrganizationLegalAddressPayload struct {
	ID      string       `validate:"required,uuid"`
	Address AddressInput `validate:"required"`
}

// UpdateOrganizationLegalAddressCommand = caller + payload.
type UpdateOrganizationLegalAddressCommand struct {
	Caller  authz.Caller
	Payload UpdateOrganizationLegalAddressPayload
}

func buildOrganizationLegalAddressChangedEnvelope(org *model.Organization) (*eventv1.Envelope, error) {
	msg := &orgv1.OrganizationLegalAddressChanged{
		LegalAddress: buildOrgAddressProto(org.LegalAddress),
		UpdatedAt:    timestamppb.New(org.UpdatedAt),
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In("services.orgstructure.organization").
			Code(ErrCodeOrganizationSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(org.UpdatedAt),
		AggregateType: "organization",
		AggregateId:   org.ID.String(),
		Payload:       payload,
	}, nil
}

// UpdateLegalAddress replaces the organization's legal address. Returns
// nil without writing when the new address is equal to the stored one.
//
// See: docs/services/OrgStructure.md
func (s *OrganizationService) UpdateLegalAddress(
	ctx context.Context,
	cmd UpdateOrganizationLegalAddressCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.ID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Organization(id)); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var org model.Organization
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&org, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.organization").
					Code(ErrCodeOrganizationNotFound).
					Public("Organization not found.").
					With("organization_id", id).
					Errorf("organization not found")
			}
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationLoadFailed).
				With("organization_id", id).
				Wrap(err)
		}

		newAddress := model.Address{
			Text: strings.TrimSpace(cmd.Payload.Address.Text),
		}
		if cmd.Payload.Address.Point != nil {
			newAddress.Point = null.ValueFrom(model.Point{
				Longitude: cmd.Payload.Address.Point.Longitude,
				Latitude:  cmd.Payload.Address.Point.Latitude,
			})
		}
		if org.LegalAddress.Equal(newAddress) {
			return nil
		}
		org.LegalAddress = newAddress
		if err := tx.Save(&org).Error; err != nil {
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationSaveFailed).
				With("organization_id", id).
				Wrap(err)
		}

		env, err := buildOrganizationLegalAddressChangedEnvelope(&org)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.organization.v1.legal_address_changed", env)
	})
}
