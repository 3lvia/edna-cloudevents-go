package ednaevents

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

const (
	specVersion = "1.0"
	contentType = "application/avro"
)

type handler struct {
	schema          string
	schemaReference string
	config          *Config
}

func (h *handler) start(ch <-chan Message) {
	for {
		obj := <- ch
		ce := h.getCloudEvent(obj)

		_ = ce
	}
}

func (h *handler) getCloudEvent(m Message) cloudevents.Event {
	id := m.ID
	if id == "" {
		id = uuid.New().String()
	}

	ce := cloudevents.NewEvent()

	ce.SetID(id)

	if m.EntityID != "" {
		ce.SetSubject(m.EntityID)
	}

	ce.SetSpecVersion(specVersion)
	ce.SetSource(h.config.Source)
	ce.SetType(h.config.Type)
	ce.SetDataSchema(h.schemaReference)

	ce.SetData(contentType, nil)

	return ce
}