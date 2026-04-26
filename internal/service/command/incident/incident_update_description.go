package incident

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
	"github.com/medincident/medincident-backend/internal/service/validation"
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

		desc := null.StringFrom(strings.TrimSpace(*cmd.Payload.Description))
		inc.Description = desc
		inc.UpdatedAt = now
		if err := tx.Save(inc).Error; err != nil {
			return oops.In(scope).Code(ErrCodeIncidentSaveFailed).Wrap(err)
		}
		return projector.IncidentDescriptionUpdated(tx, inc.ID, desc, now)
	})
}
