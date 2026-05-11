package classifier

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	classifierv1 "github.com/medincident/medincident-backend/pkg/event/incident/classifier/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// UpdateIncidentTypeDetailsPayload carries the new name and (optional)
// description for an existing incident type.
type UpdateIncidentTypeDetailsPayload struct {
	TypeID      string  `validate:"required,uuid"`
	Name        string  `validate:"required,no_extra_ws,min=2,max=256"`
	Description *string `validate:"omitnil,no_extra_ws,min=8,max=2048"`
}

// UpdateIncidentTypeDetailsCommand = caller + payload.
type UpdateIncidentTypeDetailsCommand struct {
	Caller  authz.Caller
	Payload UpdateIncidentTypeDetailsPayload
}

// UpdateIncidentTypeDetailsResult is empty.
type UpdateIncidentTypeDetailsResult struct{}

// UpdateDetails changes the name and description of an incident type.
//
// See: docs/services/incident/Classifier.md
func (s *IncidentTypeService) UpdateDetails(
	ctx context.Context,
	cmd UpdateIncidentTypeDetailsCommand,
) (UpdateIncidentTypeDetailsResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return UpdateIncidentTypeDetailsResult{}, err
	}
	typeID := uuid.MustParse(cmd.Payload.TypeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.IncidentType(typeID)); err != nil {
		return UpdateIncidentTypeDetailsResult{}, err
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model.IncidentType
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&row, "id = ?", typeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In("services.incident.classifier.type").
					Code(ErrCodeIncidentTypeNotFound).
					Public("Incident type not found.").
					With("incident_type_id", typeID).
					Wrap(err)
			}
			return oops.In("services.incident.classifier.type").
				Code(ErrCodeIncidentTypeLoadFailed).
				With("incident_type_id", typeID).
				Wrap(err)
		}

		if err := lockClassifierOrg(tx, row.OrganizationID); err != nil {
			return err
		}

		newName := strings.TrimSpace(cmd.Payload.Name)
		var newDescription null.String
		if cmd.Payload.Description != nil {
			newDescription = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
		}
		if row.Name == newName && row.Description == newDescription {
			return nil
		}
		row.Name = newName
		row.Description = newDescription
		if err := tx.Save(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
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

		env, err := buildIncidentTypeDetailsUpdatedEnvelope(&row)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.incident_type.v1.details_updated", env)
	})
	return UpdateIncidentTypeDetailsResult{}, err
}

func buildIncidentTypeDetailsUpdatedEnvelope(t *model.IncidentType) (*eventv1.Envelope, error) {
	msg := &classifierv1.IncidentTypeDetailsUpdated{TypeId: t.ID.String(), Name: t.Name, UpdatedAt: timestamppb.New(t.UpdatedAt)}
	if t.Description.Valid {
		msg.Description = wrapperspb.String(t.Description.String)
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In("services.incident.classifier.type").Code(ErrCodeIncidentTypeSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{OccurredAt: timestamppb.New(t.UpdatedAt), AggregateType: "incident_type", AggregateId: t.ID.String(), Payload: payload}, nil
}
