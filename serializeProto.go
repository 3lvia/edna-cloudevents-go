package ednaevents

import "github.com/gogo/protobuf/proto"

type protoSerializer struct {

}

func (s *protoSerializer) SetSchema(config *SchemaConfig) error {
	return nil
}

func (s *protoSerializer) ContentType() string {
	return "application/protobuf"
}

func (s *protoSerializer) Serialize(input interface{}) (interface{}, error) {
	message := input.(proto.Message)

	b, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}

	return b, nil
}