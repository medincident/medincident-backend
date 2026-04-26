package membership

import (
	"github.com/google/uuid"
	"github.com/samber/oops"
)

// Error codes emitted by handler-layer request parsing.
const (
	ErrCodeHandlerInvalidEmployeeID     = "handler_employee_id_invalid"
	ErrCodeHandlerInvalidDepartmentID   = "handler_department_id_invalid"
	ErrCodeHandlerInvalidClinicID       = "handler_clinic_id_invalid"
	ErrCodeHandlerInvalidOrganizationID = "handler_organization_id_invalid"
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
