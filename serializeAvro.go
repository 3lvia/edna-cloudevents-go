package ednaevents

import "github.com/linkedin/goavro/v2"

// https://psrc-dozoy.westeurope.azure.confluent.cloud/subjects/no.elvia.dp.securityauditlog
type avroSerializer struct {
	schema string
}

func (s *avroSerializer) SetSchema(config *SchemaConfig) error {
	reader := &schemaReaderImpl{
		endpointAddress: config.SchemaAPIEndpoint,
		username:        config.SchemaAPIUsername,
		password:        config.SchemaAPIPassword,
	}

	schema, err := reader.getSchema(config.Type)
	if err != nil {
		return err
	}

	s.schema = schema
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