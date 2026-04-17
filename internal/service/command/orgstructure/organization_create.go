package orgstructure

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/service/command/projector"
)

// Organization name and description invariant limits.
const (
	organizationMinNameLen = 4
	organizationMaxNameLen = 256
	organizationMinDescLen = 8
	organizationMaxDescLen = 2048
)

// Error codes used by every Organization service method.
const (
	ErrCodeOrganizationNameEmpty           = "organization_name_empty"
	ErrCodeOrganizationNameTooShort        = "organization_name_too_short"
	ErrCodeOrganizationNameTooLong         = "organization_name_too_long"
	ErrCodeOrganizationDescriptionTooShort = "organization_description_too_short"
	ErrCodeOrganizationDescriptionTooLong  = "organization_description_too_long"

	ErrCodeOrganizationIDGenerationFailed = "organization_id_generation_failed"
	ErrCodeOrganizationSaveFailed         = "organization_save_failed"
	ErrCodeOrganizationLoadFailed         = "organization_load_failed"
	ErrCodeOrganizationNotFound           = "organization_not_found"
)

// CreateOrganizationCommand is the input of OrganizationService.Create.
type CreateOrganizationCommand struct {
	Name         string
	Description  *string // nil = no description
	LegalAddress AddressInput
}

// CreateOrganizationResult is the output of OrganizationService.Create.
type CreateOrganizationResult struct {
	ID uuid.UUID
}

// validateOrganizationName ensures the name is present and within range.
func validateOrganizationName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return oops.In("services.orgstructure.organization").
			Code(ErrCodeOrganizationNameEmpty).
			Public("Organization name is required.").
			With("field", "name").
			Errorf("name is empty")
	}
	n := utf8.RuneCountInString(trimmed)
	if n < organizationMinNameLen {
		return oops.In("services.orgstructure.organization").
			Code(ErrCodeOrganizationNameTooShort).
			Public("Organization name is too short.").
			With("field", "name").
			With("actual_length", n).
			With("min_length", organizationMinNameLen).
			Errorf("name too short")
	}
	if n > organizationMaxNameLen {
		return oops.In("services.orgstructure.organization").
			Code(ErrCodeOrganizationNameTooLong).
			Public("Organization name is too long.").
			With("field", "name").
			With("actual_length", n).
			With("max_length", organizationMaxNameLen).
			Errorf("name too long")
	}
	return nil
}

// validateOrganizationDescription is nil-tolerant: a nil pointer means
// "no description", which is always valid. When non-nil, the trimmed
// value must satisfy length bounds.
func validateOrganizationDescription(desc *string) error {
	if desc == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*desc)
	n := utf8.RuneCountInString(trimmed)
	if n < organizationMinDescLen {
		return oops.In("services.orgstructure.organization").
			Code(ErrCodeOrganizationDescriptionTooShort).
			Public("Organization description is too short.").
			With("field", "description").
			With("actual_length", n).
			With("min_length", organizationMinDescLen).
			Errorf("description too short")
	}
	if n > organizationMaxDescLen {
		return oops.In("services.orgstructure.organization").
			Code(ErrCodeOrganizationDescriptionTooLong).
			Public("Organization description is too long.").
			With("field", "description").
			With("actual_length", n).
			With("max_length", organizationMaxDescLen).
			Errorf("description too long")
	}
	return nil
}

// Create persists a new Organization and writes the matching
// projection row in one transaction via the synchronous projector.
func (s *OrganizationService) Create(
	ctx context.Context,
	cmd CreateOrganizationCommand,
) (CreateOrganizationResult, error) {
	var errs []error
	if err := validateOrganizationName(cmd.Name); err != nil {
		errs = append(errs, err)
	}
	if err := validateOrganizationDescription(cmd.Description); err != nil {
		errs = append(errs, err)
	}
	if err := validateAddressInput(cmd.LegalAddress); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return CreateOrganizationResult{}, errors.Join(errs...)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return CreateOrganizationResult{}, oops.In("services.orgstructure.organization").
			Code(ErrCodeOrganizationIDGenerationFailed).
			Public("Failed to create organization.").
			Wrap(err)
	}

	org := model.Organization{
		ID:          id,
		Name:        strings.TrimSpace(cmd.Name),
		Description: null.StringFromPtr(cmd.Description),
		LegalAddress: model.Address{
			Text: strings.TrimSpace(cmd.LegalAddress.Text),
		},
	}
	if cmd.LegalAddress.Point != nil {
		org.LegalAddress.Point = &model.Point{
			Longitude: cmd.LegalAddress.Point.Longitude,
			Latitude:  cmd.LegalAddress.Point.Latitude,
		}
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
