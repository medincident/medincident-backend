package orgstructure

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
	"github.com/medincident/medincident-command-service/internal/service/validation"
)

// Error codes emitted by Clinic-aggregate commands that are not
// primitive validation. Validation codes are generic and live in
// internal/validation.
const (
	ErrCodeClinicIDGenerationFailed   = "clinic_id_generation_failed"
	ErrCodeClinicSaveFailed           = "clinic_save_failed"
	ErrCodeClinicLoadFailed           = "clinic_load_failed"
	ErrCodeClinicNotFound             = "clinic_not_found"
	ErrCodeClinicOrganizationNotFound = "clinic_organization_not_found"
)

// CreateClinicCommand is the input of ClinicService.Create.
type CreateClinicCommand struct {
	OrganizationID  uuid.UUID `validate:"required"`
	Name            string    `validate:"required,min=4,max=256"`
	Description     *string   `validate:"omitnil,min=8,max=2048"`
	PhysicalAddress AddressInput
}

// CreateClinicResult is the output of ClinicService.Create.
type CreateClinicResult struct {
	ID uuid.UUID
}

// Create persists a new Clinic under the given organization.
func (s *ClinicService) Create(
	ctx context.Context,
	cmd CreateClinicCommand,
) (CreateClinicResult, error) {
	if err := validation.Struct(cmd); err != nil {
		return CreateClinicResult{}, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return CreateClinicResult{}, oops.In("services.orgstructure.clinic").
			Code(ErrCodeClinicIDGenerationFailed).
			Public("Failed to create clinic.").
			Wrap(err)
	}

	clinic := model.Clinic{
		ID:             id,
		OrganizationID: cmd.OrganizationID,
		Name:           strings.TrimSpace(cmd.Name),
		PhysicalAddress: model.Address{
			Text: strings.TrimSpace(cmd.PhysicalAddress.Text),
		},
	}
	if cmd.Description != nil {
		clinic.Description = null.StringFrom(strings.TrimSpace(*cmd.Description))
	}
	if cmd.PhysicalAddress.Point != nil {
		clinic.PhysicalAddress.Point = null.ValueFrom(model.Point{
			Longitude: cmd.PhysicalAddress.Point.Longitude,
			Latitude:  cmd.PhysicalAddress.Point.Latitude,
		})
	}

	var result CreateClinicResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&clinic).Error; err != nil {
			if errors.Is(err, gorm.ErrForeignKeyViolated) {
				return oops.In("services.orgstructure.clinic").
					Code(ErrCodeClinicOrganizationNotFound).
					Public("Organization not found.").
					With("organization_id", cmd.OrganizationID).
					Wrap(err)
			}
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicSaveFailed).
				With("clinic_id", id).
				Wrap(err)
		}

		if err := projector.ClinicCreated(tx, &clinic); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
