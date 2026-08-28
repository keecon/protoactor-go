package remote

import (
	"testing"

	"github.com/keecon/protoactor-go/actor"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//func TestJsonSerializer_round_trip(t *testing.T) {
//	m := &ActorPidRequest{
//		Kind: "abc",
//		Name: "def",
//	}
//	b, typeName, _ := Serialize(m, 1)
//	res, err := Deserialize(b, typeName, 1)
//
//	assert.Nil(t, err)
//
//	var typed = res.(*ActorPidRequest)
//	assert.Equal(t, "remote.ActorPidRequest", typeName)
//	assert.Equal(t, m, typed)
//}
//
//func TestJsonSerializer_Serialize_PID_raw(t *testing.T) {
//	system := actor.NewActorSystem()
//	m, _ := system.Root.SpawnNamed(actor.PropsFromFunc(func(ctx actor.Context) {}), "actorpid")
//	var ser = jsonpb.Marshaler{}
//	res, _ := ser.MarshalToString(m)
//	assert.Equal(t, "{\"Address\":\"nonhost\",\"Id\":\"actorpid\"}", res)
//}
//
//func TestJsonSerializer_Serialize_PID(t *testing.T) {
//	system := actor.NewActorSystem()
//	m := system.NewLocalPID("foo")
//	b, typeName, _ := Serialize(m, 1)
//	res, err := Deserialize(b, typeName, 1)
//
//	assert.Nil(t, err)
//
//	var typed = res.(*actor.PID)
//	assert.Equal(t, "actor.PID", typeName)
//	assert.Equal(t, m, typed)
//}

func TestProtobufSerializer_Serialize_PID(t *testing.T) {
	system := actor.NewActorSystem()
	m := system.NewLocalPID("foo")
	b, typeName, _ := Serialize(m, 0)
	res, err := Deserialize(b, typeName, 0)

	assert.Nil(t, err)

	typed := res.(*actor.PID)
	assert.Equal(t, "actor.PID", typeName)
	assert.True(t, m.Equal(typed))
}

func TestProtobufSerializerRoundTripsGRPCStatus(t *testing.T) {
	serializer := newProtoSerializer()
	want := status.New(codes.InvalidArgument, "invalid request")

	payload, err := serializer.Serialize(want)
	if err != nil {
		t.Fatalf("serialize gRPC status: %v", err)
	}
	typeName, err := serializer.GetTypeName(want)
	if err != nil {
		t.Fatalf("get gRPC status type name: %v", err)
	}
	if typeName != "google.rpc.Status" {
		t.Fatalf("type name = %q, want %q", typeName, "google.rpc.Status")
	}

	value, err := serializer.Deserialize(typeName, payload)
	if err != nil {
		t.Fatalf("deserialize gRPC status: %v", err)
	}
	got, ok := value.(*status.Status)
	if !ok {
		t.Fatalf("deserialized type = %T, want *status.Status", value)
	}
	if got.Code() != want.Code() || got.Message() != want.Message() {
		t.Fatalf("deserialized status = (%s, %q), want (%s, %q)", got.Code(), got.Message(), want.Code(), want.Message())
	}
}

func TestSerialize_InvalidSerializerID(t *testing.T) {
	_, _, err := Serialize("msg", int32(len(serializers)))
	assert.Error(t, err)
}

func TestDeserialize_InvalidSerializerID(t *testing.T) {
	_, err := Deserialize([]byte("{}"), "", int32(len(serializers)))
	assert.Error(t, err)
}

// TestProtobufSerializer_Deserialize_InvalidType ensures an error is returned when the message type is unknown.
func TestProtobufSerializer_Deserialize_InvalidType(t *testing.T) {
	_, err := Deserialize([]byte{}, "unknown.Type", 0)
	assert.Error(t, err)
}
