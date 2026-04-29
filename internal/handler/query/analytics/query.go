// Package analytics is the gRPC transport for AnalyticsQueryService.
package analytics

import (
	"context"
	"database/sql"
	"time"

	"github.com/medincident/medincident-backend/internal/middleware/grpcmw"
	"github.com/medincident/medincident-backend/internal/service/authz"
	analyticsread "github.com/medincident/medincident-backend/internal/service/query/analytics"
	analyticsqueryv1 "github.com/medincident/medincident-backend/pkg/query/analytics/v1"
)

// AnalyticsQueryHandler implements analyticsqueryv1.AnalyticsQueryServiceServer.
type AnalyticsQueryHandler struct {
	analyticsqueryv1.UnimplementedAnalyticsQueryServiceServer

	reader *analyticsread.Reader
}

// NewAnalyticsQueryHandler wires the handler with a Reader.
func NewAnalyticsQueryHandler(reader *analyticsread.Reader) *AnalyticsQueryHandler {
	return &AnalyticsQueryHandler{reader: reader}
}

// GetSnapshot returns raw analytics records for the given filter.
//
// See: docs/services/analytics/Analytics.md
func (h *AnalyticsQueryHandler) GetSnapshot(
	ctx context.Context,
	req *analyticsqueryv1.GetSnapshotRequest,
) (*analyticsqueryv1.GetSnapshotResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}

	result, err := h.reader.GetSnapshot(
		ctx,
		caller,
		req.GetOrganizationId(),
		req.GetFrom(),
		req.GetTo(),
		req.ClinicId,
		req.DepartmentId,
		req.GetIncludePatientBuffer(),
	)
	if err != nil {
		return nil, err
	}

	incidents := make([]*analyticsqueryv1.SnapshotIncident, len(result.Incidents))
	for i := range result.Incidents {
		inc := &result.Incidents[i]
		si := &analyticsqueryv1.SnapshotIncident{
			CreatedAt:           inc.CreatedAt.Format(time.RFC3339),
			OccurredAt:          inc.OccurredAt.Format(time.RFC3339),
			Status:              inc.Status,
			Priority:            inc.Priority,
			CategoryId:          inc.CategoryID.String(),
			CategoryName:        inc.CategoryName,
			TypeId:              inc.TypeID.String(),
			TypeName:            inc.TypeName,
			ClinicId:            inc.ClinicID.String(),
			DepartmentId:        inc.DepartmentID.String(),
			IsPatientSource:     inc.IsPatientSource,
			IsReopened:          inc.IsReopened,
			LinkedRequestsCount: inc.LinkedRequestsCount,
		}
		if inc.ClosedAt != nil {
			s := inc.ClosedAt.Format(time.RFC3339)
			si.ClosedAt = &s
		}
		incidents[i] = si
	}

	requests := make([]*analyticsqueryv1.SnapshotRequest, len(result.Requests))
	for i, r := range result.Requests {
		sr := &analyticsqueryv1.SnapshotRequest{
			CreatedAt:         r.CreatedAt.Format(time.RFC3339),
			Status:            r.Status,
			TypeId:            r.TypeID.String(),
			TypeName:          r.TypeName,
			DepartmentId:      r.DepartmentID.String(),
			HasLinkedIncident: r.HasLinkedIncident,
		}
		if r.CompletedAt != nil {
			s := r.CompletedAt.Format(time.RFC3339)
			sr.CompletedAt = &s
		}
		requests[i] = sr
	}

	buf := make([]*analyticsqueryv1.SnapshotPatientBuffer, len(result.PatientBuffer))
	for i, b := range result.PatientBuffer {
		sb := &analyticsqueryv1.SnapshotPatientBuffer{
			CreatedAt: b.CreatedAt.Format(time.RFC3339),
			Status:    b.Status,
		}
		if b.CategoryID != nil {
			s := b.CategoryID.String()
			sb.CategoryId = &s
		}
		buf[i] = sb
	}

	return &analyticsqueryv1.GetSnapshotResponse{
		Incidents:     incidents,
		Requests:      requests,
		PatientBuffer: buf,
	}, nil
}

