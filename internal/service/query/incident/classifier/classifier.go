package classifier

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"

	"github.com/medincident/medincident-backend/internal/cursor"
	"github.com/medincident/medincident-backend/internal/service/authz"
)

// Error codes emitted by the classifier reader.
const (
	ErrCodeCategoryNotFound   = "incident_category_not_found"
	ErrCodeCategoryLoadFailed = "incident_category_load_failed"
	ErrCodeTypeNotFound       = "incident_type_not_found"
	ErrCodeTypeLoadFailed     = "incident_type_load_failed"
)

// CategoryView mirrors projections.incident_categories.
type CategoryView struct {
	ID               uuid.UUID
	OrganizationID   uuid.UUID
	ParentCategoryID *uuid.UUID
	Name             string
	Description      *string
	IsActive         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// TypeView mirrors projections.incident_types.
type TypeView struct {
	ID                   uuid.UUID
	OrganizationID       uuid.UUID
	CategoryID           uuid.UUID
	Name                 string
	Description          *string
	IsActive             bool
	IsAllowedForPatients bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

const selectCategory = `
	SELECT id, organization_id, parent_category_id, name, description,
	       is_active, created_at, updated_at
	  FROM projections.incident_categories`

const selectType = `
	SELECT id, organization_id, category_id, name, description,
	       is_active, is_allowed_for_patients, created_at, updated_at
	  FROM projections.incident_types`

// scanCategory scans one row into a CategoryView.
func scanCategory(scanner interface {
	Scan(dest ...any) error
}, out *CategoryView,
) error {
	return scanner.Scan(
		&out.ID, &out.OrganizationID, &out.ParentCategoryID, &out.Name, &out.Description,
		&out.IsActive, &out.CreatedAt, &out.UpdatedAt,
	)
}

// scanType scans one row into a TypeView.
func scanType(scanner interface {
	Scan(dest ...any) error
}, out *TypeView,
) error {
	return scanner.Scan(
		&out.ID, &out.OrganizationID, &out.CategoryID, &out.Name, &out.Description,
		&out.IsActive, &out.IsAllowedForPatients, &out.CreatedAt, &out.UpdatedAt,
	)
}

// GetCategory returns one incident category by id. Authorization:
// authz.ReaderOf.Category(id).
//
// See: docs/services/incident/Classifier.md
func (r *Reader) GetCategory(
	ctx context.Context,
	caller authz.Caller,
	id uuid.UUID,
) (*CategoryView, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.Category(id)); err != nil {
		return nil, err
	}
	var out CategoryView
	err := scanCategory(r.db.WithContext(ctx).Raw(selectCategory+` WHERE id = ?`, id).Row(), &out)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In("reader.incident.classifier.category").
				Code(ErrCodeCategoryNotFound).
				Public("Incident category not found.").
				With("category_id", id).
				Errorf("not found")
		}
		return nil, oops.In("reader.incident.classifier.category").
			Code(ErrCodeCategoryLoadFailed).
			With("category_id", id).
			Wrap(err)
	}
	return &out, nil
}

// ListCategoriesByOrganization paginates categories for one org.
// Employees (ReaderOf.Organization) receive all categories.
// Authenticated non-employees (patients) receive only active categories
// whose subtree contains at least one active, patient-allowed type.
//
// See: docs/services/incident/Classifier.md
func (r *Reader) ListCategoriesByOrganization(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
	q ListQuery,
) (CategoryListResult, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.Authenticated); err != nil {
		return CategoryListResult{}, err
	}
	if err := q.normalize(); err != nil {
		return CategoryListResult{}, err
	}
	isEmployee, err := r.authz.Satisfies(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID))
	if err != nil {
		return CategoryListResult{}, err
	}
	if isEmployee {
		return r.listCategoriesByOrgEmployee(ctx, orgID, q)
	}
	return r.listPatientVisibleCategories(ctx, orgID, q, false)
}

