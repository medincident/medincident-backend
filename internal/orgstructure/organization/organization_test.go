package organization_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/stretchr/testify/require"

	"github.com/medincident/medincident-command-service/internal/orgstructure/organization"
	"github.com/medincident/medincident-command-service/internal/shared/geo"
)

func validAddress(t *testing.T) *geo.Address {
	t.Helper()
	p, err := geo.NewPoint(30.5, 50.5)
	require.NoError(t, err)
	addr, err := geo.NewAddress("Main St 1", &p)
	require.NoError(t, err)
	return &addr
}

func TestNewValidRaisesCreated(t *testing.T) {
	now := time.Date(2026, 4, 10, 12, 0, 0, 0, time.UTC)
	o, err := organization.New("Acme Clinic", "a healthcare org", validAddress(t), now)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, o.ID)
	require.Equal(t, "Acme Clinic", o.Name)
	require.Equal(t, "a healthcare org", o.Description)
	require.NotNil(t, o.LegalAddress)
	require.Equal(t, now, o.CreatedAt)
	require.Equal(t, now, o.UpdatedAt)

	events := o.PullEvents()
	require.Len(t, events, 1)
	created, ok := events[0].(*organization.Created)
	require.True(t, ok)
	require.Equal(t, o.ID, created.ID)
	require.Equal(t, "Acme Clinic", created.Name)
	require.Equal(t, now, created.At)
}

func TestNewTrimsNameAndDescription(t *testing.T) {
	now := time.Now().UTC()
	o, err := organization.New("  Acme  ", "  desc  ", nil, now)
	require.NoError(t, err)
	require.Equal(t, "Acme", o.Name)
	require.Equal(t, "desc", o.Description)
}

func TestNewAcceptsNilLegalAddress(t *testing.T) {
	now := time.Now().UTC()
	o, err := organization.New("Acme", "", nil, now)
	require.NoError(t, err)
	require.Nil(t, o.LegalAddress)
}

func TestNewRejectsEmptyName(t *testing.T) {
	_, err := organization.New("   ", "", nil, time.Now())
	require.Error(t, err)
	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	require.Equal(t, organization.ErrCodeOrganizationNameEmpty, oe.Code())
	require.Equal(t, "name", oe.Context()["field"])
}

func TestNewRejectsShortName(t *testing.T) {
	_, err := organization.New("abc", "", nil, time.Now())
	require.Error(t, err)
	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	require.Equal(t, organization.ErrCodeOrganizationNameTooShort, oe.Code())
}

func TestNewRejectsShortDescription(t *testing.T) {
	_, err := organization.New("Acme", "abc", nil, time.Now())
	require.Error(t, err)
	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	require.Equal(t, organization.ErrCodeOrganizationDescriptionTooShort, oe.Code())
}

func TestNewCollectsNameAndDescriptionErrors(t *testing.T) {
	_, err := organization.New("", strings.Repeat("x", 2001), nil, time.Now())
	require.Error(t, err)
	var joined interface{ Unwrap() []error }
	require.True(t, errors.As(err, &joined))
	require.GreaterOrEqual(t, len(joined.Unwrap()), 2)
}

func TestHydrateDoesNotValidateOrRaise(t *testing.T) {
	id, err := uuid.NewV7()
	require.NoError(t, err)
	created := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	updated := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	// Name empty — would fail New, but Hydrate trusts the store.
	o := organization.Hydrate(id, "", "", nil, created, updated)
	require.Equal(t, id, o.ID)
	require.Equal(t, created, o.CreatedAt)
	require.Equal(t, updated, o.UpdatedAt)
	require.Empty(t, o.PullEvents())
}

func TestOrganizationSatisfiesEventSourceContract(t *testing.T) {
	now := time.Now().UTC()
	o, err := organization.New("Acme", "", nil, now)
	require.NoError(t, err)
	require.Equal(t, "organization", o.AggregateType())
	require.Equal(t, o.ID.String(), o.AggregateID())
	require.Len(t, o.PullEvents(), 1) // inherited from aggregate.Root
}

