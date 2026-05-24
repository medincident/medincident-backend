package buffer

import (
	"context"
	"errors"
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

type SubmitPayload struct {
	OrganizationID string  `validate:"required,uuid"`
	CategoryID     *string `validate:"omitnil,uuid"`
	TypeID         *string `validate:"omitnil,uuid"`
	Description    string  `validate:"required,no_extra_ws,min=1,max=10000"`
	Summary        string  `validate:"required,no_extra_ws,min=1,max=10000"`
	Priority       string  `validate:"required,oneof=normal high"`
	OccurredAt     *string
}

type SubmitCommand struct {
	Caller  authz.Caller
	Payload SubmitPayload
}

type SubmitResult struct {
	ID uuid.UUID
}

// Submit creates a new buffer entry for the authenticated patient.
// Authorization: any authenticated caller (Authenticated policy).
//
// See: docs/services/incident/Buffer.md
func (s *BufferService) Submit(ctx context.Context, cmd SubmitCommand) (SubmitResult, error) { //nolint:gocritic // hugeParam: Command is passed by value across the whole service layer for consistency.
	if err := validation.Struct(cmd.Payload); err != nil {
		return SubmitResult{}, err
	}
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.Authenticated); err != nil {
		return SubmitResult{}, err
	}
	orgID := uuid.MustParse(cmd.Payload.OrganizationID)
	now := time.Now()
	var occurredTime time.Time
	var hasOccurred bool
	if cmd.Payload.OccurredAt != nil {
		t, err := time.Parse(time.RFC3339Nano, *cmd.Payload.OccurredAt)
		if err != nil {
			return SubmitResult{}, oops.In(scope).
				Code(ErrCodeBufferOccurredAtInvalid).
				Public("occurred_at is not a valid RFC3339 timestamp.").
				With("occurred_at", *cmd.Payload.OccurredAt).Wrap(err)
		}
		if err := validateOccurredAt(t, now); err != nil {
			return SubmitResult{}, err
		}
		occurredTime = t
		hasOccurred = true
	}

	var categoryID, typeID uuid.NullUUID
	if cmd.Payload.CategoryID != nil {
		categoryID = uuid.NullUUID{UUID: uuid.MustParse(*cmd.Payload.CategoryID), Valid: true}
	}
	if cmd.Payload.TypeID != nil {
		typeID = uuid.NullUUID{UUID: uuid.MustParse(*cmd.Payload.TypeID), Valid: true}
	}

	id, err := uuid.NewV7()
	if err != nil {
		return SubmitResult{}, oops.In(scope).Code(ErrCodeBufferIDGenerationFailed).Wrap(err)
	}

	var result SubmitResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var org model.Organization
		if err := tx.First(&org, "id = ?", orgID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scope).Code(ErrCodeBufferOrgNotFound).
					Public("Organization not found.").
					With("organization_id", orgID).Wrap(err)
			}
			return oops.In(scope).Code(ErrCodeBufferLoadFailed).Wrap(err)
		}
		if err := validatePatientCategoryType(tx, orgID, categoryID, typeID); err != nil {
			return err
		}

		b := model.PatientIncidentBuffer{
			ID:                   id,
			OrganizationID:       orgID,
			PatientZitadelUserID: cmd.Caller.ZitadelUserID,
			CategoryID:           categoryID,
			TypeID:               typeID,
			Description:          strings.TrimSpace(cmd.Payload.Description),
			Summary:              strings.TrimSpace(cmd.Payload.Summary),
			Priority:             model.BufferPriority(cmd.Payload.Priority),
			Status:               model.BufferStatusPending,
			CreatedAt:            now,
			UpdatedAt:            now,
		}
		if hasOccurred {
			b.OccurredAt = null.TimeFrom(occurredTime)
		}
		if err := tx.Create(&b).Error; err != nil {
			return oops.In(scope).Code(ErrCodeBufferSaveFailed).Wrap(err)
		}
		env, err := buildPatientIncidentBufferCreatedEnvelope(&b)
		if err != nil {
			return err
		}
		if err := outbox.Append(tx, "medincident.event.patient_incident_buffer.v1.created", env); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}

func buildPatientIncidentBufferCreatedEnvelope(b *model.PatientIncidentBuffer) (*eventv1.Envelope, error) {
	msg := &bufferv1.PatientIncidentBufferCreated{
		BufferId:             b.ID.String(),
		OrganizationId:       b.OrganizationID.String(),
		PatientZitadelUserId: b.PatientZitadelUserID,
		Description:          b.Description,
		Summary:              b.Summary,
		Priority:             string(b.Priority),
		Status:               string(b.Status),
		CreatedAt:            timestamppb.New(b.CreatedAt),
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
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scope).Code(ErrCodeBufferSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(b.CreatedAt),
		AggregateType: "patient_incident_buffer",
		AggregateId:   b.ID.String(),
		Payload:       payload,
	}, nil
}
