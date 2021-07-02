package ednaevents

import "testing"

func Test_handler_getCloudEvent(t *testing.T) {
	// Arrange
	config := &Config{
		Source:            "http://elvia.no/kunde/salesforce",
		Type:              "no.elvia.kunde.customer",
	}

	h := &handler{
		schema:          schemaCustomer,
		schemaReference: "http://schema.repo.com/kunde/customer",
		config:          config,
	}

	c := &customer{
		ID:          "123",
		Surname:     "Hansen",
		DisplayName: "Sigvald Hansen",
	}

	m := Message{
		EntityID: "123",
		Payload:  c,
	}

	// Act
	e, err := h.getCloudEvent(m)

	// Assert
	if err != nil {
		t.Errorf("unexpected error, %v", err)
	}
	if len(e.ID()) != 36 {
		t.Errorf("expected guid ID, got %s", e.ID())
	}
	if len(e.Data()) != 36 {
		t.Errorf("unexpected data, got bytes of length %d", len(e.Data()))
	}
	if e.Source() != "http://elvia.no/kunde/salesforce" {
		t.Errorf("unexpected source, got %s", e.Source())
	}
	if e.Type() != "no.elvia.kunde.customer" {
		t.Errorf("unexpected type, got %s", e.Type())
	}
	if e.Subject() != "123" {
		t.Errorf("unexpected subject, got %s", e.Subject())
	}
	if e.DataContentType() != "application/avro" {
		t.Errorf("unexpected data content type, got %s", e.DataContentType())
	}
	if e.DataSchema() != "http://schema.repo.com/kunde/customer" {
		t.Errorf("unexpected data schema, got %s", e.DataSchema())
	}
}
