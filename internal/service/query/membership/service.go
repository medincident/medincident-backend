// Package membership exposes read methods over the membership
// projections: the denormalised employee_cards view, vacations, and
// the role tables (clinic heads, department responsibles, org admins,
// org dispatchers, org heads, system admins). All queries use raw SQL
// via gorm's Raw(...) helper; views are per-query structs co-located
// with the method that returns them.
package membership

import (
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"gorm.io/gorm"
)

// Bounds applied to List pagination.
const (
	listMinLimit     = 1
	listMaxLimit     = 500
	listDefaultLimit = 50
)

// Error codes emitted by pagination validators.
const (
	ErrCodeListLimitOutOfRange  = "list_limit_out_of_range"
	ErrCodeListOffsetOutOfRange = "list_offset_out_of_range"
)

// ListQuery is the input shared by every List method.
type ListQuery struct {
	Limit  int
	Offset int
}

// normalize validates and normalizes the pagination fields.
func (q *ListQuery) normalize() error {
	if q.Offset < 0 {
		return oops.In("reader.membership").
			Code(ErrCodeListOffsetOutOfRange).
			Public("List offset must be non-negative.").
			With("field", "offset").
			With("actual_value", q.Offset).
			Errorf("offset out of range")
	}
	if q.Limit == 0 {
		q.Limit = listDefaultLimit
		return nil
	}
	if q.Limit < listMinLimit || q.Limit > listMaxLimit {
		return oops.In("reader.membership").
			Code(ErrCodeListLimitOutOfRange).
			Public("List limit is out of range.").
			With("field", "limit").
			With("actual_value", q.Limit).
			With("min_value", listMinLimit).
			With("max_value", listMaxLimit).
			Errorf("limit out of range")
	}
	return nil
}

// EmployeeReader exposes read methods over projections.employee_cards
// and projections.employee_vacations.
type EmployeeReader struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewEmployeeReader returns an EmployeeReader bound to the given db.
func NewEmployeeReader(db *gorm.DB, logger *zerolog.Logger) *EmployeeReader {
	return &EmployeeReader{db: db, logger: logger}
}

// RoleReader exposes read methods over the role tables.
type RoleReader struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewRoleReader returns a RoleReader bound to the given db.
func NewRoleReader(db *gorm.DB, logger *zerolog.Logger) *RoleReader {
	return &RoleReader{db: db, logger: logger}
}
