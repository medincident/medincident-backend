package outbox_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/medincident/medincident-command-service/internal/outbox"
)

type (
	regStubA struct{}
	regStubB struct{}
)

func TestRegistryRegisterAndLookup(t *testing.T) {
	reg := outbox.NewRegistry()
	info := outbox.EventInfo{
		Subject: "test.v1.a",
		ToProto: func(any) (proto.Message, error) { return &emptypb.Empty{}, nil },
	}
	reg.Register(reflect.TypeOf(&regStubA{}), info)

	got, ok := reg.Lookup(reflect.TypeOf(&regStubA{}))
	require.True(t, ok)
	require.Equal(t, "test.v1.a", got.Subject)
}

func TestRegistryMissingLookup(t *testing.T) {
	reg := outbox.NewRegistry()
	_, ok := reg.Lookup(reflect.TypeOf(&regStubA{}))
	require.False(t, ok)
}

func TestRegistryPanicsOnDuplicateGoType(t *testing.T) {
	reg := outbox.NewRegistry()
	info := outbox.EventInfo{
		Subject: "test.v1.a",
		ToProto: func(any) (proto.Message, error) { return &emptypb.Empty{}, nil },
	}
	reg.Register(reflect.TypeOf(&regStubA{}), info)

	require.Panics(t, func() {
		reg.Register(reflect.TypeOf(&regStubA{}), info)
	})
}

// Different Go types registered with different subjects do not conflict.
func TestRegistryAllowsDifferentTypesWithDifferentSubjects(t *testing.T) {
	reg := outbox.NewRegistry()
	reg.Register(reflect.TypeOf(&regStubA{}), outbox.EventInfo{
		Subject: "test.v1.a",
		ToProto: func(any) (proto.Message, error) { return &emptypb.Empty{}, nil },
	})
	require.NotPanics(t, func() {
		reg.Register(reflect.TypeOf(&regStubB{}), outbox.EventInfo{
			Subject: "test.v1.b",
			ToProto: func(any) (proto.Message, error) { return &emptypb.Empty{}, nil },
		})
	})
}
