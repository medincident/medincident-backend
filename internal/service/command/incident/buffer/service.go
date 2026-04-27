// Package buffer is the write-side service for the patient incident
// buffer. Patients submit/update/cancel buffer entries; dispatchers
// publish or reject them. Once published the buffer row stays as a
// historical anchor; the materialised incident lives in
// internal/service/command/incident.
//
// See: docs/services/incident/Buffer.md
package buffer

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
)

const (
	ErrCodeBufferIDGenerationFailed = "buffer_id_generation_failed"
	ErrCodeBufferSaveFailed         = "buffer_save_failed"
	ErrCodeBufferLoadFailed         = "buffer_load_failed"
	ErrCodeBufferNotFound           = "buffer_not_found"
	ErrCodeBufferOrgNotFound        = "buffer_organization_not_found"
	ErrCodeBufferCategoryNotFound   = "buffer_category_not_found"
	ErrCodeBufferTypeNotFound       = "buffer_type_not_found"
	ErrCodeBufferTypeNotForPatients = "buffer_type_not_allowed_for_patients"
	ErrCodeBufferOccurredAtInvalid  = "buffer_occurred_at_invalid"
	ErrCodeBufferOccurredAtTooOld   = "buffer_occurred_at_too_old"
	ErrCodeBufferOccurredAtFuture   = "buffer_occurred_at_in_future"
	ErrCodeBufferNotPending         = "buffer_not_pending"
	ErrCodeBufferNotPatient         = "buffer_not_patient_owner"
	ErrCodeBufferDeptNotFound       = "buffer_department_not_found"
	ErrCodeBufferDispatcherNotFound = "buffer_dispatcher_not_found"
)

const bufferMaxOccurredAtAge = 7 * 24 * time.Hour

const scope = "services.command.incident.buffer"

type BufferService struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

func NewBufferService(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *BufferService {
	return &BufferService{db: db, authz: az, logger: logger}
}

// loadBuffer reads a buffer row by id; not-found is mapped to a public error.
func (s *BufferService) loadBuffer(tx *gorm.DB, id uuid.UUID) (*model.PatientIncidentBuffer, error) {
	var b model.PatientIncidentBuffer
	if err := tx.First(&b, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, oops.In(scope).
				Code(ErrCodeBufferNotFound).
				Public("Patient incident not found.").
				With("buffer_id", id).Wrap(err)
		}
		return nil, oops.In(scope).Code(ErrCodeBufferLoadFailed).Wrap(err)
	}
	return &b, nil
}

// validateOccurredAt enforces the 7-day window.
func validateOccurredAt(t, now time.Time) error {
	if t.After(now) {
		return oops.In(scope).Code(ErrCodeBufferOccurredAtFuture).
			Public("occurred_at must not be in the future.").
			With("occurred_at", t).Errorf("future occurred_at")
	}
	if now.Sub(t) > bufferMaxOccurredAtAge {
		return oops.In(scope).Code(ErrCodeBufferOccurredAtTooOld).
			Public("occurred_at exceeds the 7-day patient submission window.").
			With("occurred_at", t).Errorf("too old")
	}
	return nil
}

// occurredAtOrNow returns the buffer's occurred_at if set, else `now`.
// Used by Publish to materialise an incident from a buffer entry that
// the patient may have submitted without a specific timestamp.
func occurredAtOrNow(t null.Time, now time.Time) time.Time {
	if t.Valid {
		return t.Time
	}
	return now
}

// validatePatientCategoryType ensures (when supplied) the category +
// type belong to the org and the type is patient-allowed and active.
func validatePatientCategoryType(
	tx *gorm.DB, orgID uuid.UUID, categoryID, typeID uuid.NullUUID,
) error {
	if categoryID.Valid {
		var cat model.IncidentCategory
		if err := tx.First(&cat, "id = ?", categoryID.UUID).Error; err != nil {
			return oops.In(scope).Code(ErrCodeBufferCategoryNotFound).
				With("category_id", categoryID.UUID).Wrap(err)
		}
		if cat.OrganizationID != orgID || !cat.IsActive {
			return oops.In(scope).Code(ErrCodeBufferCategoryNotFound).
				Public("Category is not available.").Errorf("invalid category")
		}
	}
	if typeID.Valid {
		var typ model.IncidentType
		if err := tx.First(&typ, "id = ?", typeID.UUID).Error; err != nil {
			return oops.In(scope).Code(ErrCodeBufferTypeNotFound).
				With("type_id", typeID.UUID).Wrap(err)
		}
		if typ.OrganizationID != orgID || !typ.IsActive {
			return oops.In(scope).Code(ErrCodeBufferTypeNotFound).
				Public("Type is not available.").Errorf("invalid type")
		}
		if !typ.IsAllowedForPatients {
			return oops.In(scope).Code(ErrCodeBufferTypeNotForPatients).
				Public("Type is not available for patients.").
				With("type_id", typ.ID).Errorf("not patient-allowed")
		}
		if categoryID.Valid && typ.CategoryID != categoryID.UUID {
			return oops.In(scope).Code(ErrCodeBufferTypeNotFound).
				Public("Type does not belong to the chosen category.").
				Errorf("category mismatch")
		}
	}
	return nil
}
