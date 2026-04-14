// Package membership is the gRPC transport for MembershipService.
// Handler methods translate proto request objects into service Command
// structs, call the service, and translate the result back. Error
// semantics are owned by the service layer.
package membership

import (
	"github.com/google/uuid"
	"github.com/samber/oops"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// MembershipHandler implements membershipv1.MembershipServiceServer.
type MembershipHandler struct {
	membershipv1.UnimplementedMembershipServiceServer

	empSvc *membership.EmployeeService
}

// NewMembershipHandler wires the handler with EmployeeService.
func NewMembershipHandler(empSvc *membership.EmployeeService) *MembershipHandler {
	return &MembershipHandler{empSvc: empSvc}
}

func parseEmployeeID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.membership").
			Code(ErrCodeHandlerInvalidEmployeeID).
			Public("employee_id is not a valid UUID.").
			Wrap(err)
	}
	return id, nil
}

func parseDepartmentID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.membership").
			Code(ErrCodeHandlerInvalidDepartmentID).
			Public("department_id is not a valid UUID.").
			Wrap(err)
	}
	return id, nil
}

func parseVacationID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.membership").
			Code(ErrCodeHandlerInvalidVacationID).
			Public("vacation_id is not a valid UUID.").
			Wrap(err)
	}
	return id, nil
}
