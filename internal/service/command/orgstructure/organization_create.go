package orgstructure

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

// Error codes emitted by Organization-aggregate commands that are not
// primitive validation (infrastructure / existence / concurrency).
// Primitive validation codes (string_required, string_too_short,
// float_out_of_range, …) are produced by the validation translator.
const (
	ErrCodeOrganizationIDGenerationFailed = "organization_id_generation_failed"
	ErrCodeOrganizationSaveFailed         = "organization_save_failed"
	ErrCodeOrganizationLoadFailed         = "organization_load_failed"
	ErrCodeOrganizationNotFound           = "organization_not_found"
)

// CreateOrganizationPayload is the validated client-facing payload of
// CreateOrganization. All fields are primitives so the transport layer
// can hand raw proto values straight through.
type CreateOrganizationPayload struct {
	Name         string  `validate:"required,no_extra_ws,min=4,max=256"`
	Description  *string `validate:"omitnil,no_extra_ws,min=8,max=2048"`
	LegalAddress AddressInput
}

// CreateOrganizationCommand is the input of OrganizationService.Create.
// It bundles the authenticated caller (already validated at the
// transport boundary) with the request payload so the service layer is
// self-contained: validate → authorize → execute.
type CreateOrganizationCommand struct {
	Caller  authz.Caller
	Payload CreateOrganizationPayload
}

// CreateOrganizationResult is the output of OrganizationService.Create.
type CreateOrganizationResult struct {
	ID uuid.UUID
}

// Create persists a new Organization and writes the matching
// projection row in one transaction via the synchronous projector.
func (s *OrganizationService) Create(
	ctx context.Context,
	cmd CreateOrganizationCommand,
) (CreateOrganizationResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return CreateOrganizationResult{}, err
	}
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.SystemAdmin); err != nil {
		return CreateOrganizationResult{}, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return CreateOrganizationResult{}, oops.In("services.orgstructure.organization").
			Code(ErrCodeOrganizationIDGenerationFailed).
			Public("Failed to create organization.").
			Wrap(err)
	}

	org := model.Organization{
		ID:   id,
		Name: strings.TrimSpace(cmd.Payload.Name),
		LegalAddress: model.Address{
			Text: strings.TrimSpace(cmd.Payload.LegalAddress.Text),
		},
	}
	if cmd.Payload.Description != nil {
		org.Description = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
	}
	if cmd.Payload.LegalAddress.Point != nil {
		org.LegalAddress.Point = null.ValueFrom(model.Point{
			Longitude: cmd.Payload.LegalAddress.Point.Longitude,
			Latitude:  cmd.Payload.LegalAddress.Point.Latitude,
		})
	}

	var result CreateOrganizationResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&org).Error; err != nil {
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationSaveFailed).
				With("organization_id", id).
				Wrap(err)
		}

		if err := projector.OrganizationCreated(tx, &org); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
