package ednaevents

import "github.com/linkedin/goavro/v2"

type avroSerializer struct {
	schema string
}

func (s *avroSerializer) SetSchema(schemaReference string) error {
	return nil
}

func (s *avroSerializer) ContentType() string {
	return "application/avro"
}

func (s *avroSerializer) Serialize(input interface{}) (interface{}, error) {
	codec, err := goavro.NewCodec(s.schema)
	if err != nil {
		return nil, err
	}

	obj := input.(Serializable)
	m := obj.ToMap()

	b, err := codec.SingleFromNative(nil, m)
	if err != nil {
		return nil, err
	}

	return b, nil
}