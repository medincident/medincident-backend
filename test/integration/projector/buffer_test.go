//go:build integration

package projector_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"

	qprojector "github.com/medincident/medincident-backend/internal/service/query/projector"
	bufferv1 "github.com/medincident/medincident-backend/pkg/event/incident/buffer/v1"
)

func TestBuffer_ProjectorCreated(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	bufID := uuid.Must(uuid.NewV7())
	orgID := uuid.Must(uuid.NewV7())

	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.PatientIncidentBufferCreated(tx, bufID.String(), now,
			&bufferv1.PatientIncidentBufferCreated{
				OrganizationId:       orgID.String(),
				PatientZitadelUserId: "zitadel-user-1",
				Description:          "raw patient text",
				Summary:              "AI formalized text",
				Priority:             "high",
				Status:               "pending",
				CreatedAt:            timestamppb.New(now),
			})
	}))

	var gotDesc, gotSummary, gotPriority string
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT description, summary, priority FROM projections.patient_incident_buffer WHERE id = ?`, bufID,
	).Row().Scan(&gotDesc, &gotSummary, &gotPriority))
	require.Equal(t, "raw patient text", gotDesc)
	require.Equal(t, "AI formalized text", gotSummary)
	require.Equal(t, "high", gotPriority)
}

func TestBuffer_ProjectorUpdated_SummaryAndPriority(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	bufID := uuid.Must(uuid.NewV7())
	orgID := uuid.Must(uuid.NewV7())

	// Seed a buffer row directly.
	require.NoError(t, testDB.WithContext(ctx).Exec(`
		INSERT INTO projections.patient_incident_buffer
		(id, organization_id, patient_zitadel_user_id, description, summary, priority, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		bufID, orgID, "z1", "original", "original summary", "normal", "pending", now, now,
	).Error)

	// Update only priority — summary wrapper is nil, so summary should remain unchanged.
	require.NoError(t, testDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return qprojector.PatientIncidentBufferUpdated(tx, bufID.String(), now.Add(time.Second),
			&bufferv1.PatientIncidentBufferUpdated{
				Description: "original",
				Status:      "pending",
				Priority:    wrapperspb.String("high"),
				UpdatedAt:   timestamppb.New(now.Add(time.Second)),
			})
	}))

	var gotSummary, gotPriority string
	require.NoError(t, testDB.WithContext(ctx).Raw(
		`SELECT summary, priority FROM projections.patient_incident_buffer WHERE id = ?`, bufID,
	).Row().Scan(&gotSummary, &gotPriority))
	require.Equal(t, "original summary", gotSummary) // unchanged
	require.Equal(t, "high", gotPriority)            // updated
}
