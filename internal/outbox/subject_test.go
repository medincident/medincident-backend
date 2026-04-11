package outbox_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/medincident/medincident-command-service/internal/outbox"
)

func TestSubjectFromProto(t *testing.T) {
	cases := []struct {
		name string
		msg  proto.Message
		want string
	}{
		{
			name: "google.protobuf.Empty → empty",
			msg:  &emptypb.Empty{},
			want: "google.protobuf.empty",
		},
		{
			name: "google.protobuf.Timestamp → timestamp",
			msg:  &timestamppb.Timestamp{},
			want: "google.protobuf.timestamp",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, outbox.SubjectFromProto(tc.msg))
		})
	}
}
