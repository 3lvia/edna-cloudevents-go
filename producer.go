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
	kConfig         *kafka.ConfigMap
	schema          string
	schemaReference string
	config          *Config
	logChannels     telemetry.LogChannels
}

func (p *producer) start(ch <-chan *Message) {
	kProducer, err := kafka.NewProducer(kConfig)
	defer kProducer.Close()
	if err != nil {
		p.logChannels.ErrorChan <- err
		return
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
				return

			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

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

		kProducer.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &p.config.Topic, Partition: kafka.PartitionAny}, Value: []byte(js)}
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
