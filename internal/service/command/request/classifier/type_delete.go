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

// DeleteRequestTypePayload identifies the request type to delete.
type DeleteRequestTypePayload struct {
	TypeID string `validate:"required,uuid"`
}

// DeleteRequestTypeCommand = caller + payload.
type DeleteRequestTypeCommand struct {
	Caller  authz.Caller
	Payload DeleteRequestTypePayload
}

// DeleteRequestTypeResult is empty.
type DeleteRequestTypeResult struct{}

// Delete permanently removes a request type.
// See: docs/services/request/Classifier.md
func (s *RequestTypeService) Delete(
	ctx context.Context,
	cmd DeleteRequestTypeCommand,
) (DeleteRequestTypeResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return DeleteRequestTypeResult{}, err
	}
	typeID := uuid.MustParse(cmd.Payload.TypeID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.RequestType(typeID)); err != nil {
		return DeleteRequestTypeResult{}, err
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
		deletedAt := time.Now()
		env, err := buildRequestTypeDeletedEnvelope(row.ID, deletedAt)
		if err != nil {
			return err
		}
		if err := tx.Delete(&model.RequestType{}, "id = ?", row.ID).Error; err != nil {
			return oops.In(scope).Code(ErrCodeRequestTypeSaveFailed).
				With("request_type_id", row.ID).Wrap(err)
		}
		return outbox.Append(tx, "medincident.event.request_type.v1.deleted", env)
	})
	return DeleteRequestTypeResult{}, err
}

func buildRequestTypeDeletedEnvelope(typeID uuid.UUID, deletedAt time.Time) (*eventv1.Envelope, error) {
	msg := &requesttypev1.RequestTypeDeleted{
		TypeId:    typeID.String(),
		DeletedAt: timestamppb.New(deletedAt),
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scope).Code(ErrCodeRequestTypeSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(deletedAt),
		AggregateType: "request_type",
		AggregateId:   typeID.String(),
		Payload:       payload,
	}, nil
}
