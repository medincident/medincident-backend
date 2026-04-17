//go:build integration

package membership_integration_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/medincident/medincident-command-service/internal/model"
	envelopev1 "github.com/medincident/medincident-command-service/pkg/event/v1"
)

// fixture holds the IDs of pre-seeded org/clinic/department rows that
// every test uses. takeFixture creates a fresh tree and registers a
// t.Cleanup to wipe it at the end of the test.
type fixture struct {
	OrgA     uuid.UUID
	ClinicA1 uuid.UUID
	DeptA1a  uuid.UUID
	DeptA1b  uuid.UUID
	ClinicA2 uuid.UUID
	DeptA2a  uuid.UUID
	OrgB     uuid.UUID
	ClinicB1 uuid.UUID
	DeptB1a  uuid.UUID
}

func takeFixture(t *testing.T) fixture {
	t.Helper()
	f := fixture{
		OrgA:     uuid.Must(uuid.NewV7()),
		ClinicA1: uuid.Must(uuid.NewV7()),
		DeptA1a:  uuid.Must(uuid.NewV7()),
		DeptA1b:  uuid.Must(uuid.NewV7()),
		ClinicA2: uuid.Must(uuid.NewV7()),
		DeptA2a:  uuid.Must(uuid.NewV7()),
		OrgB:     uuid.Must(uuid.NewV7()),
		ClinicB1: uuid.Must(uuid.NewV7()),
		DeptB1a:  uuid.Must(uuid.NewV7()),
	}
	must := func(err error) { require.NoError(t, err) }
	must(testDB.Exec(`INSERT INTO domain.organizations (id, name, legal_address) VALUES (?, ?, ROW(?, NULL)::domain.address)`, f.OrgA, "Org A", "addr-a").Error)
	must(testDB.Exec(`INSERT INTO domain.clinics (id, organization_id, name, physical_address) VALUES (?, ?, ?, ROW(?, NULL)::domain.address)`, f.ClinicA1, f.OrgA, "Clinic A1", "caddr-a1").Error)
	must(testDB.Exec(`INSERT INTO domain.departments (id, clinic_id, name) VALUES (?, ?, ?)`, f.DeptA1a, f.ClinicA1, "Dept A1a").Error)
	must(testDB.Exec(`INSERT INTO domain.departments (id, clinic_id, name) VALUES (?, ?, ?)`, f.DeptA1b, f.ClinicA1, "Dept A1b").Error)
	must(testDB.Exec(`INSERT INTO domain.clinics (id, organization_id, name, physical_address) VALUES (?, ?, ?, ROW(?, NULL)::domain.address)`, f.ClinicA2, f.OrgA, "Clinic A2", "caddr-a2").Error)
	must(testDB.Exec(`INSERT INTO domain.departments (id, clinic_id, name) VALUES (?, ?, ?)`, f.DeptA2a, f.ClinicA2, "Dept A2a").Error)
	must(testDB.Exec(`INSERT INTO domain.organizations (id, name, legal_address) VALUES (?, ?, ROW(?, NULL)::domain.address)`, f.OrgB, "Org B", "addr-b").Error)
	must(testDB.Exec(`INSERT INTO domain.clinics (id, organization_id, name, physical_address) VALUES (?, ?, ?, ROW(?, NULL)::domain.address)`, f.ClinicB1, f.OrgB, "Clinic B1", "cbaddr").Error)
	must(testDB.Exec(`INSERT INTO domain.departments (id, clinic_id, name) VALUES (?, ?, ?)`, f.DeptB1a, f.ClinicB1, "Dept B1a").Error)

	t.Cleanup(func() {
		// Role tables first — they have FK RESTRICT on employees, so
		// every employee delete downstream fails silently if any role
		// row survives. Child-to-parent order mirrors the FK chain.
		_ = testDB.Exec(`DELETE FROM domain.department_responsibles`).Error
		_ = testDB.Exec(`DELETE FROM domain.clinic_heads`).Error
		_ = testDB.Exec(`DELETE FROM domain.org_admins`).Error
		_ = testDB.Exec(`DELETE FROM domain.org_heads`).Error
		_ = testDB.Exec(`DELETE FROM domain.org_dispatchers`).Error
		_ = testDB.Exec(`DELETE FROM domain.system_admins`).Error
		_ = testDB.Exec(`DELETE FROM domain.employee_vacations`).Error
		_ = testDB.Exec(`DELETE FROM domain.employees`).Error
		_ = testDB.Exec(`DELETE FROM domain.departments WHERE id IN (?, ?, ?, ?)`, f.DeptA1a, f.DeptA1b, f.DeptA2a, f.DeptB1a).Error
		_ = testDB.Exec(`DELETE FROM domain.clinics WHERE id IN (?, ?, ?)`, f.ClinicA1, f.ClinicA2, f.ClinicB1).Error
		_ = testDB.Exec(`DELETE FROM domain.organizations WHERE id IN (?, ?)`, f.OrgA, f.OrgB).Error
		_ = testDB.Exec(`DELETE FROM outbox.events`).Error
		// Projection tables — sync projector writes these alongside
		// every domain mutation, so the cleanup needs to sweep them
		// too or cross-test id reuse trips unique constraints.
		_ = testDB.Exec(`TRUNCATE TABLE projections.organization_counters,
		                                 projections.clinic_counters,
		                                 projections.department_counters,
		                                 projections.employee_cards,
		                                 projections.employee_vacations,
		                                 projections.employees,
		                                 projections.departments,
		                                 projections.clinics,
		                                 projections.organizations,
		                                 projections.clinic_heads,
		                                 projections.department_responsibles,
		                                 projections.org_admins,
		                                 projections.org_dispatchers,
		                                 projections.org_heads,
		                                 projections.system_admins
		                         CASCADE`).Error
	})
	return f
}

