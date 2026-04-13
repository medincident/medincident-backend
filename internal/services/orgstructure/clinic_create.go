package orgstructure

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	clinicv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/clinic/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
)

const (
	clinicMinNameLen = 4
	clinicMaxNameLen = 256
	clinicMinDescLen = 8
	clinicMaxDescLen = 2048
)

const (
	ErrCodeClinicNameEmpty           = "clinic_name_empty"
	ErrCodeClinicNameTooShort        = "clinic_name_too_short"
	ErrCodeClinicNameTooLong         = "clinic_name_too_long"
	ErrCodeClinicDescriptionTooShort = "clinic_description_too_short"
	ErrCodeClinicDescriptionTooLong  = "clinic_description_too_long"

	ErrCodeClinicIDGenerationFailed   = "clinic_id_generation_failed"
	ErrCodeClinicSaveFailed           = "clinic_save_failed"
	ErrCodeClinicLoadFailed           = "clinic_load_failed"
	ErrCodeClinicNotFound             = "clinic_not_found"
	ErrCodeClinicEventBuildFailed     = "clinic_event_build_failed"
	ErrCodeClinicOrganizationNotFound = "clinic_organization_not_found"
)

const (
	SubjectClinicCreated                = "medincident.event.clinic.v1.created"
	SubjectClinicDetailsChanged         = "medincident.event.clinic.v1.details_changed"
	SubjectClinicPhysicalAddressChanged = "medincident.event.clinic.v1.physical_address_changed"
)

// pgErrCodeForeignKeyViolation is Postgres error code 23503.
const pgErrCodeForeignKeyViolation = "23503"

type CreateClinicCommand struct {
	OrganizationID  uuid.UUID
	Name            string
	Description     *string
	PhysicalAddress AddressInput
}

type CreateClinicResult struct {
	ID uuid.UUID
}

func validateClinicName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return oops.In("services.orgstructure.clinic").
			Code(ErrCodeClinicNameEmpty).
			Public("Clinic name is required.").
			With("field", "name").
			Errorf("name is empty")
	}
	n := utf8.RuneCountInString(trimmed)
	if n < clinicMinNameLen {
		return oops.In("services.orgstructure.clinic").
			Code(ErrCodeClinicNameTooShort).
			Public("Clinic name is too short.").
			With("field", "name").
			With("actual_length", n).
			With("min_length", clinicMinNameLen).
			Errorf("name too short")
	}
	if n > clinicMaxNameLen {
		return oops.In("services.orgstructure.clinic").
			Code(ErrCodeClinicNameTooLong).
			Public("Clinic name is too long.").
			With("field", "name").
			With("actual_length", n).
			With("max_length", clinicMaxNameLen).
			Errorf("name too long")
	}
	return nil
}

func validateClinicDescription(desc *string) error {
	if desc == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*desc)
	n := utf8.RuneCountInString(trimmed)
	if n < clinicMinDescLen {
		return oops.In("services.orgstructure.clinic").
			Code(ErrCodeClinicDescriptionTooShort).
			Public("Clinic description is too short.").
			With("field", "description").
			With("actual_length", n).
			With("min_length", clinicMinDescLen).
			Errorf("description too short")
	}
	if n > clinicMaxDescLen {
		return oops.In("services.orgstructure.clinic").
			Code(ErrCodeClinicDescriptionTooLong).
			Public("Clinic description is too long.").
			With("field", "description").
			With("actual_length", n).
			With("max_length", clinicMaxDescLen).
			Errorf("description too long")
	}
	return nil
}

func buildClinicCreatedEvent(c *model.Clinic) *clinicv1.ClinicCreated {
	ev := &clinicv1.ClinicCreated{
		OrganizationId:  c.OrganizationID.String(),
		Name:            c.Name,
		PhysicalAddress: &clinicv1.Address{Text: c.PhysicalAddress.Text},
	}
	if c.Description.Valid {
		desc := c.Description.String
		ev.Description = &desc
	}
	if c.PhysicalAddress.Point.Longitude.Valid && c.PhysicalAddress.Point.Latitude.Valid {
		ev.PhysicalAddress.Point = &clinicv1.Point{
			Longitude: c.PhysicalAddress.Point.Longitude.Float64,
			Latitude:  c.PhysicalAddress.Point.Latitude.Float64,
		}
	}
	return ev
}

func (s *ClinicService) Create(
	ctx context.Context,
	cmd CreateClinicCommand,
) (CreateClinicResult, error) {
	var errs []error
	if err := validateClinicName(cmd.Name); err != nil {
		errs = append(errs, err)
	}
	if err := validateClinicDescription(cmd.Description); err != nil {
		errs = append(errs, err)
	}
	if err := validateAddressInput(cmd.PhysicalAddress); err != nil {
		errs = append(errs, err)
	}
	if err := validatePointInput(cmd.PhysicalAddress.Point); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return CreateClinicResult{}, errors.Join(errs...)
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
		Description:    null.StringFromPtr(cmd.Description),
		PhysicalAddress: model.Address{
			Text: strings.TrimSpace(cmd.PhysicalAddress.Text),
		},
	}
	if cmd.PhysicalAddress.Point != nil {
		clinic.PhysicalAddress.Point = model.Point{
			Longitude: null.FloatFrom(cmd.PhysicalAddress.Point.Longitude),
			Latitude:  null.FloatFrom(cmd.PhysicalAddress.Point.Latitude),
		}
	}

	var result CreateClinicResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&clinic).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgErrCodeForeignKeyViolation {
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

		event := buildClinicCreatedEvent(&clinic)
		payload, err := anypb.New(event)
		if err != nil {
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicEventBuildFailed).
				With("clinic_id", id).
				Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(clinic.UpdatedAt),
			AggregateType: "clinic",
			AggregateId:   clinic.ID.String(),
			Payload:       payload,
		}
		if err := AppendOutboxEvent(tx, SubjectClinicCreated, envelope, nil); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
