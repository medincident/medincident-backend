//go:build integration

package organization_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	organizationapp "github.com/medincident/medincident-command-service/internal/orgstructure/organization/app"
	organizationinfra "github.com/medincident/medincident-command-service/internal/orgstructure/organization/infra"
	"github.com/medincident/medincident-command-service/internal/outbox"
	"github.com/medincident/medincident-command-service/internal/storage/postgres"
)

// wallClock is a production-shaped clock used to prove the service
// works with real time. Tests assert only that CreatedAt / UpdatedAt
// are set and consistent, not their absolute values.
type wallClock struct{}

func (wallClock) Now() time.Time { return time.Now() }

// newService builds a Service bound to the shared testPool. Shared
// across test functions so each test exercises the real full pipeline
// (aggregate → repository → outbox → commit) without re-wiring.
func newService(t *testing.T) organizationapp.Service {
	t.Helper()
	log := zerolog.Nop()
	pool := &postgres.Pool{Pool: testPool}
	repo := postgres.NewOrganizationRepo(pool)
	store := postgres.NewOutboxStore()
	reg := outbox.NewRegistry()
	organizationinfra.RegisterOutboxMappers(reg)
	beg := postgres.NewBeginner(pool)
	clk := wallClock{}
	pub := outbox.NewPublisher(store, reg, clk)
	return organizationapp.NewService(&log, beg, repo, pub, clk)
}

// cleanup truncates the two tables this test suite writes to, so each
// test starts with a predictable blank state without tearing down the
// whole container between cases.
func cleanup(t *testing.T) {
	t.Helper()
	_, err := testPool.Exec(context.Background(), `
		TRUNCATE TABLE outbox.events, domain.organizations RESTART IDENTITY CASCADE
	`)
	require.NoError(t, err)
}

func TestCreateOrganizationFullPipeline(t *testing.T) {
	cleanup(t)
	svc := newService(t)

	res, err := svc.Create(context.Background(), organizationapp.CreateCommand{
		Name:        "Acme Clinic",
		Description: "integration test",
		LegalAddress: &organizationapp.AddressInput{
			Text:  "Main St 1",
			Point: &organizationapp.PointInput{Longitude: 30.5234, Latitude: 50.4501},
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, res.ID)

	// Verify the aggregate is in domain.organizations with the expected
	// flat-column values.
	var (
		name     string
		desc     *string
		addrText *string
		lng, lat *float64
	)
	err = testPool.QueryRow(context.Background(), `
		SELECT name, description, legal_address_text, legal_address_lng, legal_address_lat
		FROM domain.organizations
		WHERE id = $1
	`, res.ID).Scan(&name, &desc, &addrText, &lng, &lat)
	require.NoError(t, err)
	require.Equal(t, "Acme Clinic", name)
	require.NotNil(t, desc)
	require.Equal(t, "integration test", *desc)
	require.NotNil(t, addrText)
	require.Equal(t, "Main St 1", *addrText)
	require.NotNil(t, lng)
	require.InDelta(t, 30.5234, *lng, 1e-9)
	require.NotNil(t, lat)
	require.InDelta(t, 50.4501, *lat, 1e-9)

	// Verify outbox.events has a matching row with all envelope metadata populated.
	var (
		aggType     string
		aggID       string
		subject     string
		eventID     uuid.UUID
		headersJSON []byte
		payload     []byte
	)
	err = testPool.QueryRow(context.Background(), `
		SELECT aggregate_type, aggregate_id, subject, event_id, headers, payload
		FROM outbox.events
		ORDER BY created_at DESC
		LIMIT 1
	`).Scan(&aggType, &aggID, &subject, &eventID, &headersJSON, &payload)
	require.NoError(t, err)
	require.Equal(t, "organization", aggType)
	require.Equal(t, res.ID.String(), aggID)
	require.Equal(t, "medincident.orgstructure.v1.organization.created", subject)
	require.NotEqual(t, uuid.Nil, eventID)
	require.Contains(t, string(headersJSON), eventID.String(), "Nats-Msg-Id header should carry the event id")
	require.NotEmpty(t, payload)
}

func TestCreateOrganizationRejectsInvalidCoordinatesBeforeWrite(t *testing.T) {
	cleanup(t)
	svc := newService(t)

	_, err := svc.Create(context.Background(), organizationapp.CreateCommand{
		Name: "Acme",
		LegalAddress: &organizationapp.AddressInput{
			Text:  "Main St",
			Point: &organizationapp.PointInput{Longitude: 999, Latitude: 50},
		},
	})
	require.Error(t, err)

	var count int
	require.NoError(t, testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM domain.organizations`).Scan(&count))
	require.Equal(t, 0, count)
	require.NoError(t, testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM outbox.events`).Scan(&count))
	require.Equal(t, 0, count)
}

func TestRenameAppendsOutboxRowInSameTransaction(t *testing.T) {
	cleanup(t)
	svc := newService(t)

	createRes, err := svc.Create(context.Background(), organizationapp.CreateCommand{
		Name: "OldOrg",
	})
	require.NoError(t, err)

	_, err = svc.Rename(context.Background(), organizationapp.RenameCommand{
		ID:   createRes.ID,
		Name: "NewOrg",
	})
	require.NoError(t, err)

	var name string
	require.NoError(t, testPool.QueryRow(context.Background(),
		`SELECT name FROM domain.organizations WHERE id = $1`, createRes.ID,
	).Scan(&name))
	require.Equal(t, "NewOrg", name)

	// Two outbox rows: Created + Renamed, both for the same aggregate.
	var rowCount int
	require.NoError(t, testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM outbox.events WHERE aggregate_id = $1`, createRes.ID.String(),
	).Scan(&rowCount))
	require.Equal(t, 2, rowCount)
}
