package classifier

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	anypb "google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	requesttypev1 "github.com/medincident/medincident-backend/pkg/event/request_type/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// CreateRequestTypePayload is the validated client-facing payload.
type CreateRequestTypePayload struct {
	OrganizationID string  `validate:"required,uuid"`
	Name           string  `validate:"required,no_extra_ws,min=2,max=256"`
	Description    *string `validate:"omitnil,no_extra_ws,min=8,max=2048"`
}

// CreateRequestTypeCommand = caller + payload.
type CreateRequestTypeCommand struct {
	Caller  authz.Caller
	Payload CreateRequestTypePayload
}

// CreateRequestTypeResult is the output of RequestTypeService.Create.
type CreateRequestTypeResult struct {
	ID uuid.UUID
}

// Create creates a new request type for an organization.
// See: docs/services/request/Classifier.md
func (s *RequestTypeService) Create(
	ctx context.Context,
	cmd CreateRequestTypeCommand,
) (CreateRequestTypeResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return CreateRequestTypeResult{}, err
	}
	orgID := uuid.MustParse(cmd.Payload.OrganizationID)
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.Organization(orgID)); err != nil {
		return CreateRequestTypeResult{}, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return CreateRequestTypeResult{}, oops.In(scope).
			Code(ErrCodeRequestTypeIDGenerationFailed).
			Public("Failed to create request type.").
			Wrap(err)
	}

	var result CreateRequestTypeResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		row := model.RequestType{
			ID:             id,
			OrganizationID: orgID,
			Name:           strings.TrimSpace(cmd.Payload.Name),
			IsActive:       true,
			CreatedAt:      tx.NowFunc(),
			UpdatedAt:      tx.NowFunc(),
		}
		if cmd.Payload.Description != nil {
			row.Description = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
		}
		if err := tx.Create(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return oops.In(scope).
					Code(ErrCodeRequestTypeNameConflict).
					Public("An active request type with this name already exists.").
					With("organization_id", orgID).
					With("name", row.Name).
					Wrap(err)
			}
			return oops.In(scope).
				Code(ErrCodeRequestTypeSaveFailed).
				With("request_type_id", id).
				Wrap(err)
		}
		env, err := buildRequestTypeCreatedEnvelope(&row)
		if err != nil {
			return err
		}
		if err := outbox.Append(tx, "medincident.event.request_type.v1.created", env); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}

func buildRequestTypeCreatedEnvelope(rt *model.RequestType) (*eventv1.Envelope, error) {
	msg := &requesttypev1.RequestTypeCreated{
		TypeId:         rt.ID.String(),
		OrganizationId: rt.OrganizationID.String(),
		Name:           rt.Name,
		IsActive:       rt.IsActive,
		CreatedAt:      timestamppb.New(rt.CreatedAt),
	}
	if rt.Description.Valid {
		msg.Description = wrapperspb.String(rt.Description.String)
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In(scope).Code(ErrCodeRequestTypeSaveFailed).Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(rt.CreatedAt),
		AggregateType: "request_type",
		AggregateId:   rt.ID.String(),
		Payload:       payload,
	}, nil
}
