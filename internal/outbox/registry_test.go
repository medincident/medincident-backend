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

func TestRegistryRegisterAndLookupBothKeys(t *testing.T) {
	reg := outbox.NewRegistry()
	info := outbox.EventInfo{
		TypeName: "test.v1.A",
		Zero:     func() any { return &regStubA{} },
		ToProto:  func(any) (proto.Message, error) { return &emptypb.Empty{}, nil },
	}
	reg.Register(reflect.TypeOf(&regStubA{}), info)

	byType, ok := reg.ByGoType(reflect.TypeOf(&regStubA{}))
	require.True(t, ok)
	require.Equal(t, "test.v1.A", byType.TypeName)

	byName, ok := reg.ByTypeName("test.v1.A")
	require.True(t, ok)
	require.Equal(t, "test.v1.A", byName.TypeName)
}

func TestRegistryMissingLookup(t *testing.T) {
	reg := outbox.NewRegistry()
	_, ok := reg.ByGoType(reflect.TypeOf(&regStubA{}))
	require.False(t, ok)
	_, ok = reg.ByTypeName("nope")
	require.False(t, ok)
}

func TestRegistryPanicsOnDuplicateGoType(t *testing.T) {
	reg := outbox.NewRegistry()
	info := outbox.EventInfo{
		TypeName: "test.v1.A",
		Zero:     func() any { return &regStubA{} },
		ToProto:  func(any) (proto.Message, error) { return &emptypb.Empty{}, nil },
	}
	reg.Register(reflect.TypeOf(&regStubA{}), info)

	require.Panics(t, func() {
		reg.Register(reflect.TypeOf(&regStubA{}), info)
	})
}

func TestRegistryPanicsOnDuplicateTypeName(t *testing.T) {
	reg := outbox.NewRegistry()
	shared := outbox.EventInfo{
		TypeName: "test.v1.A",
		Zero:     func() any { return &regStubA{} },
		ToProto:  func(any) (proto.Message, error) { return &emptypb.Empty{}, nil },
	}
	reg.Register(reflect.TypeOf(&regStubA{}), shared)

	// Different Go type, same TypeName string.
	conflicting := outbox.EventInfo{
		TypeName: "test.v1.A",
		Zero:     func() any { return &regStubB{} },
		ToProto:  func(any) (proto.Message, error) { return &emptypb.Empty{}, nil },
	}
	require.Panics(t, func() {
		reg.Register(reflect.TypeOf(&regStubB{}), conflicting)
	})
}