// latestOutbox returns all outbox rows ordered by id ascending.
func latestOutbox(t *testing.T) []model.OutboxEvent {
	t.Helper()
	var rows []model.OutboxEvent
	require.NoError(t, testDB.Order("id ASC").Find(&rows).Error)
	return rows
}

func truncateOutbox(t *testing.T) {
	t.Helper()
	require.NoError(t, testDB.Exec(`DELETE FROM outbox.events`).Error)
}

// decodePayload unmarshals an outbox row into its envelope and then
// into the caller-provided target proto message.
func decodePayload(t *testing.T, row model.OutboxEvent, target proto.Message) *envelopev1.Envelope {
	t.Helper()
	env := &envelopev1.Envelope{}
	require.NoError(t, proto.Unmarshal(row.Payload, env))
	require.NoError(t, env.Payload.UnmarshalTo(target))
	return env
}

// ctxT returns a 10-second context cancelled on test cleanup.
func ctxT(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func mustParseUUID(t *testing.T, s string) uuid.UUID {
	t.Helper()
	id, err := uuid.Parse(s)
	require.NoError(t, err)
	return id
}

func uuidMustV7() uuid.UUID { return uuid.Must(uuid.NewV7()) }

// oopsCode extracts the domain error code from an oops-wrapped error,
// failing the test if err is not an oops error or if the code is not
// a string.
func oopsCode(t *testing.T, err error) string {
	t.Helper()
	require.Error(t, err)
	var oe oops.OopsError
	require.True(t, errors.As(err, &oe), "error is not an oops error: %v", err)
	code, ok := oe.Code().(string)
	require.True(t, ok, "oops.Code() returned non-string: %T %v", oe.Code(), oe.Code())
	return code
}

// oopsCodes walks a joined error chain (errors.Join output) and
// collects the oops .Code() string from every oops error in the tree.
// Used by tests that assert multiple field-level violations are all
// surfaced by a single multi-error return.
func oopsCodes(t *testing.T, err error) []string {
	t.Helper()
	require.Error(t, err)
	var out []string
	var walk func(error)
	walk = func(e error) {
		if e == nil {
			return
		}
		if joined, ok := e.(interface{ Unwrap() []error }); ok {
			for _, child := range joined.Unwrap() {
				walk(child)
			}
			return
		}
		var oe oops.OopsError
		if errors.As(e, &oe) {
			if code, ok := oe.Code().(string); ok && code != "" {
				out = append(out, code)
			}
		}
	}
	walk(err)
	return out
}
