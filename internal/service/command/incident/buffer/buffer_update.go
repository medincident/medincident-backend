package buffer

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

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	bufferv1 "github.com/medincident/medincident-backend/pkg/event/incident/buffer/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

type UpdatePayload struct {
	BufferID    string  `validate:"required,uuid"`
	CategoryID  *string `validate:"omitnil,uuid"`
	TypeID      *string `validate:"omitnil,uuid"`
	Description *string `validate:"omitnil,no_extra_ws,min=1,max=10000"`
	Summary     *string `validate:"omitnil,no_extra_ws,min=1,max=10000"`
	Priority    *string `validate:"omitnil,oneof=normal high"`
	OccurredAt  *string
}

type UpdateCommand struct {
	Caller  authz.Caller
	Payload UpdatePayload
}

// Update lets the patient owner edit a pending buffer entry.
// All four input fields are independently mutable; passing nil keeps
// the current value (no clear-to-null semantics).
//
// See: docs/services/incident/Buffer.md
func (s *BufferService) Update(ctx context.Context, cmd UpdateCommand) error { //nolint:gocritic // hugeParam: Command is passed by value across the whole service layer for consistency.
	if err := validation.Struct(cmd.Payload); err != nil {
		return err
	}
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.Authenticated); err != nil {
		return err
	}
	id := uuid.MustParse(cmd.Payload.BufferID)
	now := time.Now()
	var occurred *null.Time
	if cmd.Payload.OccurredAt != nil {
		t, err := time.Parse(time.RFC3339Nano, *cmd.Payload.OccurredAt)
		if err != nil {
			return oops.In(scope).
				Code(ErrCodeBufferOccurredAtInvalid).
				Public("occurred_at is not a valid RFC3339 timestamp.").
				With("occurred_at", *cmd.Payload.OccurredAt).Wrap(err)
		}
		if err := validateOccurredAt(t, now); err != nil {
			return err
		}
		v := null.TimeFrom(t)
		occurred = &v
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		b, err := s.loadBuffer(tx, id)
		if err != nil {
			return err
		}
		if b.PatientZitadelUserID != cmd.Caller.ZitadelUserID {
			return oops.In(scope).Code(ErrCodeBufferNotPatient).
				Public("Only the submitting patient may edit this entry.").
				Errorf("not owner")
		}
		if b.Status != model.BufferStatusPending {
			return oops.In(scope).Code(ErrCodeBufferNotPending).
				Public("Buffer entry is no longer editable.").
				With("status", b.Status).Errorf("not pending")
		}

		if cmd.Payload.CategoryID != nil {
			b.CategoryID = uuid.NullUUID{UUID: uuid.MustParse(*cmd.Payload.CategoryID), Valid: true}
		}
		if cmd.Payload.TypeID != nil {
			b.TypeID = uuid.NullUUID{UUID: uuid.MustParse(*cmd.Payload.TypeID), Valid: true}
		}
		if err := validatePatientCategoryType(tx, b.OrganizationID, b.CategoryID, b.TypeID); err != nil {
			return err
		}
		summaryChanged := cmd.Payload.Summary != nil
		priorityChanged := cmd.Payload.Priority != nil
		if cmd.Payload.Description != nil {
			b.Description = strings.TrimSpace(*cmd.Payload.Description)
		}
		if cmd.Payload.Summary != nil {
			b.Summary = strings.TrimSpace(*cmd.Payload.Summary)
		}
		if cmd.Payload.Priority != nil {
			b.Priority = model.BufferPriority(*cmd.Payload.Priority)
		}
		if occurred != nil {
			b.OccurredAt = *occurred
		}
		b.UpdatedAt = now

		if err := tx.Save(b).Error; err != nil {
			return oops.In(scope).Code(ErrCodeBufferSaveFailed).Wrap(err)
		}
		env, err := buildPatientIncidentBufferUpdatedEnvelope(b, summaryChanged, priorityChanged)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.patient_incident_buffer.v1.updated", env)
	})
}

func buildPatientIncidentBufferUpdatedEnvelope(b *model.PatientIncidentBuffer, summaryChanged, priorityChanged bool) (*eventv1.Envelope, error) {
	msg := &bufferv1.PatientIncidentBufferUpdated{
		BufferId:    b.ID.String(),
		Description: b.Description,
		Status:      string(b.Status),
		UpdatedAt:   timestamppb.New(b.UpdatedAt),
	}
	if b.OccurredAt.Valid {
		msg.OccurredAt = timestamppb.New(b.OccurredAt.Time)
	}
	if b.CategoryID.Valid {
		msg.CategoryId = wrapperspb.String(b.CategoryID.UUID.String())
	}
	if b.TypeID.Valid {
		msg.TypeId = wrapperspb.String(b.TypeID.UUID.String())
	}
	if b.PublishedIncidentID.Valid {
		msg.PublishedIncidentId = wrapperspb.String(b.PublishedIncidentID.UUID.String())
	}
	if summaryChanged {
		msg.Summary = wrapperspb.String(b.Summary)
	}
	if priorityChanged {
		msg.Priority = wrapperspb.String(string(b.Priority))
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scope).Code(ErrCodeBufferSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(b.UpdatedAt),
		AggregateType: "patient_incident_buffer",
		AggregateId:   b.ID.String(),
		Payload:       payload,
	}, nil
}
