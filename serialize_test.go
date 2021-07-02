package ednaevents

import (
	"encoding/json"
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

func Test_serialize_deserialize(t *testing.T) {
	// Arrange
	c := &customer{
		ID:          "123",
		Surname:     "Hansen",
		DisplayName: "Sigvald Hansen",
	}

	// Act
	b, err := serialize(c, schemaCustomer)
	if err != nil {
		t.Errorf("unexpected error when serializing %v", err)
	}

	d := func(m interface{}) (Serializable, error) {
		b, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}

		var c customer
		err = json.Unmarshal(b, &c)
		if err != nil {
			return nil, err
		}

		return &c, nil
	}

	obj, err := deserialize(b, schemaCustomer, d)
	if err != nil {
		t.Errorf("unexpected error when deserializing %v", err)
	}

	// Assert
	cc := obj.(*customer)
	if cc.ID != c.ID {
		t.Errorf("unexpected ID, got %s", cc.ID)
	}
	if cc.Surname != c.Surname {
		t.Errorf("unexpected Surname, got %s", cc.ID)
	}
	if cc.DisplayName != c.DisplayName {
		t.Errorf("unexpected DisplayName, got %s", cc.ID)
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
