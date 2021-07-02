package ednaevents

import (
	"fmt"
	"github.com/riferrei/srclient"
)

type schemaReader interface {
	getSchema(id int) (string, error)
}

type schemaReaderImpl struct {
	endpointAddress string
	username        string
	password        string
}

func (r *schemaReaderImpl) getSchema(id int) (string, string, error) {
	client := srclient.CreateSchemaRegistryClient(r.endpointAddress)
	client.SetCredentials(r.username, r.password)

	schema, err := client.GetSchema(id)
	if err != nil {
		return "", "", err
	}

	reference := fmt.Sprintf("%s/schemas/ids/%d", r.endpointAddress, id)

	return reference, schema.Schema(), nil
}