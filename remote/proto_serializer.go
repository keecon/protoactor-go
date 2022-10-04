package remote

import (
	"fmt"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type protoSerializer struct{}

func newProtoSerializer() *protoSerializer {
	return &protoSerializer{}
}

func (p *protoSerializer) Serialize(msg interface{}) ([]byte, error) {
	if message, ok := msg.(*status.Status); ok {
		msg = message.Proto()
	}

	if message, ok := msg.(proto.Message); ok {
		bytes, err := proto.Marshal(message)
		if err != nil {
			return nil, err
		}

		return bytes, nil
	}
	return nil, fmt.Errorf("msg must be proto.Message")
}

func (p *protoSerializer) Deserialize(typeName string, bytes []byte) (interface{}, error) {
	n, _ := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(typeName))

	pm := n.New().Interface()

	if err := proto.Unmarshal(bytes, pm); err != nil {
		return nil, err
	}

	if message, ok := pm.(*spb.Status); ok {
		return status.FromProto(message), nil
	}
	return pm, nil
}

func (protoSerializer) GetTypeName(msg interface{}) (string, error) {
	if message, ok := msg.(*status.Status); ok {
		msg = message.Proto()
	}

	if message, ok := msg.(proto.Message); ok {
		typeName := proto.MessageName(message)

		return string(typeName), nil
	}
	return "", fmt.Errorf("msg must be proto.Message")
}
