//go:build integration

package identity_query_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	zitadelv1 "github.com/medincident/medincident-zitadel-actions/gen/zitadel/events/v1"
	usersv1 "github.com/medincident/medincident-zitadel-actions/gen/zitadel/users/v1"

	"github.com/medincident/medincident-backend/internal/config"
	identityread "github.com/medincident/medincident-backend/internal/service/query/identity"
)

// TestConsumer_UserHumanAdded_EndToEnd publishes a Zitadel user-added
// event onto the JetStream stream and asserts the projector writes
// the projections.users row within 5 seconds.
func TestConsumer_UserHumanAdded_EndToEnd(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	nc, err := nats.Connect(natsURL)
	require.NoError(t, err)
	t.Cleanup(nc.Close)
	js, err := jetstream.New(nc)
	require.NoError(t, err)

	projector := identityread.NewProjector(testDB, &logger)
	cfg := &config.NATSConfig{
		URL:         natsURL,
		Stream:      "zitadel",
		Subjects:    []string{"zitadel.>"},
		DurableName: "test-query-server-identity",
	}
	consumer := identityread.NewConsumer(js, cfg, projector, &logger)
	require.NoError(t, consumer.Start(ctx))
	t.Cleanup(func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = consumer.Shutdown(shutdownCtx)
	})

	userID := "zit-user-123"
	payload := &usersv1.UserHumanAdded{
		UserName:          "jdoe",
		FirstName:         "Jane",
		LastName:          "Doe",
		DisplayName:       "Jane Doe",
		Email:             "jane@example.com",
		PreferredLanguage: "en",
	}
	anyPayload, err := anypb.New(payload)
	require.NoError(t, err)

	env := &zitadelv1.Envelope{
		AggregateId: userID,
		EventType:   identityread.ZitadelEventUserHumanAdded,
		CreatedAt:   timestamppb.New(time.Now().UTC()),
		Payload:     anyPayload,
	}
	envBytes, err := proto.Marshal(env)
	require.NoError(t, err)

	_, err = js.Publish(ctx, "zitadel.users.added", envBytes)
	require.NoError(t, err)

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		var count int64
		if err := testDB.WithContext(ctx).
			Raw(`SELECT count(*) FROM projections.users WHERE id = ?`, userID).
			Row().Scan(&count); err == nil && count == 1 {
			var email string
			var displayName string
			require.NoError(t, testDB.WithContext(ctx).
				Raw(`SELECT email, display_name FROM projections.users WHERE id = ?`, userID).
				Row().Scan(&email, &displayName))
			require.Equal(t, "jane@example.com", email)
			require.Equal(t, "Jane Doe", displayName)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("projections.users row did not appear within deadline")
}