func (r *Reader) listCategoriesByOrgEmployee(
	ctx context.Context,
	orgID uuid.UUID,
	q ListQuery,
) (CategoryListResult, error) {
	sqlBuf := selectCategory + ` WHERE organization_id = ?`
	args := make([]any, 0, 4)
	args = append(args, orgID)
	if q.After != nil {
		c, err := cursor.Decode(*q.After)
		if err != nil {
			return CategoryListResult{}, oops.In("reader.incident.classifier.category").
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		sqlBuf += ` AND (updated_at, id) < (?, ?)`
		args = append(args, c.Time(), c.I)
	}
	sqlBuf += ` ORDER BY updated_at DESC, id DESC LIMIT ?`
	args = append(args, q.Limit+1)
	rows, err := r.db.WithContext(ctx).Raw(sqlBuf, args...).Rows()
	if err != nil {
		return CategoryListResult{}, oops.In("reader.incident.classifier.category").
			Code(ErrCodeCategoryLoadFailed).
			With("organization_id", orgID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]CategoryView, 0, q.Limit+1)
	for rows.Next() {
		var v CategoryView
		if err := scanCategory(rows, &v); err != nil {
			return CategoryListResult{}, oops.In("reader.incident.classifier.category").
				Code(ErrCodeCategoryLoadFailed).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return CategoryListResult{}, oops.In("reader.incident.classifier.category").
			Code(ErrCodeCategoryLoadFailed).
			Wrap(err)
	}
	var nextCursor *string
	if len(out) > q.Limit {
		out = out[:q.Limit]
		last := out[len(out)-1]
		s := cursor.Encode(last.UpdatedAt, last.ID.String())
		nextCursor = &s
	}
	return CategoryListResult{Items: out, NextCursor: nextCursor}, nil
}

// listPatientVisibleCategories returns categories the patient may see: active
// categories whose subtree (following active ancestors up) contains at least
// one active, patient-allowed type. When rootOnly is true only top-level
// categories (parent_category_id IS NULL) are returned.
func (r *Reader) listPatientVisibleCategories(
	ctx context.Context,
	orgID uuid.UUID,
	q ListQuery,
	rootOnly bool,
) (CategoryListResult, error) {
	cursorClause := ""
	args := []any{orgID, orgID, orgID, orgID}
	if q.After != nil {
		c, err := cursor.Decode(*q.After)
		if err != nil {
			return CategoryListResult{}, oops.In("reader.incident.classifier.category").
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		cursorClause = ` AND (updated_at, id) < (?, ?)`
		args = append(args, c.Time(), c.I)
	}
	args = append(args, q.Limit+1)

	rootClause := ""
	if rootOnly {
		rootClause = ` AND parent_category_id IS NULL`
	}

	query := `
		WITH RECURSIVE
		  direct AS (
		    SELECT DISTINCT c.id, c.parent_category_id
		      FROM projections.incident_categories c
		      JOIN projections.incident_types t ON t.category_id = c.id
		     WHERE c.organization_id = ?
		       AND t.organization_id = ?
		       AND c.is_active = TRUE
		       AND t.is_active = TRUE
		       AND t.is_allowed_for_patients = TRUE
		  ),
		  visible(id, parent_category_id) AS (
		    SELECT id, parent_category_id FROM direct
		    UNION
		    SELECT c.id, c.parent_category_id
		      FROM projections.incident_categories c
		      JOIN visible v ON v.parent_category_id = c.id
		     WHERE c.organization_id = ?
		       AND c.is_active = TRUE
		  )
		SELECT id, organization_id, parent_category_id, name, description,
		       is_active, created_at, updated_at
		  FROM projections.incident_categories
		 WHERE id IN (SELECT id FROM visible)
		   AND organization_id = ?` +
		rootClause +
		cursorClause +
		` ORDER BY updated_at DESC, id DESC
		 LIMIT ?`

	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return CategoryListResult{}, oops.In("reader.incident.classifier.category").
			Code(ErrCodeCategoryLoadFailed).
			With("organization_id", orgID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]CategoryView, 0, q.Limit+1)
	for rows.Next() {
		var v CategoryView
		if err := scanCategory(rows, &v); err != nil {
			return CategoryListResult{}, oops.In("reader.incident.classifier.category").
				Code(ErrCodeCategoryLoadFailed).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return CategoryListResult{}, oops.In("reader.incident.classifier.category").
			Code(ErrCodeCategoryLoadFailed).
			Wrap(err)
	}
	var nextCursor *string
	if len(out) > q.Limit {
		out = out[:q.Limit]
		last := out[len(out)-1]
		s := cursor.Encode(last.UpdatedAt, last.ID.String())
		nextCursor = &s
	}
	return CategoryListResult{Items: out, NextCursor: nextCursor}, nil
}

// ListActiveRootCategories returns top-level active categories for an org.
// Employees receive all active roots. Patients receive only roots whose
// subtree contains at least one active, patient-allowed type.
//
// See: docs/services/incident/Classifier.md
func (r *Reader) ListActiveRootCategories(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
	q ListQuery,
) (CategoryListResult, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.Authenticated); err != nil {
		return CategoryListResult{}, err
	}
	if err := q.normalize(); err != nil {
		return CategoryListResult{}, err
	}
	isEmployee, err := r.authz.Satisfies(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID))
	if err != nil {
		return CategoryListResult{}, err
	}
	if isEmployee {
		return r.listActiveRootCategoriesEmployee(ctx, orgID, q)
	}
	return r.listPatientVisibleCategories(ctx, orgID, q, true)
}

func (r *Reader) listActiveRootCategoriesEmployee(
	ctx context.Context,
	orgID uuid.UUID,
	q ListQuery,
) (CategoryListResult, error) {
	sqlBuf := selectCategory + `
		 WHERE organization_id = ?
		   AND parent_category_id IS NULL
		   AND is_active = TRUE`
	args := make([]any, 0, 4)
	args = append(args, orgID)
	if q.After != nil {
		c, err := cursor.Decode(*q.After)
		if err != nil {
			return CategoryListResult{}, oops.In("reader.incident.classifier.category").
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		sqlBuf += ` AND (updated_at, id) < (?, ?)`
		args = append(args, c.Time(), c.I)
	}
	sqlBuf += ` ORDER BY updated_at DESC, id DESC LIMIT ?`
	args = append(args, q.Limit+1)
	rows, err := r.db.WithContext(ctx).Raw(sqlBuf, args...).Rows()
	if err != nil {
		return CategoryListResult{}, oops.In("reader.incident.classifier.category").
			Code(ErrCodeCategoryLoadFailed).
			With("organization_id", orgID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]CategoryView, 0, q.Limit+1)
	for rows.Next() {
		var v CategoryView
		if err := scanCategory(rows, &v); err != nil {
			return CategoryListResult{}, oops.In("reader.incident.classifier.category").
				Code(ErrCodeCategoryLoadFailed).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return CategoryListResult{}, oops.In("reader.incident.classifier.category").
			Code(ErrCodeCategoryLoadFailed).
			Wrap(err)
	}
	var nextCursor *string
	if len(out) > q.Limit {
		out = out[:q.Limit]
		last := out[len(out)-1]
		s := cursor.Encode(last.UpdatedAt, last.ID.String())
		nextCursor = &s
	}
	return CategoryListResult{Items: out, NextCursor: nextCursor}, nil
}

// ListCategorySubtree returns every descendant of rootID (inclusive).
// Employees receive the full subtree regardless of active/patient status.
// Patients receive only active categories in the subtree whose descendant
// subtree contains at least one active, patient-allowed type.
//
// See: docs/services/incident/Classifier.md
func (r *Reader) ListCategorySubtree(
	ctx context.Context,
	caller authz.Caller,
	rootID uuid.UUID,
) ([]CategoryView, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.Authenticated); err != nil {
		return nil, err
	}
	isEmployee, err := r.authz.Satisfies(ctx, caller.ZitadelUserID, authz.ReaderOf.Category(rootID))
	if err != nil {
		return nil, err
	}
	if isEmployee {
		return r.listCategorySubtreeEmployee(ctx, rootID)
	}
	return r.listCategorySubtreePatient(ctx, rootID)
}

func (r *Reader) listCategorySubtreeEmployee(ctx context.Context, rootID uuid.UUID) ([]CategoryView, error) {
	const query = `
		WITH RECURSIVE tree AS (
		  SELECT id, organization_id, parent_category_id, name, description,
		         is_active, created_at, updated_at
		    FROM projections.incident_categories
		   WHERE id = ?
		  UNION ALL
		  SELECT c.id, c.organization_id, c.parent_category_id, c.name, c.description,
		         c.is_active, c.created_at, c.updated_at
		    FROM projections.incident_categories c
		    JOIN tree t ON c.parent_category_id = t.id
		)
		SELECT id, organization_id, parent_category_id, name, description,
		       is_active, created_at, updated_at
		  FROM tree
		 ORDER BY created_at ASC, id ASC`
	rows, err := r.db.WithContext(ctx).Raw(query, rootID).Rows()
	if err != nil {
		return nil, oops.In("reader.incident.classifier.category").
			Code(ErrCodeCategoryLoadFailed).
			With("root_category_id", rootID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]CategoryView, 0)
	for rows.Next() {
		var v CategoryView
		if err := scanCategory(rows, &v); err != nil {
			return nil, oops.In("reader.incident.classifier.category").
				Code(ErrCodeCategoryLoadFailed).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, oops.In("reader.incident.classifier.category").
			Code(ErrCodeCategoryLoadFailed).
			Wrap(err)
	}
	return out, nil
}

// listCategorySubtreePatient walks down from rootID collecting only active
// categories that have at least one active patient-allowed type in their
// own descendant subtree. The CTE walks DOWN first (subtree), then marks
// directly-qualifying categories (direct), then walks UP within the subtree
// to include their ancestors (visible).
func (r *Reader) listCategorySubtreePatient(ctx context.Context, rootID uuid.UUID) ([]CategoryView, error) {
	const query = `
		WITH RECURSIVE
		  subtree AS (
		    SELECT id, parent_category_id
		      FROM projections.incident_categories
		     WHERE id = ?
		    UNION ALL
		    SELECT c.id, c.parent_category_id
		      FROM projections.incident_categories c
		      JOIN subtree s ON c.parent_category_id = s.id
		     WHERE c.is_active = TRUE
		  ),
		  direct AS (
		    SELECT DISTINCT c.id, c.parent_category_id
		      FROM projections.incident_categories c
		      JOIN projections.incident_types t ON t.category_id = c.id
		     WHERE c.id IN (SELECT id FROM subtree)
		       AND c.is_active = TRUE
		       AND t.is_active = TRUE
		       AND t.is_allowed_for_patients = TRUE
		  ),
		  visible(id, parent_category_id) AS (
		    SELECT id, parent_category_id FROM direct
		    UNION
		    SELECT c.id, c.parent_category_id
		      FROM projections.incident_categories c
		      JOIN visible v ON v.parent_category_id = c.id
		     WHERE c.id IN (SELECT id FROM subtree)
		       AND c.is_active = TRUE
		  )
		SELECT id, organization_id, parent_category_id, name, description,
		       is_active, created_at, updated_at
		  FROM projections.incident_categories
		 WHERE id IN (SELECT id FROM visible)
		 ORDER BY created_at ASC, id ASC`
	rows, err := r.db.WithContext(ctx).Raw(query, rootID).Rows()
	if err != nil {
		return nil, oops.In("reader.incident.classifier.category").
			Code(ErrCodeCategoryLoadFailed).
			With("root_category_id", rootID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]CategoryView, 0)
	for rows.Next() {
		var v CategoryView
		if err := scanCategory(rows, &v); err != nil {
			return nil, oops.In("reader.incident.classifier.category").
				Code(ErrCodeCategoryLoadFailed).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, oops.In("reader.incident.classifier.category").
			Code(ErrCodeCategoryLoadFailed).
			Wrap(err)
	}
	return out, nil
}

// GetType returns one incident type by id. Authorization:
// authz.ReaderOf.IncidentType(id).
//
// See: docs/services/incident/Classifier.md
func (r *Reader) GetType(
	ctx context.Context,
	caller authz.Caller,
	id uuid.UUID,
) (*TypeView, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.ReaderOf.IncidentType(id)); err != nil {
		return nil, err
	}
	var out TypeView
	err := scanType(r.db.WithContext(ctx).Raw(selectType+` WHERE id = ?`, id).Row(), &out)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In("reader.incident.classifier.type").
				Code(ErrCodeTypeNotFound).
				Public("Incident type not found.").
				With("type_id", id).
				Errorf("not found")
		}
		return nil, oops.In("reader.incident.classifier.type").
			Code(ErrCodeTypeLoadFailed).
			With("type_id", id).
			Wrap(err)
	}
	return &out, nil
}

// ListTypesByCategory returns every type under one category.
// Authorization: authz.ReaderOf.Category(categoryID).
//
// See: docs/services/incident/Classifier.md
// ListTypesByCategory paginates incident types under a category.
// Employees (ReaderOf.Category) receive all types.
// Authenticated non-employees (patients) receive only active,
// patient-allowed types.
//
// See: docs/services/incident/Classifier.md
func (r *Reader) ListTypesByCategory(
	ctx context.Context,
	caller authz.Caller,
	categoryID uuid.UUID,
	q ListQuery,
) (TypeListResult, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.Authenticated); err != nil {
		return TypeListResult{}, err
	}
	if err := q.normalize(); err != nil {
		return TypeListResult{}, err
	}
	isEmployee, err := r.authz.Satisfies(ctx, caller.ZitadelUserID, authz.ReaderOf.Category(categoryID))
	if err != nil {
		return TypeListResult{}, err
	}
	if isEmployee {
		return r.listTypesByCategoryEmployee(ctx, categoryID, q)
	}
	return r.listTypesByCategoryPatient(ctx, categoryID, q)
}

