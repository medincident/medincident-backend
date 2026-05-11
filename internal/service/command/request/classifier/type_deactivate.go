package classifier

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	anypb "google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	requesttypev1 "github.com/medincident/medincident-backend/pkg/event/request_type/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// DeactivateRequestTypePayload identifies the request type to deactivate.
type DeactivateRequestTypePayload struct {
	TypeID string `validate:"required,uuid"`
}

// DeactivateRequestTypeCommand = caller + payload.
type DeactivateRequestTypeCommand struct {
	Caller  authz.Caller
	Payload DeactivateRequestTypePayload
}

// DeactivateRequestTypeResult is empty.
type DeactivateRequestTypeResult struct{}

// Deactivate marks a request type as inactive. Idempotent.
// See: docs/services/request/Classifier.md
func (s *RequestTypeService) Deactivate(
	ctx context.Context,
	cmd DeactivateRequestTypeCommand,
) (DeactivateRequestTypeResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return DeactivateRequestTypeResult{}, err
	}
	typeID := uuid.MustParse(cmd.Payload.TypeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.RequestType(typeID)); err != nil {
		return DeactivateRequestTypeResult{}, err
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row model.RequestType
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			First(&row, "id = ?", typeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return oops.In(scope).
					Code(ErrCodeRequestTypeNotFound).
					Public("Request type not found.").
					With("request_type_id", typeID).Wrap(err)
			}
			return oops.In(scope).Code(ErrCodeRequestTypeLoadFailed).
				With("request_type_id", typeID).Wrap(err)
		}
		if !row.IsActive {
			return nil
		}
		now := time.Now()
		row.IsActive = false
		row.UpdatedAt = now
		if err := tx.Save(&row).Error; err != nil {
			return oops.In(scope).Code(ErrCodeRequestTypeSaveFailed).
				With("request_type_id", row.ID).Wrap(err)
		}
		env, err := buildRequestTypeDeactivatedEnvelope(row.ID, now)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.request_type.v1.deactivated", env)
	})
	return DeactivateRequestTypeResult{}, err
}

func buildRequestTypeDeactivatedEnvelope(typeID uuid.UUID, updatedAt time.Time) (*eventv1.Envelope, error) {
	msg := &requesttypev1.RequestTypeDeactivated{
		TypeId:    typeID.String(),
		UpdatedAt: timestamppb.New(updatedAt),
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scope).Code(ErrCodeRequestTypeSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(updatedAt),
		AggregateType: "request_type",
		AggregateId:   typeID.String(),
		Payload:       payload,
	}, nil
}
