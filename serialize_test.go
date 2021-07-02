package ednaevents

import (
	"testing"
)

func Test_serialize_happyDays(t *testing.T) {
	// Arrange
	c := &customer{
		ID:          "123",
		Surname:     "Hansen",
		DisplayName: "Sigvald Hansen",
	}

	// Act
	b, err := serialize(c, schemaCustomer)

	// Assert
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if len(b) != 36 {
		t.Errorf("unexpected length of serialized bytes, got %d", len(b))
	}
}

func Test_serialize_nils(t *testing.T) {
	// Arrange
	c := &customer{
	}

	// Act
	b, err := serialize(c, schemaCustomer)

	// Assert
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if len(b) != 13 {
		t.Errorf("unexpected length of serialized bytes, got %d", len(b))
	}
}

type customer struct {
	ID          string `json:"customer_id"`
	Surname     string `json:"surname"`
	DisplayName string `json:"display_name"`
}

func (c *customer) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"customer_id": c.ID,
		"surname": c.Surname,
		"display_name": c.DisplayName,
	}
}

const (
	schemaCustomer = `{
  "doc": "Custumer.",
  "fields": [
    {
      "doc": "",
      "name": "customer_id",
      "type": "string"
    },
    {
      "doc": "",
      "name": "surname",
      "type": "string"
    },
    {
      "doc": "",
      "name": "display_name",
      "type": "string"
    }
  ],
  "name": "customer",
  "namespace": "no.elvia.kunde.salesforce",
  "type": "record"
}`
)
