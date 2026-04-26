// Package request is the gRPC transport for ServiceRequestCommandService.
package request

import (
	requestsvc "github.com/medincident/medincident-backend/internal/service/command/request"
	requestv1 "github.com/medincident/medincident-backend/pkg/command/request/v1"
)

// ServiceRequestHandler implements requestv1.ServiceRequestCommandServiceServer.
type ServiceRequestHandler struct {
	requestv1.UnimplementedServiceRequestCommandServiceServer
	svc *requestsvc.ServiceRequestService
}

// NewServiceRequestHandler wires the handler with the service.
func NewServiceRequestHandler(svc *requestsvc.ServiceRequestService) *ServiceRequestHandler {
	return &ServiceRequestHandler{svc: svc}
}
