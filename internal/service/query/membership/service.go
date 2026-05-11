// Package membership exposes read methods over the membership
// projections: the denormalised employee_cards view, vacations, and
// the role tables (clinic heads, department responsibles, org admins,
// org dispatchers, org heads, system admins). All queries use raw SQL
// via gorm's Raw(...) helper; views are per-query structs co-located
// with the method that returns them.
//
// Authorization model: every read is gated by authz. Card and role
// lookups scoped to {org,clinic,dept} use authz.ReaderOf.X; the
// single-employee view resolves its scope through domain.employees
// via authz.ReaderOf.Employee. Vacation history is stricter — only
// system admins, the owning organization's admins, and the employee
// themselves may read it, so the policy is composed inline from
// SystemAdmin + OrgAdminOf.Employee + SelfEmployee. ListSystemAdmins
// is scope-less and requires authz.SystemAdmin.
//
// See: docs/services/Membership.md
package membership

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/service/authz"
	"github.com/medincident/medincident-backend/internal/service/query"
)

// Error codes emitted by pagination validators.
const (
	ErrCodeListLimitOutOfRange = "list_limit_out_of_range"
	ErrCodeListBadCursor       = "list_bad_cursor"
)

// employeeCursor is the keyset pagination token for employee_cards
// lists ordered by (updated_at DESC, employee_id DESC).
type employeeCursor struct {
	UpdatedAt  time.Time `json:"updated_at"`
	EmployeeID uuid.UUID `json:"employee_id"`
}

// vacationCursor is the keyset pagination token for vacation lists
// ordered by (starts_at DESC, id DESC).
type vacationCursor struct {
	StartsAt time.Time `json:"starts_at"`
	ID       uuid.UUID `json:"id"`
}

// roleCursor is the keyset pagination token for role-assignment lists
// ordered by (employee_id ASC).
type roleCursor struct {
	EmployeeID uuid.UUID `json:"employee_id"`
}

// systemAdminCursor is the keyset pagination token for system-admin
// lists ordered by (created_at DESC, zitadel_user_id DESC).
type systemAdminCursor struct {
	CreatedAt     time.Time `json:"created_at"`
	ZitadelUserID string    `json:"zitadel_user_id"`
}

func encodeCursor[T any](c T) string {
	b, _ := json.Marshal(c)
	return base64.StdEncoding.EncodeToString(b)
}

func decodeCursor[T any](s string) (T, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	var zero T
	if err != nil {
		return zero, err
	}
	var c T
	if err := json.Unmarshal(b, &c); err != nil {
		return zero, err
	}
	return c, nil
}

// ListQuery is the input shared by every List method.
type ListQuery struct {
	Limit int
	After *string
}

// normalize validates and normalizes the pagination fields.
func (q *ListQuery) normalize() error {
	if q.Limit == 0 {
		q.Limit = query.DefaultLimit
		return nil
	}
	if q.Limit < query.MinLimit || q.Limit > query.MaxLimit {
		return oops.In("reader.membership").
			Code(ErrCodeListLimitOutOfRange).
			Public("List limit is out of range.").
			With("field", "limit").
			With("actual_value", q.Limit).
			With("min_value", query.MinLimit).
			With("max_value", query.MaxLimit).
			Errorf("limit out of range")
	}
	return nil
}

// EmployeeReader exposes read methods over projections.employee_cards
// and projections.employee_vacations.
type EmployeeReader struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewEmployeeReader returns an EmployeeReader bound to the given db
// and authorization service.
func NewEmployeeReader(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *EmployeeReader {
	return &EmployeeReader{db: db, authz: az, logger: logger}
}

// RoleReader exposes read methods over the role tables.
type RoleReader struct {
	db     *gorm.DB
	authz  *authz.Authz
	logger *zerolog.Logger
}

// NewRoleReader returns a RoleReader bound to the given db and
// authorization service.
func NewRoleReader(db *gorm.DB, az *authz.Authz, logger *zerolog.Logger) *RoleReader {
	return &RoleReader{db: db, authz: az, logger: logger}
}
