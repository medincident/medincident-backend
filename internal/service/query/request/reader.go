package request

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/authz"
)

// ServiceRequestView is the projection row returned to handlers.
type ServiceRequestView struct {
	ID                uuid.UUID
	OrganizationID    uuid.UUID
	ClinicID          uuid.UUID
	DepartmentID      uuid.UUID
	TypeID            uuid.UUID
	IncidentID        *uuid.UUID
	Description       string
	Status            model.ServiceRequestStatus
	AuthorID          string
	AuthorDisplayName string
	Executors         []ExecutorView
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ExecutorView mirrors domain.service_request_executors joined to projections.
type ExecutorView struct {
	EmployeeID   uuid.UUID
	AssignedAt   time.Time
	AssignedByID string
}

// StatusHistoryView mirrors projections.service_request_status_history.
type StatusHistoryView struct {
	ID        uuid.UUID
	OldStatus *string
	NewStatus string
	ActorID   string
	ActorName string
	ChangedAt time.Time
}

// ExecutorHistoryView mirrors projections.service_request_executor_history.
type ExecutorHistoryView struct {
	ID           uuid.UUID
	Action       string
	EmployeeID   uuid.UUID
	EmployeeName string
	ActorID      string
	ActorName    string
	ChangedAt    time.Time
}

// ServiceRequestHistory bundles both timelines.
type ServiceRequestHistory struct {
	StatusHistory   []StatusHistoryView
	ExecutorHistory []ExecutorHistoryView
}

const selectServiceRequest = `
	SELECT id, organization_id, clinic_id, department_id, type_id, incident_id,
	       description, status, author_id, author_display_name,
	       created_at, updated_at
	  FROM projections.service_requests`

func scanServiceRequest(scanner interface{ Scan(...any) error }, v *ServiceRequestView) error {
	return scanner.Scan(
		&v.ID, &v.OrganizationID, &v.ClinicID, &v.DepartmentID,
		&v.TypeID, &v.IncidentID, &v.Description, &v.Status,
		&v.AuthorID, &v.AuthorDisplayName, &v.CreatedAt, &v.UpdatedAt,
	)
}

func (r *Reader) loadExecutors(ctx context.Context, requestID uuid.UUID) ([]ExecutorView, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT employee_id, assigned_at, assigned_by_id
		FROM domain.service_request_executors
		WHERE request_id = ?
		ORDER BY assigned_at ASC`, requestID,
	).Rows()
	if err != nil {
		return nil, wrapRead(err, "load executors")
	}
	defer func() { _ = rows.Close() }()
	var out []ExecutorView
	for rows.Next() {
		var v ExecutorView
		if err := rows.Scan(&v.EmployeeID, &v.AssignedAt, &v.AssignedByID); err != nil {
			return nil, wrapRead(err, "scan executor")
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, wrapRead(err, "iterate executors")
	}
	return out, nil
}

// GetServiceRequest loads one service request with executors.
// Authorization: authz.ReaderOf.Organization(orgID) — the org is
// resolved from the projection row.
//
// See: docs/services/request/Requests.md
func (r *Reader) GetServiceRequest(
	ctx context.Context,
	caller authz.Caller,
	id uuid.UUID,
) (*ServiceRequestView, error) {
	var v ServiceRequestView
	err := scanServiceRequest(
		r.db.WithContext(ctx).Raw(selectServiceRequest+` WHERE id = ?`, id).Row(), &v)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In(scope).Code(ErrCodeServiceRequestNotFound).
				Public("Service request not found.").With("service_request_id", id).
				Errorf("not found")
		}
		return nil, wrapRead(err, "scan service request")
	}
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(v.OrganizationID)); err != nil {
		return nil, err
	}
	execs, err := r.loadExecutors(ctx, v.ID)
	if err != nil {
		return nil, err
	}
	v.Executors = execs
	return &v, nil
}

// ListServiceRequests returns service requests for an organization.
// Authorization: authz.ReaderOf.Organization(orgID).
//
// See: docs/services/request/Requests.md
func (r *Reader) ListServiceRequests(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
	q ListQuery,
) (ServiceRequestListResult, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID)); err != nil {
		return ServiceRequestListResult{}, err
	}
	if err := q.normalize(); err != nil {
		return ServiceRequestListResult{}, err
	}
	sqlBuf := selectServiceRequest + ` WHERE organization_id = ?`
	args := make([]any, 0, 4)
	args = append(args, orgID)
	if q.After != nil {
		c, err := decodeCursor[createdAtCursor](*q.After)
		if err != nil {
			return ServiceRequestListResult{}, oops.In(scope).
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		sqlBuf += ` AND (created_at, id) < (?, ?)`
		args = append(args, c.CreatedAt, c.ID)
	}
	sqlBuf += ` ORDER BY created_at DESC, id DESC LIMIT ?`
	args = append(args, q.Limit+1)
	rows, err := r.db.WithContext(ctx).Raw(sqlBuf, args...).Rows()
	if err != nil {
		return ServiceRequestListResult{}, wrapRead(err, "list service requests")
	}
	defer func() { _ = rows.Close() }()
	out := make([]ServiceRequestView, 0, q.Limit+1)
	for rows.Next() {
		var v ServiceRequestView
		if err := scanServiceRequest(rows, &v); err != nil {
			return ServiceRequestListResult{}, wrapRead(err, "scan service request row")
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return ServiceRequestListResult{}, wrapRead(err, "iterate service request rows")
	}
	var nextCursor *string
	if len(out) > q.Limit {
		out = out[:q.Limit]
		last := out[len(out)-1]
		s := encodeCursor(createdAtCursor{CreatedAt: last.CreatedAt, ID: last.ID})
		nextCursor = &s
	}
	return ServiceRequestListResult{Items: out, NextCursor: nextCursor}, nil
}

// ListServiceRequestsByIncident returns service requests linked to an incident.
// Authorization: resolves the incident's org, then authz.ReaderOf.Organization.
//
// See: docs/services/request/Requests.md
func (r *Reader) ListServiceRequestsByIncident(
	ctx context.Context,
	caller authz.Caller,
	incidentID uuid.UUID,
	q ListQuery,
) (ServiceRequestListResult, error) {
	if err := q.normalize(); err != nil {
		return ServiceRequestListResult{}, err
	}
	var orgID uuid.UUID
	err := r.db.WithContext(ctx).Raw(
		`SELECT organization_id FROM projections.incidents WHERE id = ?`,
		incidentID,
	).Row().Scan(&orgID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ServiceRequestListResult{}, oops.In(scope).
				Code(ErrCodeIncidentNotFound).
				Public("Incident not found.").
				With("incident_id", incidentID).
				Errorf("incident not found")
		}
		return ServiceRequestListResult{}, wrapRead(err, "resolve incident org")
	}
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID)); err != nil {
		return ServiceRequestListResult{}, err
	}

	sqlBuf := selectServiceRequest + ` WHERE incident_id = ?`
	args := make([]any, 0, 4)
	args = append(args, incidentID)
	if q.After != nil {
		c, err := decodeCursor[createdAtCursor](*q.After)
		if err != nil {
			return ServiceRequestListResult{}, oops.In(scope).
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		sqlBuf += ` AND (created_at, id) < (?, ?)`
		args = append(args, c.CreatedAt, c.ID)
	}
	sqlBuf += ` ORDER BY created_at DESC, id DESC LIMIT ?`
	args = append(args, q.Limit+1)
	rows, err := r.db.WithContext(ctx).Raw(sqlBuf, args...).Rows()
	if err != nil {
		return ServiceRequestListResult{}, wrapRead(err, "list service requests by incident")
	}
	defer func() { _ = rows.Close() }()
	out := make([]ServiceRequestView, 0, q.Limit+1)
	for rows.Next() {
		var v ServiceRequestView
		if err := scanServiceRequest(rows, &v); err != nil {
			return ServiceRequestListResult{}, wrapRead(err, "scan service request row")
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return ServiceRequestListResult{}, wrapRead(err, "iterate service request rows")
	}
	var nextCursor *string
	if len(out) > q.Limit {
		out = out[:q.Limit]
		last := out[len(out)-1]
		s := encodeCursor(createdAtCursor{CreatedAt: last.CreatedAt, ID: last.ID})
		nextCursor = &s
	}
	return ServiceRequestListResult{Items: out, NextCursor: nextCursor}, nil
}

// GetServiceRequestHistory returns both timelines for one service request.
//
// See: docs/services/request/Requests.md
func (r *Reader) GetServiceRequestHistory(
	ctx context.Context,
	caller authz.Caller,
	serviceRequestID uuid.UUID,
) (*ServiceRequestHistory, error) {
	// Verify the caller can see this request.
	var orgID uuid.UUID
	err := r.db.WithContext(ctx).Raw(
		`SELECT organization_id FROM projections.service_requests WHERE id = ?`,
		serviceRequestID,
	).Row().Scan(&orgID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In(scope).Code(ErrCodeServiceRequestNotFound).
				Public("Service request not found.").With("service_request_id", serviceRequestID).
				Errorf("not found")
		}
		return nil, wrapRead(err, "resolve request org")
	}
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID)); err != nil {
		return nil, err
	}

	tx := r.db.WithContext(ctx)

	statusRows, err := tx.Raw(
		`SELECT id, old_status, new_status, actor_id, actor_name, changed_at
		 FROM projections.service_request_status_history
		 WHERE request_id = ? ORDER BY changed_at ASC`, serviceRequestID,
	).Rows()
	if err != nil {
		return nil, wrapRead(err, "list status history")
	}
	defer func() { _ = statusRows.Close() }()
	statuses := []StatusHistoryView{}
	for statusRows.Next() {
		var e StatusHistoryView
		if err := statusRows.Scan(&e.ID, &e.OldStatus, &e.NewStatus,
			&e.ActorID, &e.ActorName, &e.ChangedAt); err != nil {
			return nil, wrapRead(err, "scan status history row")
		}
		statuses = append(statuses, e)
	}
	if err := statusRows.Err(); err != nil {
		return nil, wrapRead(err, "iterate status history")
	}

	execRows, err := tx.Raw(
		`SELECT id, action, employee_id, employee_name, actor_id, actor_name, changed_at
		 FROM projections.service_request_executor_history
		 WHERE request_id = ? ORDER BY changed_at ASC`, serviceRequestID,
	).Rows()
	if err != nil {
		return nil, wrapRead(err, "list executor history")
	}
	defer func() { _ = execRows.Close() }()
	executors := []ExecutorHistoryView{}
	for execRows.Next() {
		var e ExecutorHistoryView
		if err := execRows.Scan(&e.ID, &e.Action, &e.EmployeeID, &e.EmployeeName,
			&e.ActorID, &e.ActorName, &e.ChangedAt); err != nil {
			return nil, wrapRead(err, "scan executor history row")
		}
		executors = append(executors, e)
	}
	if err := execRows.Err(); err != nil {
		return nil, wrapRead(err, "iterate executor history")
	}
	return &ServiceRequestHistory{StatusHistory: statuses, ExecutorHistory: executors}, nil
}
