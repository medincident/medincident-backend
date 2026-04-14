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

type UpdateClinicPhysicalAddressCommand struct {
	ID      uuid.UUID
	Address AddressInput
}

func buildClinicPhysicalAddressChangedEvent(c *model.Clinic) *clinicv1.ClinicPhysicalAddressChanged {
	ev := &clinicv1.ClinicPhysicalAddressChanged{
		PhysicalAddress: &clinicv1.Address{Text: c.PhysicalAddress.Text},
	}
	if c.PhysicalAddress.Point.Longitude.Valid && c.PhysicalAddress.Point.Latitude.Valid {
		ev.PhysicalAddress.Point = &clinicv1.Point{
			Longitude: c.PhysicalAddress.Point.Longitude.Float64,
			Latitude:  c.PhysicalAddress.Point.Latitude.Float64,
		}
	}
	return ev
}

func (s *ClinicService) UpdatePhysicalAddress(
	ctx context.Context,
	cmd UpdateClinicPhysicalAddressCommand,
) error {
	var errs []error
	if err := validateAddressInput(cmd.Address); err != nil {
		errs = append(errs, err)
	}
	if err := validatePointInput(cmd.Address.Point); err != nil {
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

		newAddress := model.Address{Text: strings.TrimSpace(cmd.Address.Text)}
		if cmd.Address.Point != nil {
			newAddress.Point = model.Point{
				Longitude: null.FloatFrom(cmd.Address.Point.Longitude),
				Latitude:  null.FloatFrom(cmd.Address.Point.Latitude),
			}
		}
		if clinic.PhysicalAddress == newAddress {
			return nil
		}
		clinic.PhysicalAddress = newAddress
		if err := tx.Save(&clinic).Error; err != nil {
			return oops.In("services.orgstructure.clinic").
				Code(ErrCodeClinicSaveFailed).
				With("clinic_id", cmd.ID).
				Wrap(err)
		}

		event := buildClinicPhysicalAddressChangedEvent(&clinic)
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
		return outbox.AppendEvent(tx, SubjectClinicPhysicalAddressChanged, envelope)
	})
}
