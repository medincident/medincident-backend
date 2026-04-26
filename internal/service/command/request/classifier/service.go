// Package classifier owns the request type classifier command-side logic.
// A flat, per-organisation classifier with active-name uniqueness.
//
// See: docs/services/request/Classifier.md
package classifier

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/service/authz"
)

const (
	ErrCodeRequestTypeIDGenerationFailed = "request_type_id_generation_failed"
	ErrCodeRequestTypeSaveFailed         = "request_type_save_failed"
	ErrCodeRequestTypeLoadFailed         = "request_type_load_failed"
	ErrCodeRequestTypeNotFound           = "request_type_not_found"
	ErrCodeRequestTypeNameConflict       = "request_type_name_conflict"
)

const scope = "services.command.request.classifier"

// RequestTypeService handles mutations of domain.request_types.
type RequestTypeService struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewRequestTypeService returns a service bound to the given gorm DB
// and authorization service.
func NewRequestTypeService(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *RequestTypeService {
	return &RequestTypeService{db: db, authz: az, logger: logger}
}
