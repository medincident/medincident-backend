package organizationinfra_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	orgstructurev1 "github.com/medincident/medincident-command-service/gen/medincident/orgstructure/v1"

	"github.com/medincident/medincident-command-service/internal/orgstructure/organization"
	organizationinfra "github.com/medincident/medincident-command-service/internal/orgstructure/organization/infra"
	"github.com/medincident/medincident-command-service/internal/outbox"
	"github.com/medincident/medincident-command-service/internal/shared/geo"
)

func TestRegisterOutboxMappersCoversAllEvents(t *testing.T) {
	reg := outbox.NewRegistry()
	organizationinfra.RegisterOutboxMappers(reg)

	for _, tp := range []reflect.Type{
		reflect.TypeOf(&organization.Created{}),
		reflect.TypeOf(&organization.Renamed{}),
		reflect.TypeOf(&organization.DescriptionUpdated{}),
		reflect.TypeOf(&organization.LegalAddressRelocated{}),
	} {
		info, ok := reg.ByGoType(tp)
		require.True(t, ok, "missing mapper for %s", tp)
		require.NotEmpty(t, info.TypeName)
		require.NotNil(t, info.Zero)
		require.NotNil(t, info.ToProto)
	}
}

func TestCreatedMapperFullAddress(t *testing.T) {
	reg := outbox.NewRegistry()
	organizationinfra.RegisterOutboxMappers(reg)

	info, ok := reg.ByGoType(reflect.TypeOf(&organization.Created{}))
	require.True(t, ok)

	addr := &geo.Address{Text: "Main 1", Point: &geo.Point{Longitude: 30.5, Latitude: 50.5}}
	msg, err := info.ToProto(&organization.Created{
		ID:           uuid.MustParse("01910000-0000-0000-0000-000000000000"),
		Name:         "Acme",
		Description:  "desc",
		LegalAddress: addr,
		At:           time.Unix(0, 0),
	})
	require.NoError(t, err)

	oc, ok := msg.(*orgstructurev1.OrganizationCreated)
	require.True(t, ok)
	require.Equal(t, "Acme", oc.Name)
	require.NotNil(t, oc.Description)
	require.Equal(t, "desc", *oc.Description)
	require.NotNil(t, oc.LegalAddress)
	require.Equal(t, "Main 1", oc.LegalAddress.Text)
	require.NotNil(t, oc.LegalAddress.Point)
	require.InDelta(t, 30.5, oc.LegalAddress.Point.Lng, 1e-9)
	require.InDelta(t, 50.5, oc.LegalAddress.Point.Lat, 1e-9)
}

func TestCreatedMapperOmitsEmptyDescription(t *testing.T) {
	reg := outbox.NewRegistry()
	organizationinfra.RegisterOutboxMappers(reg)
	info, _ := reg.ByGoType(reflect.TypeOf(&organization.Created{}))

	msg, err := info.ToProto(&organization.Created{Name: "Acme"})
	require.NoError(t, err)
	oc := msg.(*orgstructurev1.OrganizationCreated)
	require.Nil(t, oc.Description, "empty string must not set the optional proto field")
	require.Nil(t, oc.LegalAddress)
}

func TestRenamedMapper(t *testing.T) {
	reg := outbox.NewRegistry()
	organizationinfra.RegisterOutboxMappers(reg)
	info, _ := reg.ByGoType(reflect.TypeOf(&organization.Renamed{}))

	msg, err := info.ToProto(&organization.Renamed{Name: "AcmeCo"})
	require.NoError(t, err)
	require.Equal(t, "AcmeCo", msg.(*orgstructurev1.OrganizationRenamed).Name)
}

func TestDescriptionUpdatedMapperClearsWhenEmpty(t *testing.T) {
	reg := outbox.NewRegistry()
	organizationinfra.RegisterOutboxMappers(reg)
	info, _ := reg.ByGoType(reflect.TypeOf(&organization.DescriptionUpdated{}))

	msg, err := info.ToProto(&organization.DescriptionUpdated{Description: ""})
	require.NoError(t, err)
	require.Nil(t, msg.(*orgstructurev1.OrganizationDescriptionUpdated).Description)

	msg, err = info.ToProto(&organization.DescriptionUpdated{Description: "new"})
	require.NoError(t, err)
	require.Equal(t, "new", *msg.(*orgstructurev1.OrganizationDescriptionUpdated).Description)
}

func TestLegalAddressRelocatedMapperClearsWhenNil(t *testing.T) {
	reg := outbox.NewRegistry()
	organizationinfra.RegisterOutboxMappers(reg)
	info, _ := reg.ByGoType(reflect.TypeOf(&organization.LegalAddressRelocated{}))

	msg, err := info.ToProto(&organization.LegalAddressRelocated{LegalAddress: nil})
	require.NoError(t, err)
	require.Nil(t, msg.(*orgstructurev1.OrganizationLegalAddressRelocated).LegalAddress)

	addr := &geo.Address{Text: "New St"}
	msg, err = info.ToProto(&organization.LegalAddressRelocated{LegalAddress: addr})
	require.NoError(t, err)
	got := msg.(*orgstructurev1.OrganizationLegalAddressRelocated).LegalAddress
	require.NotNil(t, got)
	require.Equal(t, "New St", got.Text)
}