// GetSummary returns aggregated KPIs and distributions.
//
// See: docs/services/analytics/Analytics.md
func (h *AnalyticsQueryHandler) GetSummary(
	ctx context.Context,
	req *analyticsqueryv1.GetSummaryRequest,
) (*analyticsqueryv1.GetSummaryResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}

	result, err := h.reader.GetSummary(
		ctx,
		caller,
		req.GetOrganizationId(),
		req.GetFrom(),
		req.GetTo(),
		req.ClinicId,
		req.DepartmentId,
	)
	if err != nil {
		return nil, err
	}

	ia := result.IncidentAgg
	ip := result.IncidentPercentiles

	var incResolution *analyticsqueryv1.ResolutionStats
	if ia.Done > 0 {
		incResolution = &analyticsqueryv1.ResolutionStats{
			AvgMinutes: orZero(ia.ResolutionAvg),
			MinMinutes: orZero(ia.ResolutionMin),
			MaxMinutes: orZero(ia.ResolutionMax),
			P50Minutes: orZero(ip.P50),
			P90Minutes: orZero(ip.P90),
			P95Minutes: orZero(ip.P95),
		}
	}

	incSummary := &analyticsqueryv1.IncidentSummary{
		Total: ia.Total,
		ByStatus: &analyticsqueryv1.IncidentStatusBreakdown{
			Pending:    ia.Pending,
			InProgress: ia.InProgress,
			Done:       ia.Done,
			Rejected:   ia.Rejected,
			Cancelled:  ia.Cancelled,
		},
		ByPriority: &analyticsqueryv1.IncidentPriorityBreakdown{
			Low:      ia.Low,
			Normal:   ia.Normal,
			High:     ia.High,
			Critical: ia.Critical,
		},
		BySource: &analyticsqueryv1.IncidentSourceBreakdown{
			Staff:   ia.StaffSource,
			Patient: ia.PatientSource,
		},
		Reopened:           ia.Reopened,
		WithLinkedRequests: ia.WithLinkedRequests,
		Resolution:         incResolution,
		TopCategories:      toProtoCategories(result.IncidentCategories),
		TopTypes:           toProtoTypes(result.IncidentTypes),
		TopDepartments:     toProtoDepartments(result.IncidentDepartments),
	}

	ra := result.RequestAgg
	rp := result.RequestPercentiles

	var reqCompletion *analyticsqueryv1.ResolutionStats
	if ra.Completed > 0 {
		reqCompletion = &analyticsqueryv1.ResolutionStats{
			AvgMinutes: orZero(ra.CompletionAvg),
			MinMinutes: orZero(ra.CompletionMin),
			MaxMinutes: orZero(ra.CompletionMax),
			P50Minutes: orZero(rp.P50),
			P90Minutes: orZero(rp.P90),
			P95Minutes: orZero(rp.P95),
		}
	}

	reqSummary := &analyticsqueryv1.RequestSummary{
		Total: ra.Total,
		ByStatus: &analyticsqueryv1.RequestStatusBreakdown{
			Created:       ra.Created,
			InWork:        ra.InWork,
			OnHold:        ra.OnHold,
			PendingReview: ra.PendingReview,
			Completed:     ra.Completed,
			Cancelled:     ra.Cancelled,
		},
		Linked:         ra.Linked,
		Unlinked:       ra.Unlinked,
		Completion:     reqCompletion,
		TopTypes:       toProtoTypes(result.RequestTypes),
		TopDepartments: toProtoDepartments(result.RequestDepartments),
	}

	ba := result.BufferAgg
	resolved := ba.Published + ba.Rejected + ba.Cancelled
	var acceptRate, rejectRate float64
	if resolved > 0 {
		acceptRate = float64(ba.Published) / float64(resolved)
		rejectRate = float64(ba.Rejected) / float64(resolved)
	}

	bufSummary := &analyticsqueryv1.PatientBufferSummary{
		Total: ba.Total,
		ByStatus: &analyticsqueryv1.PatientBufferStatusBreakdown{
			Pending:   ba.Pending,
			Published: ba.Published,
			Rejected:  ba.Rejected,
			Cancelled: ba.Cancelled,
		},
		AcceptanceRate: acceptRate,
		RejectionRate:  rejectRate,
	}

	period := &analyticsqueryv1.SummaryPeriod{
		OrganizationId: req.GetOrganizationId(),
		From:           req.GetFrom(),
		To:             req.GetTo(),
		ClinicId:       req.ClinicId,
		DepartmentId:   req.DepartmentId,
	}

	return &analyticsqueryv1.GetSummaryResponse{
		Incidents:     incSummary,
		Requests:      reqSummary,
		PatientBuffer: bufSummary,
		Period:        period,
	}, nil
}

