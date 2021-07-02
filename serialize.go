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
