package encoding

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// MarshalProtoMessageToJSON marshals a protobuf message to indented JSON.
func MarshalProtoMessageToJSON(msg proto.Message) ([]byte, error) {
	return protojson.MarshalOptions{
		Multiline:       true,
		Indent:          "  ",
		EmitUnpopulated: false,
		UseProtoNames:   true,
	}.Marshal(msg)
}
