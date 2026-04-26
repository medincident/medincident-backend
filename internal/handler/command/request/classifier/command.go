// Package classifier is the gRPC transport for RequestClassifierCommandService.
package classifier

import (
	"strings"

	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/request/classifier"
	requestclassifierv1 "github.com/medincident/medincident-backend/pkg/command/request/classifier/v1"
)

// RequestClassifierHandler implements requestclassifierv1.RequestClassifierCommandServiceServer.
type RequestClassifierHandler struct {
	requestclassifierv1.UnimplementedRequestClassifierCommandServiceServer

	typeSvc *classifiersvc.RequestTypeService
}

func trimOptionalString(s *string) *string {
	if s == nil {
		return nil
	}
	v := strings.TrimSpace(*s)
	return &v
}

// NewRequestClassifierHandler wires the handler with the service.
func NewRequestClassifierHandler(
	typeSvc *classifiersvc.RequestTypeService,
) *RequestClassifierHandler {
	return &RequestClassifierHandler{
		typeSvc: typeSvc,
	}
}
