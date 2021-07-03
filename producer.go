package ednaevents

import (
	"github.com/3lvia/telemetry-go"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
)

const (
	specVersion = "1.0"
	contentType = "application/avro"
)

type producer struct {
	produceChannel  chan<- *kafka.Message
	schema          string
	schemaReference string
	config          *Config
	logChannels     telemetry.LogChannels
}

func (p *producer) start(ch <-chan *Message) {
	for {
		obj := <-ch
		ce, err := p.getCloudEvent(obj)

		if err != nil {
			p.logChannels.ErrorChan <- err
			continue
		}

		js, err := ce.MarshalJSON()
		if err != nil {
			p.logChannels.ErrorChan <- err
			continue
		}

		p.produceChannel <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &p.config.Topic, Partition: kafka.PartitionAny}, Value: []byte(js)}
	}
}

func (p *producer) getCloudEvent(m *Message) (cloudevents.Event, error) {
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
	ce.SetSource(p.config.Source)
	ce.SetType(p.config.Type)
	ce.SetDataSchema(p.schemaReference)

	b, err := serialize(m.Payload, p.schema)
	if err != nil {
		return cloudevents.Event{}, err
	}

	ce.SetData(contentType, b)

	return ce, nil
}
