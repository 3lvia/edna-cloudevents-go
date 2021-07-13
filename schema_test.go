package ednaevents

import "testing"

func Test_getSchema(t *testing.T) {

	r := &schemaReaderImpl{
		endpointAddress: "https://psrc-dozoy.westeurope.azure.confluent.cloud",
		username:        "",
		password:        "",
	}

	s, err := r.getSchema("no.elvia.dp.securityauditlog")

	if err != nil {
		t.Fatal(err)
	}
	_ = s
}