func (r *Reader) listTypesByCategoryEmployee(
	ctx context.Context,
	categoryID uuid.UUID,
	q ListQuery,
) (TypeListResult, error) {
	sqlBuf := selectType + ` WHERE category_id = ?`
	args := make([]any, 0, 4)
	args = append(args, categoryID)
	if q.After != nil {
		c, err := cursor.Decode(*q.After)
		if err != nil {
			return TypeListResult{}, oops.In("reader.incident.classifier.type").
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		sqlBuf += ` AND (updated_at, id) < (?, ?)`
		args = append(args, c.Time(), c.I)
	}
	sqlBuf += ` ORDER BY updated_at DESC, id DESC LIMIT ?`
	args = append(args, q.Limit+1)
	rows, err := r.db.WithContext(ctx).Raw(sqlBuf, args...).Rows()
	if err != nil {
		return TypeListResult{}, oops.In("reader.incident.classifier.type").
			Code(ErrCodeTypeLoadFailed).
			With("category_id", categoryID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]TypeView, 0, q.Limit+1)
	for rows.Next() {
		var v TypeView
		if err := scanType(rows, &v); err != nil {
			return TypeListResult{}, oops.In("reader.incident.classifier.type").
				Code(ErrCodeTypeLoadFailed).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return TypeListResult{}, oops.In("reader.incident.classifier.type").
			Code(ErrCodeTypeLoadFailed).
			Wrap(err)
	}
	var nextCursor *string
	if len(out) > q.Limit {
		out = out[:q.Limit]
		last := out[len(out)-1]
		s := cursor.Encode(last.UpdatedAt, last.ID.String())
		nextCursor = &s
	}
	return TypeListResult{Items: out, NextCursor: nextCursor}, nil
}

func (r *Reader) listTypesByCategoryPatient(
	ctx context.Context,
	categoryID uuid.UUID,
	q ListQuery,
) (TypeListResult, error) {
	sqlBuf := selectType + ` WHERE category_id = ? AND is_active = TRUE AND is_allowed_for_patients = TRUE`
	args := make([]any, 0, 4)
	args = append(args, categoryID)
	if q.After != nil {
		c, err := cursor.Decode(*q.After)
		if err != nil {
			return TypeListResult{}, oops.In("reader.incident.classifier.type").
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		sqlBuf += ` AND (updated_at, id) < (?, ?)`
		args = append(args, c.Time(), c.I)
	}
	sqlBuf += ` ORDER BY updated_at DESC, id DESC LIMIT ?`
	args = append(args, q.Limit+1)
	rows, err := r.db.WithContext(ctx).Raw(sqlBuf, args...).Rows()
	if err != nil {
		return TypeListResult{}, oops.In("reader.incident.classifier.type").
			Code(ErrCodeTypeLoadFailed).
			With("category_id", categoryID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]TypeView, 0, q.Limit+1)
	for rows.Next() {
		var v TypeView
		if err := scanType(rows, &v); err != nil {
			return TypeListResult{}, oops.In("reader.incident.classifier.type").
				Code(ErrCodeTypeLoadFailed).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return TypeListResult{}, oops.In("reader.incident.classifier.type").
			Code(ErrCodeTypeLoadFailed).
			Wrap(err)
	}
	var nextCursor *string
	if len(out) > q.Limit {
		out = out[:q.Limit]
		last := out[len(out)-1]
		s := cursor.Encode(last.UpdatedAt, last.ID.String())
		nextCursor = &s
	}
	return TypeListResult{Items: out, NextCursor: nextCursor}, nil
}

// ListActiveTypesByOrganization returns active incident types for one org.
// Employees (ReaderOf.Organization) receive all active types.
// Authenticated non-employees (patients) receive only active,
// patient-allowed types.
//
// See: docs/services/incident/Classifier.md
func (r *Reader) ListActiveTypesByOrganization(
	ctx context.Context,
	caller authz.Caller,
	orgID uuid.UUID,
	q ListQuery,
) (TypeListResult, error) {
	if err := r.authz.Require(ctx, caller.ZitadelUserID, authz.Authenticated); err != nil {
		return TypeListResult{}, err
	}
	if err := q.normalize(); err != nil {
		return TypeListResult{}, err
	}
	isEmployee, err := r.authz.Satisfies(ctx, caller.ZitadelUserID, authz.ReaderOf.Organization(orgID))
	if err != nil {
		return TypeListResult{}, err
	}
	if isEmployee {
		return r.listActiveTypesByOrgEmployee(ctx, orgID, q)
	}
	return r.listActiveTypesByOrgPatient(ctx, orgID, q)
}

func (r *Reader) listActiveTypesByOrgEmployee(
	ctx context.Context,
	orgID uuid.UUID,
	q ListQuery,
) (TypeListResult, error) {
	sqlBuf := selectType + ` WHERE organization_id = ? AND is_active = TRUE`
	args := make([]any, 0, 4)
	args = append(args, orgID)
	if q.After != nil {
		c, err := cursor.Decode(*q.After)
		if err != nil {
			return TypeListResult{}, oops.In("reader.incident.classifier.type").
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		sqlBuf += ` AND (updated_at, id) < (?, ?)`
		args = append(args, c.Time(), c.I)
	}
	sqlBuf += ` ORDER BY updated_at DESC, id DESC LIMIT ?`
	args = append(args, q.Limit+1)
	rows, err := r.db.WithContext(ctx).Raw(sqlBuf, args...).Rows()
	if err != nil {
		return TypeListResult{}, oops.In("reader.incident.classifier.type").
			Code(ErrCodeTypeLoadFailed).
			With("organization_id", orgID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]TypeView, 0, q.Limit+1)
	for rows.Next() {
		var v TypeView
		if err := scanType(rows, &v); err != nil {
			return TypeListResult{}, oops.In("reader.incident.classifier.type").
				Code(ErrCodeTypeLoadFailed).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return TypeListResult{}, oops.In("reader.incident.classifier.type").
			Code(ErrCodeTypeLoadFailed).
			Wrap(err)
	}
	var nextCursor *string
	if len(out) > q.Limit {
		out = out[:q.Limit]
		last := out[len(out)-1]
		s := cursor.Encode(last.UpdatedAt, last.ID.String())
		nextCursor = &s
	}
	return TypeListResult{Items: out, NextCursor: nextCursor}, nil
}

func (r *Reader) listActiveTypesByOrgPatient(
	ctx context.Context,
	orgID uuid.UUID,
	q ListQuery,
) (TypeListResult, error) {
	sqlBuf := selectType + ` WHERE organization_id = ? AND is_active = TRUE AND is_allowed_for_patients = TRUE`
	args := make([]any, 0, 4)
	args = append(args, orgID)
	if q.After != nil {
		c, err := cursor.Decode(*q.After)
		if err != nil {
			return TypeListResult{}, oops.In("reader.incident.classifier.type").
				Code(ErrCodeListBadCursor).
				Public("Invalid pagination cursor.").
				Wrap(err)
		}
		sqlBuf += ` AND (updated_at, id) < (?, ?)`
		args = append(args, c.Time(), c.I)
	}
	sqlBuf += ` ORDER BY updated_at DESC, id DESC LIMIT ?`
	args = append(args, q.Limit+1)
	rows, err := r.db.WithContext(ctx).Raw(sqlBuf, args...).Rows()
	if err != nil {
		return TypeListResult{}, oops.In("reader.incident.classifier.type").
			Code(ErrCodeTypeLoadFailed).
			With("organization_id", orgID).
			Wrap(err)
	}
	defer func() { _ = rows.Close() }()
	out := make([]TypeView, 0, q.Limit+1)
	for rows.Next() {
		var v TypeView
		if err := scanType(rows, &v); err != nil {
			return TypeListResult{}, oops.In("reader.incident.classifier.type").
				Code(ErrCodeTypeLoadFailed).
				Wrap(err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return TypeListResult{}, oops.In("reader.incident.classifier.type").
			Code(ErrCodeTypeLoadFailed).
			Wrap(err)
	}
	var nextCursor *string
	if len(out) > q.Limit {
		out = out[:q.Limit]
		last := out[len(out)-1]
		s := cursor.Encode(last.UpdatedAt, last.ID.String())
		nextCursor = &s
	}
	return TypeListResult{Items: out, NextCursor: nextCursor}, nil
}
