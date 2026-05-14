//go:build integration

package async_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/outbox"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
	empv1 "github.com/medincident/medincident-backend/pkg/event/employee/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
)

// appendEmployeeHiredEvent inserts an employee row into domain.employees and
// appends an EmployeeHired event to the outbox in one transaction so the
// async projection pipeline can be exercised without a live Zitadel instance.
func appendEmployeeHiredEvent(
	t *testing.T,
	empID uuid.UUID,
	zitadelUserID string,
	orgID, deptID uuid.UUID,
	position *string,
) {
	t.Helper()
	hiredAt := time.Now().UTC()
	err := commandDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(
			`INSERT INTO domain.employees
			    (id, zitadel_user_id, organization_id, department_id, position, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			empID, zitadelUserID, orgID, deptID, position, hiredAt, hiredAt,
		).Error; err != nil {
			return err
		}
		msg := &empv1.EmployeeHired{
			ZitadelUserId:  zitadelUserID,
			OrganizationId: orgID.String(),
			DepartmentId:   deptID.String(),
			HiredAt:        timestamppb.New(hiredAt),
			DisplayName:    zitadelUserID,
		}
		if position != nil {
			msg.Position = *position
		}
		payload, err := anypb.New(msg)
		if err != nil {
			return err
		}
		env := &eventv1.Envelope{
			OccurredAt:    timestamppb.New(hiredAt),
			AggregateType: "employee",
			AggregateId:   empID.String(),
			Payload:       payload,
		}
		return outbox.Append(tx, "medincident.event.employee.v1.hired", env)
	})
	require.NoError(t, err, "appendEmployeeHiredEvent")
}

func TestEmployeeProjectedAsync(t *testing.T) {
	resetDBs(t)
	ctx := context.Background()

	org, err := orgSvc.Create(ctx, orgsvc.CreateOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateOrganizationPayload{Name: "Async Emp Org", LegalAddress: orgsvc.AddressInput{Text: "addr"}},
	})
	require.NoError(t, err)
	requireProjected(t, `SELECT count(*) FROM projections.organizations WHERE id = ?`, org.ID)

	clin, err := clinSvc.Create(ctx, orgsvc.CreateClinicCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateClinicPayload{OrganizationID: org.ID.String(), Name: "Clinic", PhysicalAddress: orgsvc.AddressInput{Text: "addr"}},
	})
	require.NoError(t, err)
	requireProjected(t, `SELECT count(*) FROM projections.clinics WHERE id = ?`, clin.ID)

	dept, err := deptSvc.Create(ctx, orgsvc.CreateDepartmentCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateDepartmentPayload{ClinicID: clin.ID.String(), Name: "Dept"},
	})
	require.NoError(t, err)
	requireProjected(t, `SELECT count(*) FROM projections.departments WHERE id = ?`, dept.ID)

	empID, err := uuid.NewV7()
	require.NoError(t, err)
	pos := "Врач"
	appendEmployeeHiredEvent(t, empID, "alice-zid", org.ID, dept.ID, &pos)

	requireProjected(t, `SELECT count(*) FROM projections.employees WHERE id = ?`, empID)

	var orgIDVal, deptIDVal string
	require.NoError(t, queryDB.Raw(
		`SELECT organization_id::text, department_id::text FROM projections.employees WHERE id = ?`, empID,
	).Row().Scan(&orgIDVal, &deptIDVal))
	assert.Equal(t, org.ID.String(), orgIDVal, "projections.employees.organization_id")
	assert.Equal(t, dept.ID.String(), deptIDVal, "projections.employees.department_id")
}
