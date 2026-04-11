package organizationapp_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/medincident/medincident-command-service/internal/orgstructure/organization"
	organizationapp "github.com/medincident/medincident-command-service/internal/orgstructure/organization/app"
	"github.com/medincident/medincident-command-service/internal/outbox"
	"github.com/medincident/medincident-command-service/internal/shared/geo"
	"github.com/medincident/medincident-command-service/internal/tx"
)

type fixedClock struct{ now time.Time }

func (f fixedClock) Now() time.Time { return f.now }

type inMemoryRepo struct {
	saved   map[uuid.UUID]*organization.Organization
	getErr  error
	saveErr error
}

func newInMemoryRepo() *inMemoryRepo {
	return &inMemoryRepo{saved: make(map[uuid.UUID]*organization.Organization)}
}

func (r *inMemoryRepo) GetByID(_ context.Context, id uuid.UUID) (*organization.Organization, error) {
	if r.getErr != nil {
		return nil, r.getErr
	}
	o, ok := r.saved[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return o, nil
}

func (r *inMemoryRepo) Save(_ context.Context, o *organization.Organization) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	r.saved[o.ID] = o
	return nil
}

func (r *inMemoryRepo) List(context.Context, organizationapp.ListFilter) ([]*organization.Organization, error) {
	panic("List not used in service tests")
}

// noopTx satisfies tx.Tx without touching any DB.
type noopTx struct{}

func (noopTx) Close() error                                    { return nil }
func (noopTx) Commit(context.Context) error                    { return nil }
func (noopTx) Rollback(context.Context) error                  { return nil }
func (noopTx) Savepoint(context.Context, string) error         { return nil }
func (noopTx) RollbackSavepoint(context.Context, string) error { return nil }
func (noopTx) ReleaseSavepoint(context.Context, string) error  { return nil }

type noopBeginner struct {
	err error
}

func (b noopBeginner) Begin(context.Context) (tx.Tx, error) {
	if b.err != nil {
		return nil, b.err
	}
	return noopTx{}, nil
}

type memStore struct {
	rows []*outbox.Record
}

func (s *memStore) Append(_ context.Context, _ tx.Tx, r *outbox.Record) error {
	if r == nil {
		return errors.New("nil record")
	}
	s.rows = append(s.rows, r)
	return nil
}

func registryForOrganization(t *testing.T) outbox.Registry {
	t.Helper()
	reg := outbox.NewRegistry()
	reg.Register(reflect.TypeOf(&organization.Created{}), outbox.EventInfo{
		TypeName: "medincident.orgstructure.v1.OrganizationCreated",
		Zero:     func() any { return &organization.Created{} },
		ToProto:  func(any) (proto.Message, error) { return &emptypb.Empty{}, nil },
	})
	reg.Register(reflect.TypeOf(&organization.Renamed{}), outbox.EventInfo{
		TypeName: "medincident.orgstructure.v1.OrganizationRenamed",
		Zero:     func() any { return &organization.Renamed{} },
		ToProto:  func(any) (proto.Message, error) { return &emptypb.Empty{}, nil },
	})
	reg.Register(reflect.TypeOf(&organization.DescriptionUpdated{}), outbox.EventInfo{
		TypeName: "medincident.orgstructure.v1.OrganizationDescriptionUpdated",
		Zero:     func() any { return &organization.DescriptionUpdated{} },
		ToProto:  func(any) (proto.Message, error) { return &emptypb.Empty{}, nil },
	})
	reg.Register(reflect.TypeOf(&organization.LegalAddressRelocated{}), outbox.EventInfo{
		TypeName: "medincident.orgstructure.v1.OrganizationLegalAddressRelocated",
		Zero:     func() any { return &organization.LegalAddressRelocated{} },
		ToProto:  func(any) (proto.Message, error) { return &emptypb.Empty{}, nil },
	})
	return reg
}

