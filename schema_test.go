package ednaevents

import "testing"

func Test_getSchema(t *testing.T) {

	r := &schemaReaderImpl{
		endpointAddress: "https://psrc-lg26v.westeurope.azure.confluent.cloud",
		username:        "",
		password:        "",
	}
	
	schemaID := 100001
	s, _, err := r.getSchema(schemaID)

	if err != nil {
		t.Fatal(err)
	}
	_ = s
}
