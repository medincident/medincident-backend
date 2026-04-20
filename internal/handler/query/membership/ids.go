package membership

import (
	"github.com/google/uuid"
	"github.com/samber/oops"
)

const (
	ErrCodeHandlerInvalidEmployeeID     = "employee_id_invalid"
	ErrCodeHandlerInvalidDepartmentID   = "department_id_invalid"
	ErrCodeHandlerInvalidClinicID       = "clinic_id_invalid"
	ErrCodeHandlerInvalidOrganizationID = "organization_id_invalid"
)

func parseEmployeeID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.membership").
			Code(ErrCodeHandlerInvalidEmployeeID).
			Public("employee_id is not a valid UUID.").
			Wrap(err)
	}
	return id, nil
}

func parseDepartmentID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.membership").
			Code(ErrCodeHandlerInvalidDepartmentID).
			Public("department_id is not a valid UUID.").
			Wrap(err)
	}
	return id, nil
}

func parseClinicID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.membership").
			Code(ErrCodeHandlerInvalidClinicID).
			Public("clinic_id is not a valid UUID.").
			Wrap(err)
	}
	return id, nil
}

func parseOrganizationID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.membership").
			Code(ErrCodeHandlerInvalidOrganizationID).
			Public("organization_id is not a valid UUID.").
			Wrap(err)
	}
	return id, nil
}
