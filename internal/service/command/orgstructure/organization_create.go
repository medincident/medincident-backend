package orgstructure

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
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

// CreateOrganizationCommand is the input of OrganizationService.Create.
type CreateOrganizationCommand struct {
	Name         string  `validate:"required,min=4,max=256"`
	Description  *string `validate:"omitnil,min=8,max=2048"`
	LegalAddress AddressInput
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
	if err := validation.Struct(cmd); err != nil {
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
		Name: strings.TrimSpace(cmd.Name),
		LegalAddress: model.Address{
			Text: strings.TrimSpace(cmd.LegalAddress.Text),
		},
	}
	if cmd.Description != nil {
		org.Description = null.StringFrom(strings.TrimSpace(*cmd.Description))
	}
	if cmd.LegalAddress.Point != nil {
		org.LegalAddress.Point = null.ValueFrom(model.Point{
			Longitude: cmd.LegalAddress.Point.Longitude,
			Latitude:  cmd.LegalAddress.Point.Latitude,
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