func TestRenameRaisesAndUpdatesTimestamp(t *testing.T) {
	now := time.Date(2026, 4, 10, 12, 0, 0, 0, time.UTC)
	o, err := organization.New("Acme", "", nil, now)
	require.NoError(t, err)
	_ = o.PullEvents() // drain Created

	later := now.Add(1 * time.Hour)
	require.NoError(t, o.Rename("AcmeCo", later))
	require.Equal(t, "AcmeCo", o.Name)
	require.Equal(t, later, o.UpdatedAt)

	events := o.PullEvents()
	require.Len(t, events, 1)
	renamed, ok := events[0].(*organization.Renamed)
	require.True(t, ok)
	require.Equal(t, "AcmeCo", renamed.Name)
}

func TestRenameNoopOnSameName(t *testing.T) {
	now := time.Now().UTC()
	o, err := organization.New("Acme", "", nil, now)
	require.NoError(t, err)
	_ = o.PullEvents()

	require.NoError(t, o.Rename("Acme", now.Add(1*time.Hour)))
	require.Equal(t, now, o.UpdatedAt, "no-op must not advance UpdatedAt")
	require.Empty(t, o.PullEvents())
}

func TestRenameRejectsEmpty(t *testing.T) {
	now := time.Now().UTC()
	o, err := organization.New("Acme", "", nil, now)
	require.NoError(t, err)

	err = o.Rename("   ", now.Add(time.Minute))
	require.Error(t, err)
	require.Equal(t, "Acme", o.Name, "bad rename must not mutate state")
}

func TestUpdateDescriptionEmptyClears(t *testing.T) {
	now := time.Now().UTC()
	o, err := organization.New("Acme", "before", nil, now)
	require.NoError(t, err)
	_ = o.PullEvents()

	later := now.Add(1 * time.Minute)
	require.NoError(t, o.UpdateDescription("  ", later))
	require.Equal(t, "", o.Description)
	events := o.PullEvents()
	require.Len(t, events, 1)
	upd, ok := events[0].(*organization.DescriptionUpdated)
	require.True(t, ok)
	require.Equal(t, "", upd.Description)
}

func TestUpdateDescriptionNoopOnSame(t *testing.T) {
	now := time.Now().UTC()
	o, err := organization.New("Acme", "hello", nil, now)
	require.NoError(t, err)
	_ = o.PullEvents()

	require.NoError(t, o.UpdateDescription("hello", now.Add(time.Minute)))
	require.Equal(t, now, o.UpdatedAt)
	require.Empty(t, o.PullEvents())
}

func TestRelocateLegalAddressWithNilClears(t *testing.T) {
	now := time.Now().UTC()
	o, err := organization.New("Acme", "", validAddress(t), now)
	require.NoError(t, err)
	_ = o.PullEvents()

	later := now.Add(1 * time.Minute)
	require.NoError(t, o.RelocateLegalAddress(nil, later))
	require.Nil(t, o.LegalAddress)

	events := o.PullEvents()
	require.Len(t, events, 1)
	rel, ok := events[0].(*organization.LegalAddressRelocated)
	require.True(t, ok)
	require.Nil(t, rel.LegalAddress)
}

func TestRelocateLegalAddressNoopOnEqual(t *testing.T) {
	now := time.Now().UTC()
	o, err := organization.New("Acme", "", validAddress(t), now)
	require.NoError(t, err)
	_ = o.PullEvents()

	// Pass an equivalent (but distinct pointer) Address.
	require.NoError(t, o.RelocateLegalAddress(validAddress(t), now.Add(time.Minute)))
	require.Equal(t, now, o.UpdatedAt)
	require.Empty(t, o.PullEvents())
}

func TestRelocateLegalAddressChangesWhenPointDiffers(t *testing.T) {
	now := time.Now().UTC()
	o, err := organization.New("Acme", "", validAddress(t), now)
	require.NoError(t, err)
	_ = o.PullEvents()

	other, _ := geo.NewPoint(31.0, 51.0)
	newAddr, _ := geo.NewAddress("Main St 1", &other)
	require.NoError(t, o.RelocateLegalAddress(&newAddr, now.Add(time.Minute)))
	require.NotNil(t, o.LegalAddress)
	require.Equal(t, 31.0, o.LegalAddress.Point.Longitude)
	require.Len(t, o.PullEvents(), 1)
}
