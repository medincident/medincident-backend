package outbox_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
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

type fixedClock struct{ t time.Time }

func (c fixedClock) Now() time.Time { return c.t }

func registryWithStub(t *testing.T) outbox.Registry {
	t.Helper()
	reg := outbox.NewRegistry()
	reg.Register(reflect.TypeOf(&stubEvent{}), outbox.EventInfo{
		Subject: "test.v1.stub_event",
		ToProto: func(any) (proto.Message, error) { return &emptypb.Empty{}, nil },
	})
	return reg
}

const (
	stubAggIDA = "01910000-0000-0000-0000-000000000001"
	stubAggIDB = "01910000-0000-0000-0000-00000000000b"
)

func newPublisher(store outbox.Store, reg outbox.Registry) outbox.Publisher {
	return outbox.NewPublisher(store, reg, fixedClock{t: time.Unix(1_700_000_000, 0).UTC()})
}

func TestPublisherWritesAllEventsWithEnvelopeMetadata(t *testing.T) {
	reg := registryWithStub(t)
	store := &stubStore{}
	src := &stubSource{
		aggType: "stub",
		aggID:   stubAggIDA,
		events:  []any{&stubEvent{Label: "a"}, &stubEvent{Label: "b"}},
	}

	pub := newPublisher(store, reg)
	require.NoError(t, pub.Publish(context.Background(), &stubTx{}, src))
	require.Len(t, store.records, 2)

	for _, rec := range store.records {
		require.Equal(t, "stub", rec.AggregateType)
		require.Equal(t, stubAggIDA, rec.AggregateID)
		require.Equal(t, "test.v1.stub_event", rec.Subject)
		require.NotEqual(t, uuid.Nil, rec.EventID)
		require.Equal(t, rec.EventID.String(), rec.Headers["Nats-Msg-Id"])
		require.False(t, rec.CreatedAt.IsZero())
		require.Equal(t, rec.CreatedAt, rec.OccurredAt)

		var wrapped anypb.Any
		require.NoError(t, proto.Unmarshal(rec.Payload, &wrapped))
		require.Equal(t, "type.googleapis.com/google.protobuf.Empty", wrapped.TypeUrl)
	}
	require.Empty(t, src.events, "source should be drained after publish")
}

func TestPublisherReturnsErrorWhenMapperMissing(t *testing.T) {
	reg := outbox.NewRegistry() // no mappers registered
	store := &stubStore{}
	src := &stubSource{
		aggType: "stub",
		aggID:   stubAggIDA,
		events:  []any{&stubEvent{Label: "a"}},
	}

	err := newPublisher(store, reg).Publish(context.Background(), &stubTx{}, src)
	require.Error(t, err)
	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	require.Equal(t, outbox.ErrCodeNoMapper, oe.Code())
	require.Empty(t, store.records)
}

func TestPublisherWrapsStoreError(t *testing.T) {
	reg := registryWithStub(t)
	storeErr := errors.New("plain store failure")
	store := &stubStore{err: storeErr}
	src := &stubSource{
		aggType: "stub",
		aggID:   stubAggIDA,
		events:  []any{&stubEvent{Label: "a"}},
	}

	err := newPublisher(store, reg).Publish(context.Background(), &stubTx{}, src)
	require.Error(t, err)
	oe, ok := oops.AsOops(err)
	require.True(t, ok)
	require.Equal(t, outbox.ErrCodeStoreAppend, oe.Code())
	require.ErrorIs(t, err, storeErr)
}

func TestPublisherDoesNothingForEmptySource(t *testing.T) {
	reg := registryWithStub(t)
	store := &stubStore{}
	src := &stubSource{
		aggType: "stub",
		aggID:   stubAggIDA,
	}

	require.NoError(t, newPublisher(store, reg).Publish(context.Background(), &stubTx{}, src))
	require.Empty(t, store.records)
}

func TestPublisherMultipleSourcesKeepMetadataSeparate(t *testing.T) {
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

	require.NoError(t, newPublisher(store, reg).Publish(context.Background(), &stubTx{}, srcA, srcB))
	require.Len(t, store.records, 2)
	require.Equal(t, "a", store.records[0].AggregateType)
	require.Equal(t, "b", store.records[1].AggregateType)
}