func newTestService(t *testing.T, clockNow time.Time) (*organizationapp.Service, *inMemoryRepo, *memStore) {
	t.Helper()
	log := zerolog.Nop()
	repo := newInMemoryRepo()
	store := &memStore{}
	reg := registryForOrganization(t)
	svc := organizationapp.NewService(&log, noopBeginner{}, repo, store, reg, fixedClock{now: clockNow})
	return svc, repo, store
}

func TestServiceCreateHappyPath(t *testing.T) {
	now := time.Date(2026, 4, 10, 12, 0, 0, 0, time.UTC)
	svc, repo, store := newTestService(t, now)

	res, err := svc.Create(context.Background(), organizationapp.CreateCommand{
		Name:        "Acme",
		Description: "desc",
		LegalAddress: &organizationapp.AddressInput{
			Text:  "Main St 1",
			Point: &organizationapp.PointInput{Longitude: 30.5, Latitude: 50.5},
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, res.ID)

	saved, ok := repo.saved[res.ID]
	require.True(t, ok)
	require.Equal(t, "Acme", saved.Name)
	require.Equal(t, "desc", saved.Description)
	require.NotNil(t, saved.LegalAddress)
	require.Equal(t, now, saved.CreatedAt)

	require.Len(t, store.rows, 1)
	require.Equal(t, "organization", store.rows[0].AggregateType)
	require.Equal(t, res.ID.String(), store.rows[0].AggregateID)
	require.Equal(t, "medincident.orgstructure.v1.OrganizationCreated", store.rows[0].EventType)
}

func TestServiceCreateNilAddressAccepted(t *testing.T) {
	svc, repo, store := newTestService(t, time.Now().UTC())
	res, err := svc.Create(context.Background(), organizationapp.CreateCommand{
		Name: "NoAddr",
	})
	require.NoError(t, err)
	require.Nil(t, repo.saved[res.ID].LegalAddress)
	require.Len(t, store.rows, 1)
}

func TestServiceCreateRejectsEmptyName(t *testing.T) {
	svc, repo, store := newTestService(t, time.Now().UTC())
	_, err := svc.Create(context.Background(), organizationapp.CreateCommand{Name: ""})
	require.Error(t, err)
	require.Empty(t, repo.saved)
	require.Empty(t, store.rows, "no outbox writes on validation failure")
}

func TestServiceCreateRejectsInvalidPoint(t *testing.T) {
	svc, repo, store := newTestService(t, time.Now().UTC())
	_, err := svc.Create(context.Background(), organizationapp.CreateCommand{
		Name: "Acme",
		LegalAddress: &organizationapp.AddressInput{
			Text:  "Main St",
			Point: &organizationapp.PointInput{Longitude: 999, Latitude: 50},
		},
	})
	require.Error(t, err)
	requireLeafCode(t, err, geo.ErrCodePointLongitudeOutOfRange)
	require.Empty(t, repo.saved)
	require.Empty(t, store.rows)
}

// TestServiceCreateCollectsAllErrorsAcrossAddressAndAggregate verifies
// that the operation-level validator does NOT stop at the first failure:
// a single Create call with a bad address AND a bad name must surface
// every violation in one joined error.
func TestServiceCreateCollectsAllErrorsAcrossAddressAndAggregate(t *testing.T) {
	svc, repo, store := newTestService(t, time.Now().UTC())
	_, err := svc.Create(context.Background(), organizationapp.CreateCommand{
		Name:        "", // organization_name_empty
		Description: "",
		LegalAddress: &organizationapp.AddressInput{
			Text:  "", // address_text_empty
			Point: &organizationapp.PointInput{Longitude: 999, Latitude: 999},
		},
	})
	require.Error(t, err)
	requireLeafCode(t, err, geo.ErrCodePointLongitudeOutOfRange)
	requireLeafCode(t, err, geo.ErrCodePointLatitudeOutOfRange)
	requireLeafCode(t, err, geo.ErrCodeAddressTextEmpty)
	requireLeafCode(t, err, organization.ErrCodeOrganizationNameEmpty)
	require.Empty(t, repo.saved)
	require.Empty(t, store.rows)
}

// requireLeafCode walks the joined error tree and asserts that at least
// one leaf carries the given oops Code. This is what the future gRPC
// error interceptor's flattenErrors will do to build BadRequest field
// violations.
func requireLeafCode(t *testing.T, err error, code string) {
	t.Helper()
	if hasLeafCode(err, code) {
		return
	}
	t.Fatalf("expected a leaf with Code(%q) in error tree, got: %v", code, err)
}

func hasLeafCode(err error, code string) bool {
	if err == nil {
		return false
	}
	if oe, ok := oops.AsOops(err); ok && oe.Code() == code {
		return true
	}
	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		for _, e := range joined.Unwrap() {
			if hasLeafCode(e, code) {
				return true
			}
		}
	}
	if u, ok := err.(interface{ Unwrap() error }); ok {
		return hasLeafCode(u.Unwrap(), code)
	}
	return false
}

func TestServiceRenameHappyPath(t *testing.T) {
	now := time.Date(2026, 4, 10, 12, 0, 0, 0, time.UTC)
	svc, repo, store := newTestService(t, now)

	// seed
	createRes, err := svc.Create(context.Background(), organizationapp.CreateCommand{Name: "Old"})
	require.NoError(t, err)
	require.Len(t, store.rows, 1)

	_, err = svc.Rename(context.Background(), organizationapp.RenameCommand{
		ID:      createRes.ID,
		NewName: "New",
	})
	require.NoError(t, err)
	require.Equal(t, "New", repo.saved[createRes.ID].Name)
	require.Len(t, store.rows, 2)
	require.Equal(t, "medincident.orgstructure.v1.OrganizationRenamed", store.rows[1].EventType)
}

func TestServiceRenameRejectsEmpty(t *testing.T) {
	svc, _, store := newTestService(t, time.Now().UTC())
	createRes, err := svc.Create(context.Background(), organizationapp.CreateCommand{Name: "Acme"})
	require.NoError(t, err)
	rowsBefore := len(store.rows)

	_, err = svc.Rename(context.Background(), organizationapp.RenameCommand{
		ID:      createRes.ID,
		NewName: "   ",
	})
	require.Error(t, err)
	require.Len(t, store.rows, rowsBefore, "failed mutation must not append to outbox")
}

func TestServiceUpdateDescriptionHappyPath(t *testing.T) {
	svc, repo, store := newTestService(t, time.Now().UTC())
	createRes, err := svc.Create(context.Background(), organizationapp.CreateCommand{
		Name: "Acme", Description: "before",
	})
	require.NoError(t, err)

	_, err = svc.UpdateDescription(context.Background(), organizationapp.UpdateDescriptionCommand{
		ID:             createRes.ID,
		NewDescription: "after",
	})
	require.NoError(t, err)
	require.Equal(t, "after", repo.saved[createRes.ID].Description)
	require.Len(t, store.rows, 2)
	require.Equal(t, "medincident.orgstructure.v1.OrganizationDescriptionUpdated", store.rows[1].EventType)
}

func TestServiceRelocateLegalAddressToNilClears(t *testing.T) {
	svc, repo, store := newTestService(t, time.Now().UTC())
	createRes, err := svc.Create(context.Background(), organizationapp.CreateCommand{
		Name: "Acme",
		LegalAddress: &organizationapp.AddressInput{
			Text:  "Main St",
			Point: &organizationapp.PointInput{Longitude: 30.5, Latitude: 50.5},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, repo.saved[createRes.ID].LegalAddress)

	_, err = svc.RelocateLegalAddress(context.Background(), organizationapp.RelocateLegalAddressCommand{
		ID:      createRes.ID,
		Address: nil,
	})
	require.NoError(t, err)
	require.Nil(t, repo.saved[createRes.ID].LegalAddress)
	require.Len(t, store.rows, 2)
	require.Equal(t, "medincident.orgstructure.v1.OrganizationLegalAddressRelocated", store.rows[1].EventType)
}
