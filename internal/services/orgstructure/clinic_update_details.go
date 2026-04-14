package orgstructure

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	clinicv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/clinic/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
)

type UpdateClinicDetailsCommand struct {
	ID          uuid.UUID
	Name        string
	Description *string
}

func buildClinicDetailsChangedEvent(c *model.Clinic) *clinicv1.ClinicDetailsChanged {
	ev := &clinicv1.ClinicDetailsChanged{Name: c.Name}
	if c.Description.Valid {
		desc := c.Description.String
		ev.Description = &desc
	}
	return ev
}

func (s *ClinicService) UpdateDetails(
	ctx context.Context,
	cmd UpdateClinicDetailsCommand,
) error {
	var errs []error
	if err := validateClinicName(cmd.Name); err != nil {
		errs = append(errs, err)
	}
	if err := validateClinicDescription(cmd.Description); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var clinic model.Clinic
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&clinic, "id = ?", cmd.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.orgstructure.clinic").
					Code(ErrCodeClinicNotFound).
					Public("Clinic not found.").
					With("clinic_id", cmd.ID).
					Errorf("clinic not found")
			}
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicLoadFailed).
				With("clinic_id", cmd.ID).
				Wrap(err)
		}

		newName := strings.TrimSpace(cmd.Name)
		newDesc := null.StringFromPtr(cmd.Description)
		if clinic.Name == newName && clinic.Description == newDesc {
			return nil
		}
		clinic.Name = newName
		clinic.Description = newDesc
		if err := tx.Save(&clinic).Error; err != nil {
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicSaveFailed).
				With("clinic_id", cmd.ID).
				Wrap(err)
		}

		event := buildClinicDetailsChangedEvent(&clinic)
		payload, err := anypb.New(event)
		if err != nil {
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicEventBuildFailed).
				With("clinic_id", cmd.ID).
				Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(clinic.UpdatedAt),
			AggregateType: AggregateTypeClinic,
			AggregateId:   clinic.ID.String(),
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectClinicDetailsChanged, envelope)
	})
}
