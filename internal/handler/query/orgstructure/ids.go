package orgstructure

import (
	"github.com/google/uuid"
	"github.com/samber/oops"
)

// Error codes emitted by handler-layer request parsing.
const (
	ErrCodeHandlerInvalidOrganizationID = "handler_invalid_organization_id"
	ErrCodeHandlerInvalidClinicID       = "handler_invalid_clinic_id"
	ErrCodeHandlerInvalidDepartmentID   = "handler_invalid_department_id"
)

func parseUUID(raw, field, code string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.orgstructure").
			Code(code).
			Public(field+" is not a valid UUID.").
			With(field, raw).
			Wrap(err)
	}
	return id, nil
}

func parseOrganizationID(raw string) (uuid.UUID, error) {
	return parseUUID(raw, "organization_id", ErrCodeHandlerInvalidOrganizationID)
}

func parseClinicID(raw string) (uuid.UUID, error) {
	return parseUUID(raw, "clinic_id", ErrCodeHandlerInvalidClinicID)
}

func parseDepartmentID(raw string) (uuid.UUID, error) {
	return parseUUID(raw, "department_id", ErrCodeHandlerInvalidDepartmentID)
}
