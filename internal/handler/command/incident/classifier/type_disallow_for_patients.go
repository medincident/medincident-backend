package classifier

import (
	"context"

	"github.com/medincident/medincident-command-service/internal/middleware/grpcmw"
	"github.com/medincident/medincident-command-service/internal/service/authz"
	classifiersvc "github.com/medincident/medincident-command-service/internal/service/command/incident/classifier"
	incidentclassifierv1 "github.com/medincident/medincident-command-service/pkg/command/incident/classifier/v1"
)

func (h *IncidentClassifierHandler) DisallowIncidentTypeForPatients(
	ctx context.Context,
	req *incidentclassifierv1.DisallowIncidentTypeForPatientsRequest,
) (*incidentclassifierv1.DisallowIncidentTypeForPatientsResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := h.typeSvc.DisallowForPatients(ctx, classifiersvc.DisallowIncidentTypeForPatientsCommand{
		Caller:  authz.Caller{ZitadelUserID: callerID},
		Payload: classifiersvc.DisallowIncidentTypeForPatientsPayload{TypeID: req.GetTypeId()},
	}); err != nil {
		return nil, err
	}
	return &incidentclassifierv1.DisallowIncidentTypeForPatientsResponse{}, nil
}
