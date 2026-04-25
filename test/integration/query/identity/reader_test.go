//go:build integration

package identity_query_integration_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"github.com/stretchr/testify/require"

	identityread "github.com/medincident/medincident-backend/internal/service/query/identity"
)

// TestReader_GetUserByEmail seeds a projections.users row directly and
// asserts the reader returns it for case-insensitive email lookups and
// reports a not-found domain error otherwise.
func TestReader_GetUserByEmail(t *testing.T) {
	resetProjections(t)
	ctx := context.Background()
	logger := zerolog.Nop()

	now := time.Now().UTC().Truncate(time.Second)
	raw, err := testDB.DB()
	require.NoError(t, err)
	_, err = raw.ExecContext(ctx, `
		INSERT INTO projections.users (
			id, user_name, first_name, last_name, display_name,
			email, email_verified, preferred_language, gender,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $10)`,
		"zit-user-1", "jdoe", "Jane", "Doe", "Jane Doe",
		"Jane@Example.com", true, "en", int32(1), now,
	)
	require.NoError(t, err)

	reader := identityread.NewReader(testDB, authzSvc, &logger)

	t.Run("exact match returns the row", func(t *testing.T) {
		v, err := reader.GetUserByEmail(ctx, sysadminCaller, "Jane@Example.com")
		require.NoError(t, err)
		require.Equal(t, "zit-user-1", v.ID)
		require.Equal(t, "Jane@Example.com", v.Email)
	})

	t.Run("case-insensitive match", func(t *testing.T) {
		v, err := reader.GetUserByEmail(ctx, sysadminCaller, "jane@example.com")
		require.NoError(t, err)
		require.Equal(t, "zit-user-1", v.ID)
	})

	t.Run("not found returns domain error", func(t *testing.T) {
		_, err := reader.GetUserByEmail(ctx, sysadminCaller, "missing@example.com")
		require.Error(t, err)
		var oe oops.OopsError
		require.True(t, errors.As(err, &oe))
		require.Equal(t, identityread.ErrCodeUserNotFound, oe.Code())
	})
}
