package self

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/oops"

	memberread "github.com/medincident/medincident-backend/internal/service/query/membership"
)

const (
	ErrCodeSelfEmploymentNotFound = "self_employment_not_found"
	ErrCodeSelfEmploymentFailed   = "self_employment_failed"
)

// GetMyEmployment returns the caller's employee card in the given
// organization. Returns a NOT_FOUND oops error (ErrCodeSelfEmploymentNotFound)
// when the caller is not an active employee of that organization.
//
// See: docs/services/self/Self.md
func (r *SelfReader) GetMyEmployment(
	ctx context.Context,
	callerID string,
	organizationID uuid.UUID,
) (*memberread.EmployeeCardView, error) {
	var out memberread.EmployeeCardView
	err := memberread.ScanEmployeeCard(
		r.db.WithContext(ctx).Raw(
			memberread.SelectEmployeeCard+`
			 WHERE zitadel_user_id = ?
			   AND organization_id = ?
			   AND terminated_at IS NULL`,
			callerID, organizationID,
		).Row(),
		&out,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.In("reader.self").
				Code(ErrCodeSelfEmploymentNotFound).
				Public("Employment record not found.").
				With("caller_id", callerID).
				With("organization_id", organizationID).
				Errorf("not found")
		}
		return nil, oops.In("reader.self").
			Code(ErrCodeSelfEmploymentFailed).
			With("caller_id", callerID).
			With("organization_id", organizationID).
			Wrap(err)
	}
	return &out, nil
}
