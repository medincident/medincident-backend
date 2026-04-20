// Package classifier is the gRPC transport for IncidentClassifierCommandService.
// Handlers translate proto → service Command + Payload, wrap the
// authenticated caller in an authz.Caller, and return the proto
// response. Authorization happens inside each service method.
package classifier

import (
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/command/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-command-service/pkg/command/incident/classifier/v1"
)

// IncidentClassifierHandler implements incidentclassifierv1.IncidentClassifierCommandServiceServer.
type IncidentClassifierHandler struct {
	incidentclassifierv1.UnimplementedIncidentClassifierCommandServiceServer

	categorySvc *classifiersvc.IncidentCategoryService
	typeSvc     *classifiersvc.IncidentTypeService
}

// NewIncidentClassifierHandler wires the handler with the two services.
func NewIncidentClassifierHandler(
	categorySvc *classifiersvc.IncidentCategoryService,
	typeSvc *classifiersvc.IncidentTypeService,
) *IncidentClassifierHandler {
	return &IncidentClassifierHandler{
		categorySvc: categorySvc,
		typeSvc:     typeSvc,
	}
}
