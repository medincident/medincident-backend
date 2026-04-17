// Package classifier is the gRPC transport for IncidentClassifierService.
package classifier

import (
	"github.com/google/uuid"
	"github.com/samber/oops"

	"github.com/medincident/medincident-command-service/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-command-service/pkg/service/incident/classifier/v1"
)

// Error codes emitted by handler-layer request parsing.
const (
	ErrCodeHandlerInvalidOrganizationID     = "handler_invalid_organization_id"
	ErrCodeHandlerInvalidIncidentCategoryID = "handler_invalid_incident_category_id"
	ErrCodeHandlerInvalidIncidentTypeID     = "handler_invalid_incident_type_id"
)

// IncidentClassifierHandler implements incidentclassifierv1.IncidentClassifierServiceServer.
type IncidentClassifierHandler struct {
	incidentclassifierv1.UnimplementedIncidentClassifierServiceServer

	authz       *authz.Authz
	categorySvc *classifiersvc.IncidentCategoryService
	typeSvc     *classifiersvc.IncidentTypeService
}

// NewIncidentClassifierHandler wires the handler with the two services.
func NewIncidentClassifierHandler(
	categorySvc *classifiersvc.IncidentCategoryService,
	typeSvc *classifiersvc.IncidentTypeService,
	az *authz.Authz,
) *IncidentClassifierHandler {
	return &IncidentClassifierHandler{
		authz:       az,
		categorySvc: categorySvc,
		typeSvc:     typeSvc,
	}
}

func parseOrganizationID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.incident.classifier").
			Code(ErrCodeHandlerInvalidOrganizationID).
			Public("Invalid organization id.").
			With("organization_id", raw).
			Wrap(err)
	}
	return id, nil
}

func parseIncidentCategoryID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.incident.classifier").
			Code(ErrCodeHandlerInvalidIncidentCategoryID).
			Public("Invalid incident category id.").
			With("incident_category_id", raw).
			Wrap(err)
	}
	return id, nil
}

func parseIncidentTypeID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.incident.classifier").
			Code(ErrCodeHandlerInvalidIncidentTypeID).
			Public("Invalid incident type id.").
			With("incident_type_id", raw).
			Wrap(err)
	}
	return id, nil
}
