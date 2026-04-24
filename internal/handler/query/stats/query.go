// Package stats is the gRPC transport for StatsQueryService. Only the
// query side exists — there is no command-side stats mutation surface.
package stats

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/service/authz"
	statsread "github.com/medincident/medincident-backend/internal/service/query/stats"
	statsqueryv1 "github.com/medincident/medincident-backend/pkg/query/stats/v1"
)

// Error codes emitted by query-handler ID parsing.
const (
	ErrCodeHandlerInvalidOrganizationID = "handler_invalid_organization_id"
	ErrCodeHandlerInvalidClinicID       = "handler_invalid_clinic_id"
	ErrCodeHandlerInvalidDepartmentID   = "handler_invalid_department_id"
)

// StatsQueryHandler implements statsqueryv1.StatsQueryServiceServer.
type StatsQueryHandler struct {
	statsqueryv1.UnimplementedStatsQueryServiceServer

	reader *statsread.Reader
}

// NewStatsQueryHandler wires the handler with a Reader.
func NewStatsQueryHandler(reader *statsread.Reader) *StatsQueryHandler {
	return &StatsQueryHandler{reader: reader}
}

// GetOrganizationStats returns aggregate stats for an organization.
func (h *StatsQueryHandler) GetOrganizationStats(
	ctx context.Context,
	req *statsqueryv1.GetOrganizationStatsRequest,
) (*statsqueryv1.GetOrganizationStatsResponse, error) {
	caller, err := authz.CallerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseOrganizationID(req.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	s, err := h.reader.GetOrganizationStats(ctx, caller, id)
	if err != nil {
		return nil, err
	}
	return &statsqueryv1.GetOrganizationStatsResponse{
		Stats: &statsqueryv1.OrganizationStats{
			OrganizationId:      s.OrganizationID.String(),
			EmployeesTotal:      s.EmployeesTotal,
			ClinicsTotal:        s.ClinicsTotal,
			DepartmentsTotal:    s.DepartmentsTotal,
			EmployeesOnVacation: s.EmployeesOnVacation,
			VacationsScheduled:  s.VacationsScheduled,
		},
	}, nil
}

// GetClinicStats returns aggregate stats for a clinic.
func (h *StatsQueryHandler) GetClinicStats(
	ctx context.Context,
	req *statsqueryv1.GetClinicStatsRequest,
) (*statsqueryv1.GetClinicStatsResponse, error) {
	caller, err := authz.CallerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseClinicID(req.GetClinicId())
	if err != nil {
		return nil, err
	}
	s, err := h.reader.GetClinicStats(ctx, caller, id)
	if err != nil {
		return nil, err
	}
	return &statsqueryv1.GetClinicStatsResponse{
		Stats: &statsqueryv1.ClinicStats{
			ClinicId:            s.ClinicID.String(),
			OrganizationId:      s.OrganizationID.String(),
			EmployeesTotal:      s.EmployeesTotal,
			DepartmentsTotal:    s.DepartmentsTotal,
			EmployeesOnVacation: s.EmployeesOnVacation,
		},
	}, nil
}

// GetDepartmentStats returns aggregate stats for a department.
func (h *StatsQueryHandler) GetDepartmentStats(
	ctx context.Context,
	req *statsqueryv1.GetDepartmentStatsRequest,
) (*statsqueryv1.GetDepartmentStatsResponse, error) {
	caller, err := authz.CallerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	id, err := parseDepartmentID(req.GetDepartmentId())
	if err != nil {
		return nil, err
	}
	s, err := h.reader.GetDepartmentStats(ctx, caller, id)
	if err != nil {
		return nil, err
	}
	out := &statsqueryv1.DepartmentStats{
		DepartmentId:        s.DepartmentID.String(),
		OrganizationId:      s.OrganizationID.String(),
		EmployeesTotal:      s.EmployeesTotal,
		EmployeesOnVacation: s.EmployeesOnVacation,
	}
	if s.ClinicID != nil {
		cid := s.ClinicID.String()
		out.ClinicId = &cid
	}
	return &statsqueryv1.GetDepartmentStatsResponse{Stats: out}, nil
}

func parseOrganizationID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.stats").
			Code(ErrCodeHandlerInvalidOrganizationID).
			Public("organization_id is not a valid UUID.").
			Wrap(err)
	}
	return id, nil
}

func parseClinicID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.stats").
			Code(ErrCodeHandlerInvalidClinicID).
			Public("clinic_id is not a valid UUID.").
			Wrap(err)
	}
	return id, nil
}

func parseDepartmentID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, oops.In("handler.query.stats").
			Code(ErrCodeHandlerInvalidDepartmentID).
			Public("department_id is not a valid UUID.").
			Wrap(err)
	}
	return id, nil
}
