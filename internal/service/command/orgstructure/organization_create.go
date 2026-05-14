package orgstructure

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	anypb "google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// Error codes emitted by Organization-aggregate commands that are not
// struct-tag validation (infrastructure / existence / concurrency).
// Struct-tag violations are collapsed into a single validation_failed
// oops error by the validation translator.
const (
	ErrCodeOrganizationIDGenerationFailed  = "organization_id_generation_failed"
	ErrCodeOrganizationSaveFailed          = "organization_save_failed"
	ErrCodeOrganizationLoadFailed          = "organization_load_failed"
	ErrCodeOrganizationNotFound            = "organization_not_found"
	ErrCodeOrganizationDeleteFailed        = "organization_delete_failed"
	ErrCodeOrganizationDeleteHasDependents = "organization_delete_has_dependents"
)

// CreateOrganizationPayload is the validated client-facing payload of
// CreateOrganization. Scalar fields are primitives the transport layer
// can hand straight through; LegalAddress is a nested validated
// AddressInput and must be mapped from the corresponding nested
// transport/proto message.
type CreateOrganizationPayload struct {
	Name         string       `validate:"required,no_extra_ws,min=2,max=256"`
	Description  *string      `validate:"omitnil,no_extra_ws,min=8,max=2048"`
	LegalAddress AddressInput `validate:"required"`
}

// CreateOrganizationCommand is the input of OrganizationService.Create.
// It bundles the authenticated caller (already validated at the
// transport boundary) with the request payload so the service layer is
// self-contained: validate → authorize → execute.
type CreateOrganizationCommand struct {
	Caller  authz.Caller
	Payload CreateOrganizationPayload
}

// CreateOrganizationResult is the output of OrganizationService.Create.
type CreateOrganizationResult struct {
	ID uuid.UUID
}

// buildOrgAddressProto converts a model.Address to the org proto Address.
func buildOrgAddressProto(a model.Address) *orgv1.Address {
	addr := &orgv1.Address{Text: a.Text}
	if a.Point.Valid {
		addr.Point = &orgv1.Point{
			Longitude: a.Point.V.Longitude,
			Latitude:  a.Point.V.Latitude,
		}
	}
	return addr
}

func buildOrganizationCreatedEnvelope(org *model.Organization) (*eventv1.Envelope, error) {
	addr := &orgv1.Address{Text: org.LegalAddress.Text}
	if org.LegalAddress.Point.Valid {
		addr.Point = &orgv1.Point{
			Longitude: org.LegalAddress.Point.V.Longitude,
			Latitude:  org.LegalAddress.Point.V.Latitude,
		}
	}
	var desc string
	if org.Description.Valid {
		desc = org.Description.String
	}
	msg := &orgv1.OrganizationCreated{
		Name:         org.Name,
		Description:  desc,
		LegalAddress: addr,
		CreatedAt:    timestamppb.New(org.CreatedAt),
	}
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, oops.In("services.orgstructure.organization").
			Code(ErrCodeOrganizationSaveFailed).
			Wrap(err)
	}
	return &eventv1.Envelope{
		OccurredAt:    timestamppb.New(org.CreatedAt),
		AggregateType: "organization",
		AggregateId:   org.ID.String(),
		Payload:       payload,
	}, nil
}

// Create persists a new Organization and writes the matching
// outbox event in one transaction.
//
// See: docs/services/OrgStructure.md
func (s *OrganizationService) Create(
	ctx context.Context,
	cmd CreateOrganizationCommand,
) (CreateOrganizationResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return CreateOrganizationResult{}, err
	}
	if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.SystemAdmin); err != nil {
		return CreateOrganizationResult{}, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return CreateOrganizationResult{}, oops.In("services.orgstructure.organization").
			Code(ErrCodeOrganizationIDGenerationFailed).
			Public("Failed to create organization.").
			Wrap(err)
	}

	org := model.Organization{
		ID:       id,
		Name:     strings.TrimSpace(cmd.Payload.Name),
		IsActive: true,
		LegalAddress: model.Address{
			Text: strings.TrimSpace(cmd.Payload.LegalAddress.Text),
		},
	}
	if cmd.Payload.Description != nil {
		org.Description = null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
	}
	if cmd.Payload.LegalAddress.Point != nil {
		org.LegalAddress.Point = null.ValueFrom(model.Point{
			Longitude: cmd.Payload.LegalAddress.Point.Longitude,
			Latitude:  cmd.Payload.LegalAddress.Point.Latitude,
		})
	}

	var result CreateOrganizationResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&org).Error; err != nil {
			return oops.In("services.orgstructure.organization").
				Code(ErrCodeOrganizationSaveFailed).
				With("organization_id", id).
				Wrap(err)
		}

		env, err := buildOrganizationCreatedEnvelope(&org)
		if err != nil {
			return err
		}
		if err := outbox.Append(tx, "medincident.event.organization.v1.created", env); err != nil {
			return err
		}
		result.ID = id
		return nil
	})
	return result, err
}
