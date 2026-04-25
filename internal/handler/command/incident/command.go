// Package incident is the gRPC transport for IncidentCommandService.
package incident

import (
	incidentsvc "github.com/medincident/medincident-backend/internal/service/command/incident"
	incidentv1 "github.com/medincident/medincident-backend/pkg/command/incident/v1"
)

// IncidentHandler implements incidentv1.IncidentCommandServiceServer.
type IncidentHandler struct {
	incidentv1.UnimplementedIncidentCommandServiceServer
	svc *incidentsvc.IncidentService
}

// NewIncidentHandler wires the handler with the service.
func NewIncidentHandler(svc *incidentsvc.IncidentService) *IncidentHandler {
	return &IncidentHandler{svc: svc}
}
