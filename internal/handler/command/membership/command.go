// Package membership is the gRPC transport for MembershipCommandService.
// Handler methods translate proto request objects into service Command
// + Payload shapes and wrap the authenticated caller in authz.Caller.
// Authorization runs inside each service method, so the handler does
// not depend on *authz.Authz, and UUID-parse guards are no longer
// needed at the transport boundary.
package membership

import (
	"github.com/medincident/medincident-backend/internal/service/command/membership"
	membershipv1 "github.com/medincident/medincident-backend/pkg/command/membership/v1"
)

// MembershipHandler implements membershipv1.MembershipCommandServiceServer.
type MembershipHandler struct {
	membershipv1.UnimplementedMembershipCommandServiceServer

	empSvc *membership.EmployeeService
}

// NewMembershipHandler wires the handler with EmployeeService.
func NewMembershipHandler(empSvc *membership.EmployeeService) *MembershipHandler {
	return &MembershipHandler{empSvc: empSvc}
}
