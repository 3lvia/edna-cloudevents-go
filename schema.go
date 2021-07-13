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

func (r *schemaReaderImpl) getSchema(typ string) (string, error) {
	client := srclient.CreateSchemaRegistryClient(r.endpointAddress)
	client.SetCredentials(r.username, r.password)

	subject := fmt.Sprintf("%s-value", typ)

	schema, err := client.GetLatestSchema(subject, false)
	if err != nil {
		return "", err
	}

	return schema.Schema(), nil
}