package membership

import (
	"github.com/google/uuid"
	"github.com/samber/oops"
)

// Error codes emitted by handler-layer request parsing.
const (
	ErrCodeHandlerInvalidEmployeeID     = "handler_invalid_employee_id"
	ErrCodeHandlerInvalidDepartmentID   = "handler_invalid_department_id"
	ErrCodeHandlerInvalidClinicID       = "handler_invalid_clinic_id"
	ErrCodeHandlerInvalidOrganizationID = "handler_invalid_organization_id"
)

func parseUUID(raw, field, code string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.membership").
			Code(code).
			Public(field+" is not a valid UUID.").
			With(field, raw).
			Wrap(err)
	}
	return id, nil
}

func parseEmployeeID(raw string) (uuid.UUID, error) {
	return parseUUID(raw, "employee_id", ErrCodeHandlerInvalidEmployeeID)
}

func parseDepartmentID(raw string) (uuid.UUID, error) {
	return parseUUID(raw, "department_id", ErrCodeHandlerInvalidDepartmentID)
}

func parseClinicID(raw string) (uuid.UUID, error) {
	return parseUUID(raw, "clinic_id", ErrCodeHandlerInvalidClinicID)
}

func parseOrganizationID(raw string) (uuid.UUID, error) {
	return parseUUID(raw, "organization_id", ErrCodeHandlerInvalidOrganizationID)
}
