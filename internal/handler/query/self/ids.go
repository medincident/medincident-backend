package self

import (
	"github.com/google/uuid"
	"github.com/samber/oops"
)

const (
	ErrCodeHandlerInvalidOrganizationID = "handler_self_organization_id_invalid"
	ErrCodeHandlerInvalidClinicID       = "handler_self_clinic_id_invalid"
	ErrCodeHandlerInvalidDepartmentID   = "handler_self_department_id_invalid"
)

func parseUUID(raw, field, code string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.self").
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
