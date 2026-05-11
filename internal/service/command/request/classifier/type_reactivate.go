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

// ReactivateRequestTypePayload identifies the request type to reactivate.
type ReactivateRequestTypePayload struct {
	TypeID string `validate:"required,uuid"`
}

// ReactivateRequestTypeCommand = caller + payload.
type ReactivateRequestTypeCommand struct {
	Caller  authz.Caller
	Payload ReactivateRequestTypePayload
}

// ReactivateRequestTypeResult is empty.
type ReactivateRequestTypeResult struct{}

// Reactivate marks a request type as active. Idempotent. May conflict
// on name uniqueness if another active type with the same name exists.
// See: docs/services/request/Classifier.md
func (s *RequestTypeService) Reactivate(
	ctx context.Context,
	cmd ReactivateRequestTypeCommand,
) (ReactivateRequestTypeResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return ReactivateRequestTypeResult{}, err
	}
	typeID := uuid.MustParse(cmd.Payload.TypeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.RequestType(typeID)); err != nil {
		return ReactivateRequestTypeResult{}, err
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
		if row.IsActive {
			return nil
		}
		now := time.Now()
		row.IsActive = true
		row.UpdatedAt = now
		if err := tx.Save(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In(scope).
					Code(ErrCodeRequestTypeNameConflict).
					Public("An active request type with this name already exists.").
					With("organization_id", row.OrganizationID).
					With("name", row.Name).Wrap(err)
			}
			return oops.In(scope).Code(ErrCodeRequestTypeSaveFailed).
				With("request_type_id", row.ID).Wrap(err)
		}
		env, err := buildRequestTypeReactivatedEnvelope(row.ID, now)
		if err != nil {
			return err
		}
		return outbox.Append(tx, "medincident.event.request_type.v1.reactivated", env)
	})
	return ReactivateRequestTypeResult{}, err
}

func buildRequestTypeReactivatedEnvelope(typeID uuid.UUID, updatedAt time.Time) (*eventv1.Envelope, error) {
	msg := &requesttypev1.RequestTypeReactivated{
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
