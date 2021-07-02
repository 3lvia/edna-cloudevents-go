package ednaevents

import "testing"

func Test_getSchema(t *testing.T) {

	r := &schemaReaderImpl{
		endpointAddress: "https://psrc-lg26v.westeurope.azure.confluent.cloud",
		username:        "5FI655RUW7R6GT7J",
		password:        "Seyv4yjij0azUNynluUccLWA5N+aeuI0CRTe58sC7EmE9dRbBmizRTxybSORUrH0",
	}
	
	schemaID := 100001
	s, err := r.getSchema(schemaID)

	if err != nil {
		t.Fatal(err)
	}
	_ = s
}
