// Package membership is the gRPC transport for MembershipService.
// Handler methods translate proto request objects into service Command
// structs, call the service, and translate the result back. Error
// semantics are owned by the service layer.
package membership

import (
	"github.com/rs/zerolog"

	membershipv1 "github.com/medincident/medincident-command-service/gen/api/medincident/service/membership/v1"
	"github.com/medincident/medincident-command-service/internal/services/membership"
)

// MembershipHandler implements membershipv1.MembershipServiceServer.
type MembershipHandler struct {
	membershipv1.UnimplementedMembershipServiceServer

	empSvc *membership.EmployeeService
	logger *zerolog.Logger
}

// NewMembershipHandler wires the handler with EmployeeService and a logger.
func NewMembershipHandler(empSvc *membership.EmployeeService, logger *zerolog.Logger) *MembershipHandler {
	return &MembershipHandler{empSvc: empSvc, logger: logger}
}
