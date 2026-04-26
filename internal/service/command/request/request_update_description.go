package request

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
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
// See: https://github.com/medincident/medincident-backend/wiki/Service-Requests#updateservicerequestdescription
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
		return projector.ServiceRequestDescriptionUpdated(tx, sr.ID, sr.Description, now)
	})
}
