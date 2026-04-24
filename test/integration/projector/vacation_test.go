//go:build integration

package projector_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/model"
	"github.com/medincident/medincident-backend/internal/service/command/projector"
)

// seedEmployee creates an org+clinic+dept+employee chain via the
// projector so vacation tests have a real employee row + card.
func seedEmployee(t *testing.T, ctx context.Context, now time.Time) *model.Employee {
	t.Helper()
	org, _, dept := seedDepartment(t, ctx, now)
	emp := &model.Employee{
		ID:             uuid.Must(uuid.NewV7()),
		ZitadelUserID:  "zit-user-" + uuid.NewString(),
		OrganizationID: org.ID,
		DepartmentID:   dept.ID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.EmployeeHired(tx, emp)
	}))
	return emp
}

// TestVacationScheduled_InsertsRowAndRefreshesCard confirms the
// projector inserts a scheduled vacation and updates the employee
// card's next_vacation_* pointers.
func TestVacationScheduled_InsertsRowAndRefreshesCard(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	emp := seedEmployee(t, ctx, now)

	starts := now.Add(24 * time.Hour)
	ends := starts.Add(48 * time.Hour)
	vac := &model.EmployeeVacation{
		ID:         uuid.Must(uuid.NewV7()),
		EmployeeID: emp.ID,
		StartsAt:   starts,
		EndsAt:     null.TimeFrom(ends),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.VacationScheduled(tx, vac)
	}))

	var state string
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT state FROM projections.employee_vacations WHERE id = ?`, vac.ID,
	).Row().Scan(&state))
	require.Equal(t, "scheduled", state)

	var nextID *uuid.UUID
	var nextStarts *time.Time
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT next_vacation_id, next_vacation_starts_at FROM projections.employee_cards WHERE employee_id = ?`, emp.ID,
	).Row().Scan(&nextID, &nextStarts))
	require.NotNil(t, nextID)
	require.Equal(t, vac.ID, *nextID)
}

// TestVacationStarted_InsertsActiveRow confirms the projector stamps
// the current_vacation_* pointers on the card.
func TestVacationStarted_InsertsActiveRow(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	emp := seedEmployee(t, ctx, now)

	ends := now.Add(48 * time.Hour)
	vac := &model.EmployeeVacation{
		ID:         uuid.Must(uuid.NewV7()),
		EmployeeID: emp.ID,
		StartsAt:   now,
		EndsAt:     null.TimeFrom(ends),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.VacationStarted(tx, vac)
	}))

	var state string
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT state FROM projections.employee_vacations WHERE id = ?`, vac.ID,
	).Row().Scan(&state))
	require.Equal(t, "active", state)

	var currentID *uuid.UUID
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT current_vacation_id FROM projections.employee_cards WHERE employee_id = ?`, emp.ID,
	).Row().Scan(&currentID))
	require.NotNil(t, currentID)
	require.Equal(t, vac.ID, *currentID)
}

// TestVacationEnded_StampsEndsAtAndClearsCard confirms ended state +
// employee_card current_vacation_id cleared.
func TestVacationEnded_StampsEndsAtAndClearsCard(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	emp := seedEmployee(t, ctx, now)

	vac := &model.EmployeeVacation{
		ID:         uuid.Must(uuid.NewV7()),
		EmployeeID: emp.ID,
		StartsAt:   now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.VacationStarted(tx, vac)
	}))

	endedAt := now.Add(4 * time.Hour)
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.VacationEnded(tx, vac, endedAt)
	}))

	var state string
	var endsAt *time.Time
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT state, ends_at FROM projections.employee_vacations WHERE id = ?`, vac.ID,
	).Row().Scan(&state, &endsAt))
	require.Equal(t, "ended", state)
	require.NotNil(t, endsAt)

	var currentID *uuid.UUID
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT current_vacation_id FROM projections.employee_cards WHERE employee_id = ?`, emp.ID,
	).Row().Scan(&currentID))
	require.Nil(t, currentID)
}

// TestVacationCancelled_SetsStateAndClearsCard confirms the projector
// marks the scheduled vacation cancelled and clears next_vacation_id.
func TestVacationCancelled_SetsStateAndClearsCard(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	emp := seedEmployee(t, ctx, now)

	starts := now.Add(48 * time.Hour)
	vac := &model.EmployeeVacation{
		ID:         uuid.Must(uuid.NewV7()),
		EmployeeID: emp.ID,
		StartsAt:   starts,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.VacationScheduled(tx, vac)
	}))

	cancelledAt := now.Add(time.Hour)
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.VacationCancelled(tx, vac, cancelledAt)
	}))

	var state string
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT state FROM projections.employee_vacations WHERE id = ?`, vac.ID,
	).Row().Scan(&state))
	require.Equal(t, "cancelled", state)

	var nextID *uuid.UUID
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT next_vacation_id FROM projections.employee_cards WHERE employee_id = ?`, emp.ID,
	).Row().Scan(&nextID))
	require.Nil(t, nextID)
}

// TestVacationEndDateChanged_UpdatesEndsAt confirms ends_at rewrites
// propagate to the projection row.
func TestVacationEndDateChanged_UpdatesEndsAt(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	emp := seedEmployee(t, ctx, now)

	ends := now.Add(24 * time.Hour)
	vac := &model.EmployeeVacation{
		ID:         uuid.Must(uuid.NewV7()),
		EmployeeID: emp.ID,
		StartsAt:   now,
		EndsAt:     null.TimeFrom(ends),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.VacationStarted(tx, vac)
	}))

	newEnd := ends.Add(24 * time.Hour)
	vac.EndsAt = null.TimeFrom(newEnd)
	vac.UpdatedAt = now.Add(time.Hour)
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return projector.VacationEndDateChanged(tx, vac)
	}))

	var gotEnd time.Time
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT ends_at FROM projections.employee_vacations WHERE id = ?`, vac.ID,
	).Row().Scan(&gotEnd))
	require.True(t, gotEnd.Equal(newEnd))
}