// GetTimeSeries returns bucketed counts for the given filter and granularity.
//
// See: docs/services/analytics/Analytics.md
func (h *AnalyticsQueryHandler) GetTimeSeries(
	ctx context.Context,
	req *analyticsqueryv1.GetTimeSeriesRequest,
) (*analyticsqueryv1.GetTimeSeriesResponse, error) {
	callerID, err := grpcmw.CallerID(ctx)
	if err != nil {
		return nil, err
	}
	caller := authz.Caller{ZitadelUserID: callerID}

	buckets, err := h.reader.GetTimeSeries(
		ctx,
		caller,
		req.GetOrganizationId(),
		req.GetFrom(),
		req.GetTo(),
		req.ClinicId,
		req.DepartmentId,
		protoGranToService(req.GetGranularity()),
	)
	if err != nil {
		return nil, err
	}

	out := make([]*analyticsqueryv1.TimeSeriesBucket, len(buckets))
	for i := range buckets {
		b := &buckets[i]
		out[i] = &analyticsqueryv1.TimeSeriesBucket{
			BucketStart: b.BucketStart.Format(time.RFC3339),
			BucketEnd:   b.BucketEnd.Format(time.RFC3339),
			Incidents: &analyticsqueryv1.TimeSeriesIncidentBucket{
				Total:         b.IncidentTotal,
				Pending:       b.IPending,
				InProgress:    b.IInProgress,
				Done:          b.IDone,
				Rejected:      b.IRejected,
				Cancelled:     b.ICancelled,
				HighCritical:  b.IHighCritical,
				PatientSource: b.IPatientSource,
				Reopened:      b.IReopened,
			},
			Requests: &analyticsqueryv1.TimeSeriesRequestBucket{
				Total:     b.ReqTotal,
				Completed: b.ReqCompleted,
				Cancelled: b.ReqCancelled,
				Linked:    b.ReqLinked,
			},
		}
	}

	return &analyticsqueryv1.GetTimeSeriesResponse{Buckets: out}, nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func orZero(n sql.NullFloat64) float64 {
	if n.Valid {
		return n.Float64
	}
	return 0
}

func protoGranToService(g analyticsqueryv1.TimeSeriesGranularity) analyticsread.TimeSeriesGranularity {
	switch g {
	case analyticsqueryv1.TimeSeriesGranularity_TIME_SERIES_GRANULARITY_WEEK:
		return analyticsread.GranularityWeek
	case analyticsqueryv1.TimeSeriesGranularity_TIME_SERIES_GRANULARITY_MONTH:
		return analyticsread.GranularityMonth
	default:
		return analyticsread.GranularityDay
	}
}

func toProtoCategories(rows []analyticsread.DistributionRow) []*analyticsqueryv1.CategoryCount {
	out := make([]*analyticsqueryv1.CategoryCount, len(rows))
	for i, r := range rows {
		out[i] = &analyticsqueryv1.CategoryCount{
			CategoryId:   r.ID.String(),
			CategoryName: r.Name,
			Count:        r.Count,
		}
	}
	return out
}

func toProtoTypes(rows []analyticsread.DistributionRow) []*analyticsqueryv1.TypeCount {
	out := make([]*analyticsqueryv1.TypeCount, len(rows))
	for i, r := range rows {
		out[i] = &analyticsqueryv1.TypeCount{
			TypeId:   r.ID.String(),
			TypeName: r.Name,
			Count:    r.Count,
		}
	}
	return out
}

func toProtoDepartments(rows []analyticsread.DistributionRow) []*analyticsqueryv1.DepartmentCount {
	out := make([]*analyticsqueryv1.DepartmentCount, len(rows))
	for i, r := range rows {
		out[i] = &analyticsqueryv1.DepartmentCount{
			DepartmentId:   r.ID.String(),
			DepartmentName: r.Name,
			Count:          r.Count,
		}
	}
	return out
}
