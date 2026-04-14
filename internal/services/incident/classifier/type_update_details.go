package classifier

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	typeeventv1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/incident/type/v1"
	envelopev1 "github.com/medincident/medincident-command-service/gen/api/medincident/event/v1"
	"github.com/medincident/medincident-command-service/internal/model"
	"github.com/medincident/medincident-command-service/internal/services/outbox"
)

type UpdateIncidentTypeDetailsCommand struct {
	TypeID      uuid.UUID
	Name        string
	Description *string
}

type UpdateIncidentTypeDetailsResult struct{}

func buildIncidentTypeDetailsChangedEvent(t *model.IncidentType) *typeeventv1.IncidentTypeDetailsChanged {
	ev := &typeeventv1.IncidentTypeDetailsChanged{Name: t.Name}
	if t.Description.Valid {
		d := t.Description.String
		ev.Description = &d
	}
	return ev
}

func (s *IncidentTypeService) UpdateDetails(
	ctx context.Context,
	cmd UpdateIncidentTypeDetailsCommand,
) (UpdateIncidentTypeDetailsResult, error) {
	var errs []error
	if err := validateIncidentTypeName(cmd.Name); err != nil {
		errs = append(errs, err)
	}
	if err := validateIncidentTypeDescription(cmd.Description); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return UpdateIncidentTypeDetailsResult{}, errors.Join(errs...)
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model.IncidentType
		if err := tx.First(&row, "id = ?", cmd.TypeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeNotFound).
					Public("Incident type not found.").
					With("incident_type_id", cmd.TypeID).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeLoadFailed).
				With("incident_type_id", cmd.TypeID).
				Wrap(err)
		}

		row.Name = strings.TrimSpace(cmd.Name)
		row.Description = null.StringFromPtr(trimmedStringPtr(cmd.Description))
		if err := tx.Save(&row).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgErrCodeUniqueViolation {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeNameConflict).
					Public("An active incident type with this name already exists.").
					With("incident_type_id", row.ID).
					With("name", row.Name).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeSaveFailed).
				With("incident_type_id", row.ID).
				Wrap(err)
		}

		event := buildIncidentTypeDetailsChangedEvent(&row)
		payload, err := anypb.New(event)
		if err != nil {
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeEventBuildFailed).
				With("incident_type_id", row.ID).
				Wrap(err)
		}
		envelope := &envelopev1.Envelope{
			OccurredAt:    timestamppb.New(row.UpdatedAt),
			AggregateType: AggregateTypeIncidentType,
			AggregateId:   row.ID.String(),
			Payload:       payload,
		}
		return outbox.AppendEvent(tx, SubjectIncidentTypeDetailsChanged, envelope, nil)
	})
	return UpdateIncidentTypeDetailsResult{}, err
}
