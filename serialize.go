package ednaevents

import "github.com/linkedin/goavro/v2"

func serialize(obj Serializable, schema string) ([]byte, error) {
	codec, err := goavro.NewCodec(schema)
	if err != nil {
		return nil, err
	}

	m := obj.ToMap()

	b, err := codec.SingleFromNative(nil, m)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func deserialize(b []byte, schema string, deserializer Deserializer) (Serializable, error) {
	codec, err := goavro.NewCodec(schema)
	if err != nil {
		return nil, err
	}

	obj, _, err := codec.NativeFromSingle(b)

	if err != nil {
		return nil, err
	}

	d, err := deserializer(obj)
	if err != nil {
		return nil, err
	}

	return d, nil
}
