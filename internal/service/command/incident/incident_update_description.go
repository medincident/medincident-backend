package incident

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	incidentv1 "github.com/medincident/medincident-backend/pkg/event/incident/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// UpdateIncidentDescriptionPayload requires a non-nil Description.
//
// Nil-as-clear was rejected by review: a buggy client that forgets the
// field would silently wipe the existing description. If the client
// genuinely wants to clear the description that is a separate
// operation; this RPC only sets a new (non-empty) value.
type UpdateIncidentDescriptionPayload struct {
	IncidentID  string  `validate:"required,uuid"`
	Description *string `validate:"required,no_extra_ws,min=1,max=10000"`
}

type UpdateIncidentDescriptionCommand struct {
	Caller  authz.Caller
	Payload UpdateIncidentDescriptionPayload
}

// UpdateDescription edits the body. Allowed for the registrar OR a
// privileged role, only while the incident is not terminal.
//
// See: docs/services/incident/Incidents.md
func (s *IncidentService) UpdateDescription(
	ctx context.Context, cmd UpdateIncidentDescriptionCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.IncidentID)
	now := time.Now()

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		inc, err := s.loadIncident(tx, id)
		if err != nil {
			return err
		}
		if inc.Status.IsTerminal() {
			return oops.In(scope).Code(ErrCodeIncidentFrozen).
				Public("Incident is frozen and cannot be modified.").
				With("status", inc.Status).Errorf("frozen")
		}

		isReg, err := s.callerIsActiveRegistrar(tx, cmd.Caller.ZitadelUserID,
			inc.RegistrarEmployeeID)
		if err != nil {
			return err
		}
		if !isReg {
			if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID,
				privilegedActorPolicy(inc.OrganizationID, inc.ClinicID, inc.DepartmentID)); err != nil {
				return err
			}
		}

		trimmed := strings.TrimSpace(*cmd.Payload.Description)
		inc.Description = null.StringFrom(trimmed)
		inc.UpdatedAt = now
		if err := tx.Save(inc).Error; err != nil {
			return oops.In(scope).Code(ErrCodeIncidentSaveFailed).Wrap(err)
		}
		env, err := buildIncidentDescriptionUpdatedEnvelope(inc.ID, &trimmed, now)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.incident.v1.description_updated", env)
	})
}

func buildIncidentDescriptionUpdatedEnvelope(incidentID uuid.UUID, description *string, updatedAt time.Time) (*eventv1.Envelope, error) {
	msg := &incidentv1.IncidentDescriptionUpdated{
		IncidentId: incidentID.String(),
		UpdatedAt:  timestamppb.New(updatedAt),
	}
	if description != nil {
		msg.Description = wrapperspb.String(*description)
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scope).Code(ErrCodeIncidentSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(updatedAt),
		AggregateType: "incident",
		AggregateId:   incidentID.String(),
		Payload:       payload,
	}, nil
}
