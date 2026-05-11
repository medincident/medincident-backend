package request

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	anypb "google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	servicerequestv1 "github.com/medincident/medincident-backend/pkg/event/service_request/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// UpdateServiceRequestDescriptionPayload carries the new description.
type UpdateServiceRequestDescriptionPayload struct {
	ServiceRequestID string `validate:"required,uuid"`
	Description      string `validate:"required,no_extra_ws,min=1,max=10000"`
}

// UpdateServiceRequestDescriptionCommand = caller + payload.
type UpdateServiceRequestDescriptionCommand struct {
	Caller  authz.Caller
	Payload UpdateServiceRequestDescriptionPayload
}

// UpdateDescription updates the description of a service request.
// See: docs/services/request/Requests.md
func (s *ServiceRequestService) UpdateDescription(
	ctx context.Context,
	cmd UpdateServiceRequestDescriptionCommand,
) error {
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.ServiceRequestID)
	now := time.Now()

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		sr, err := s.loadServiceRequest(tx, id)
		if err != nil {
			return err
		}
		if isTerminalStatus(sr.Status) {
			return oops.In(scope).Code(ErrCodeServiceRequestFrozen).
				Public("Cannot modify a completed or cancelled request.").
				With("service_request_id", id).With("status", sr.Status).
				Errorf("frozen")
		}
		if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID,
			privilegedActorPolicy(sr.OrganizationID, sr.ClinicID, sr.DepartmentID)); err != nil {
			return err
		}
		sr.Description = strings.TrimSpace(cmd.Payload.Description)
		sr.UpdatedAt = now
		if err := tx.Save(sr).Error; err != nil {
			return oops.In(scope).Code(ErrCodeServiceRequestSaveFailed).Wrap(err)
		}
		env, err := buildServiceRequestDescriptionUpdatedEnvelope(sr.ID, sr.Description, now)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.service_request.v1.description_updated", env)
	})
}

func buildServiceRequestDescriptionUpdatedEnvelope(requestID uuid.UUID, description string, updatedAt time.Time) (*eventv1.Envelope, error) {
	msg := &servicerequestv1.ServiceRequestDescriptionUpdated{
		RequestId:   requestID.String(),
		Description: description,
		UpdatedAt:   timestamppb.New(updatedAt),
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scope).Code(ErrCodeServiceRequestSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(updatedAt),
		AggregateType: "service_request",
		AggregateId:   requestID.String(),
		Payload:       payload,
	}, nil
}
