//go:build integration

package async_integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	classifiersvc "github.com/medincident/medincident-backend/internal/service/command/incident/classifier"
	orgsvc "github.com/medincident/medincident-backend/internal/service/command/orgstructure"
)

func TestIncidentCategoryProjectedAsync(t *testing.T) {
	resetDBs(t)
	ctx := context.Background()

	org, err := orgSvc.Create(ctx, orgsvc.CreateOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateOrganizationPayload{Name: "Classifier Org", LegalAddress: orgsvc.AddressInput{Text: "addr"}},
	})
	require.NoError(t, err)
	requireProjected(t, `SELECT count(*) FROM projections.organizations WHERE id = ?`, org.ID)

	cat, err := classifierCatSvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{
			OrganizationID: org.ID.String(),
			Name:           "Async Category",
		},
	})
	require.NoError(t, err)

	requireProjected(t, `SELECT count(*) FROM projections.incident_categories WHERE id = ?`, cat.ID)

	var name string
	require.NoError(t, queryDB.Raw(
		`SELECT name FROM projections.incident_categories WHERE id = ?`, cat.ID,
	).Scan(&name).Error)
	assert.Equal(t, "Async Category", name)
}

func TestIncidentTypeProjectedAsync(t *testing.T) {
	resetDBs(t)
	ctx := context.Background()

	org, err := orgSvc.Create(ctx, orgsvc.CreateOrganizationCommand{
		Caller:  sysadminCaller,
		Payload: orgsvc.CreateOrganizationPayload{Name: "Type Org", LegalAddress: orgsvc.AddressInput{Text: "addr"}},
	})
	require.NoError(t, err)
	requireProjected(t, `SELECT count(*) FROM projections.organizations WHERE id = ?`, org.ID)

	cat, err := classifierCatSvc.Create(ctx, classifiersvc.CreateIncidentCategoryCommand{
		Caller:  sysadminCaller,
		Payload: classifiersvc.CreateIncidentCategoryPayload{OrganizationID: org.ID.String(), Name: "Parent Category"},
	})
	require.NoError(t, err)
	requireProjected(t, `SELECT count(*) FROM projections.incident_categories WHERE id = ?`, cat.ID)

	typ, err := classifierTypSvc.Create(ctx, classifiersvc.CreateIncidentTypeCommand{
		Caller: sysadminCaller,
		Payload: classifiersvc.CreateIncidentTypePayload{
			CategoryID: cat.ID.String(),
			Name:       "Async Type",
		},
	})
	require.NoError(t, err)

	requireProjected(t, `SELECT count(*) FROM projections.incident_types WHERE id = ?`, typ.ID)

	var catIDVal string
	require.NoError(t, queryDB.Raw(
		`SELECT category_id::text FROM projections.incident_types WHERE id = ?`, typ.ID,
	).Scan(&catIDVal).Error)
	assert.Equal(t, cat.ID.String(), catIDVal)
}
