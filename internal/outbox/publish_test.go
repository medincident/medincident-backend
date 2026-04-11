package outbox_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/medincident/medincident-command-service/internal/outbox"
	"github.com/medincident/medincident-command-service/internal/shared/tx"
)

type stubEvent struct {
	Label string
}

type stubSource struct {
	aggType string
	aggID   string
	events  []any
}

func (s *stubSource) AggregateType() string { return s.aggType }
func (s *stubSource) AggregateID() string   { return s.aggID }
func (s *stubSource) PullEvents() []any {
	out := s.events
	s.events = nil
	return out
}

type stubStore struct {
	records []*outbox.Record
	err     error
}

func (s *stubStore) Append(_ context.Context, _ tx.Tx, r *outbox.Record) error {
	if r == nil {
		return errors.New("nil record")
	}
	if s.err != nil {
		return s.err
	}
	s.records = append(s.records, r)
	return nil
}

type stubTx struct{ tx.Tx }

func registryWithStub(t *testing.T) outbox.Registry {
	t.Helper()
	reg := outbox.NewRegistry()
	reg.Register(reflect.TypeOf(&stubEvent{}), outbox.EventInfo{
		TypeName: "test.v1.StubEvent",
		Zero:     func() any { return &stubEvent{} },
		ToProto:  func(any) (proto.Message, error) { return &emptypb.Empty{}, nil },
	})
	return reg
}

const (
	stubAggIDA = "01910000-0000-0000-0000-000000000001"
	stubAggIDB = "01910000-0000-0000-0000-00000000000b"
)

func TestPublishWritesAllEventsWithAggregateMetadata(t *testing.T) {
	reg := registryWithStub(t)
	store := &stubStore{}
	src := &stubSource{
		aggType: "stub",
		aggID:   stubAggIDA,
		events:  []any{&stubEvent{Label: "a"}, &stubEvent{Label: "b"}},
	}

	err := outbox.Publish(context.Background(), &stubTx{}, store, reg, src)
	require.NoError(t, err)
	require.Len(t, store.records, 2)

	for _, rec := range store.records {
		require.Equal(t, "stub", rec.AggregateType)
		require.Equal(t, stubAggIDA, rec.AggregateID)
		require.Equal(t, "test.v1.StubEvent", rec.EventType)
		require.NotEqual(t, uuid.Nil, rec.ID)
		require.False(t, rec.CreatedAt.IsZero())
	}
	require.JSONEq(t, `{"Label":"a"}`, string(store.records[0].Payload))
	require.JSONEq(t, `{"Label":"b"}`, string(store.records[1].Payload))
	require.Empty(t, src.events, "source should be drained after publish")
}

func TestPublishReturnsErrorWhenMapperMissing(t *testing.T) {
	reg := outbox.NewRegistry() // no mappers registered
	store := &stubStore{}
	src := &stubSource{
		aggType: "stub",
		aggID:   stubAggIDA,
		events:  []any{&stubEvent{Label: "a"}},
	}

	err := outbox.Publish(context.Background(), &stubTx{}, store, reg, src)
	require.Error(t, err)
	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	require.Equal(t, outbox.CodeNoMapper, oe.Code())
	require.Empty(t, store.records)
}

func TestPublishWrapsStoreError(t *testing.T) {
	reg := registryWithStub(t)
	storeErr := errors.New("plain store failure") // simulate a non-oops driver error
	store := &stubStore{err: storeErr}
	src := &stubSource{
		aggType: "stub",
		aggID:   stubAggIDA,
		events:  []any{&stubEvent{Label: "a"}},
	}

	err := outbox.Publish(context.Background(), &stubTx{}, store, reg, src)
	require.Error(t, err)
	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	require.Equal(t, outbox.CodeStoreAppendFailed, oe.Code())
	require.ErrorIs(t, err, storeErr)
}

func TestPublishDoesNothingForEmptySource(t *testing.T) {
	reg := registryWithStub(t)
	store := &stubStore{}
	src := &stubSource{
		aggType: "stub",
		aggID:   stubAggIDA,
	}

	err := outbox.Publish(context.Background(), &stubTx{}, store, reg, src)
	require.NoError(t, err)
	require.Empty(t, store.records)
}

func TestPublishMultipleSourcesKeepsTheirMetadataSeparate(t *testing.T) {
	reg := registryWithStub(t)
	store := &stubStore{}
	srcA := &stubSource{
		aggType: "a",
		aggID:   stubAggIDA,
		events:  []any{&stubEvent{Label: "from a"}},
	}
	srcB := &stubSource{
		aggType: "b",
		aggID:   stubAggIDB,
		events:  []any{&stubEvent{Label: "from b"}},
	}

	err := outbox.Publish(context.Background(), &stubTx{}, store, reg, srcA, srcB)
	require.NoError(t, err)
	require.Len(t, store.records, 2)
	require.Equal(t, "a", store.records[0].AggregateType)
	require.Equal(t, "b", store.records[1].AggregateType)
}
