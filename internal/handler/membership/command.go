// Package membership is the gRPC transport for MembershipService.
// Handler methods translate proto request objects into service Command
// structs, call the service, and translate the result back. Error
// semantics are owned by the service layer.
package membership

import (
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"

	"github.com/medincident/medincident-command-service/internal/service/authz"
	"github.com/medincident/medincident-command-service/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-command-service/pkg/command/membership/v1"
)

// idErrs collects UUID-parse errors from a handler so the transport
// can emit every violation in one response instead of fail-fast on
// the first bad field. Used via the parse method and joined at the
// end with err().
type idErrs []error

func (e *idErrs) parse(raw string, fn func(string) (uuid.UUID, error)) uuid.UUID {
	id, err := fn(raw)
	if err != nil {
		*e = append(*e, err)
	}
	return id
}

func (e idErrs) err() error {
	if len(e) == 0 {
		return nil
	}
	return errors.Join(e...)
}

// MembershipHandler implements membershipv1.MembershipCommandServiceServer.
type MembershipHandler struct {
	membershipv1.UnimplementedMembershipCommandServiceServer

	authz  *authz.Authz
	empSvc *membership.EmployeeService
}

// NewMembershipHandler wires the handler with EmployeeService.
func NewMembershipHandler(empSvc *membership.EmployeeService, az *authz.Authz) *MembershipHandler {
	return &MembershipHandler{authz: az, empSvc: empSvc}
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

func parseClinicID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.membership").
			Code(ErrCodeHandlerInvalidClinicID).
			Public("clinic_id is not a valid UUID.").
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

func parseDeputyEmployeeID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.membership").
			Code(ErrCodeHandlerInvalidDeputyEmployeeID).
			Public("deputy_employee_id is not a valid UUID.").
			Wrap(err)
	}
	return id, nil
}

func parseOrganizationID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.membership").
			Code(ErrCodeHandlerInvalidOrganizationID).
			Public("organization_id is not a valid UUID.").
			Wrap(err)
	}
	return id, nil
}
