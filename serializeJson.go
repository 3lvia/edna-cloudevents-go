package ednaevents

import "encoding/json"

type jsonSerializer struct {

}

func (s *jsonSerializer) SetSchema(config *SchemaConfig) error {
	return nil
}

func (s *jsonSerializer) ContentType() string {
	return "application/json"
}

func (s *jsonSerializer) Serialize(input interface{}) (interface{}, error) {
	b, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	return b, err
}