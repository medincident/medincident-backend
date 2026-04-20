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

func parseOrganizationID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.orgstructure").
			Code(ErrCodeHandlerInvalidOrganizationID).
			Public("Invalid organization id.").
			With("organization_id", raw).
			Wrap(err)
	}
	return id, nil
}

func parseClinicID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.orgstructure").
			Code(ErrCodeHandlerInvalidClinicID).
			Public("Invalid clinic id.").
			With("clinic_id", raw).
			Wrap(err)
	}
	return id, nil
}

func parseDepartmentID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.orgstructure").
			Code(ErrCodeHandlerInvalidDepartmentID).
			Public("Invalid department id.").
			With("department_id", raw).
			Wrap(err)
	}
	return id, nil
}
