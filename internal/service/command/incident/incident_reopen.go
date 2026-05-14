package incident

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/outbox"
	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/validation"
)

type ReopenIncidentPayload struct {
	IncidentID string `validate:"required,uuid"`
}

type ReopenIncidentCommand struct {
	Caller  authz.Caller
	Payload ReopenIncidentPayload
}

type ReopenIncidentResult struct {
	NewIncidentID uuid.UUID
}

// Reopen creates a fresh incident pointing back at the original via
// reopened_from_incident_id. Allowed only when the source is done or
// rejected (NOT cancelled). The new incident inherits org/clinic/dept,
// category, type, source_patient_zitadel_user_id; description is empty;
// priority resets to normal; occurred_at is now; registrar is the caller.
//
// See: docs/services/incident/Incidents.md
func (s *IncidentService) Reopen(
	ctx context.Context, cmd ReopenIncidentCommand,
) (ReopenIncidentResult, error) {
	if err := validation.Struct(cmd.Payload); err != nil {
		return ReopenIncidentResult{}, err
	}
	srcID := uuid.MustParse(cmd.Payload.IncidentID)
	now := time.Now()

	newID, err := uuid.NewV7()
	if err != nil {
		return ReopenIncidentResult{}, oops.In(scope).
			Code(ErrCodeIncidentIDGenerationFailed).Wrap(err)
	}

	var result ReopenIncidentResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		src, err := s.loadIncident(tx, srcID)
		if err != nil {
			return err
		}
		if src.Status != model.IncidentStatusDone && src.Status != model.IncidentStatusRejected {
			return oops.In(scope).Code(ErrCodeIncidentNotReopenable).
				Public("Only completed or rejected incidents can be reopened.").
				With("status", src.Status).Errorf("not reopenable")
		}
		if err := s.authz.Require(ctx, cmd.Caller.ZitadelUserID,
			privilegedActorPolicy(src.OrganizationID, src.ClinicID, src.DepartmentID)); err != nil {
			return err
		}
		registrarEmp, err := s.loadRegistrarEmployee(tx, cmd.Caller.ZitadelUserID, src.OrganizationID)
		if err != nil {
			return err
		}
		var regDept model.Department
		if err := tx.First(&regDept, "id = ?", registrarEmp.DepartmentID).Error; err != nil {
			return oops.In(scope).Code(ErrCodeIncidentLoadFailed).Wrap(err)
		}

		newInc := model.Incident{
			ID:                         newID,
			OrganizationID:             src.OrganizationID,
			ClinicID:                   src.ClinicID,
			DepartmentID:               src.DepartmentID,
			CategoryID:                 src.CategoryID,
			TypeID:                     src.TypeID,
			Status:                     model.IncidentStatusPending,
			Priority:                   model.IncidentPriorityNormal,
			Description:                null.String{},
			OccurredAt:                 now,
			RegistrarEmployeeID:        registrarEmp.ID,
			SourcePatientZitadelUserID: src.SourcePatientZitadelUserID,
			ReopenedFromIncidentID:     uuid.NullUUID{UUID: src.ID, Valid: true},
			CreatedAt:                  now,
			UpdatedAt:                  now,
		}
		if err := tx.Create(&newInc).Error; err != nil {
			return oops.In(scope).Code(ErrCodeIncidentSaveFailed).Wrap(err)
		}
		env, err := buildIncidentCreatedEnvelope(&newInc, cmd.Caller.ZitadelUserID, registrarEmp, regDept.ClinicID)
		if err != nil {
			return err
		}
		if err := outbox.Append(tx, "medincident.event.incident.v1.created", env); err != nil {
			return err
		}
		result.NewIncidentID = newID
		return nil
	})
	return result, err
}
