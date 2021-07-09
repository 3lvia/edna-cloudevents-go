package ednaevents

import (
	"context"
	"encoding/json"
	"github.com/3lvia/telemetry-go"
	"github.com/Shopify/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	//"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"time"
)

const (
	specVersion = "1.0"
	contentType = "application/avro"

	metricCountDelivered = "kafka_delivered"
	metricCountUndelivered = "kafka_undelivered"
	metricCountIgnored = "kafka_ignored"

	signaledDateFormat = "2006-01-02"
)

type producer struct {
	//kConfig         *kafka.ConfigMap
	schema          string
	schemaReference string
	config          *Config
	logChannels     telemetry.LogChannels
}

func (p *producer) start(ctx context.Context, ch <-chan *Message) {
	//kProducer, err := kafka.NewProducer(p.kConfig)
	//defer kProducer.Close()
	//if err != nil {
	//	p.logChannels.ErrorChan <- err
	//	return
	//}

	//go func() {
	//	for e := range kProducer.Events() {
	//		switch ev := e.(type) {
	//		case *kafka.Message:
	//			m := ev
	//			if m.TopicPartition.Error != nil {
	//				p.logChannels.CountChan <- telemetry.Metric{
	//					Name:        metricCountUndelivered,
	//					Value:       1,
	//					ConstLabels: map[string]string{"day": dayKey(time.Now()), "topic": *m.TopicPartition.Topic},
	//				}
	//				log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	//			} else {
	//				p.logChannels.CountChan <- telemetry.Metric{
	//					Name:        metricCountDelivered,
	//					Value:       1,
	//					ConstLabels: map[string]string{
	//						"day": dayKey(time.Now()),
	//						"topic": *m.TopicPartition.Topic,
	//						"partition": fmt.Sprintf("%d", m.TopicPartition.Partition),
	//					},
	//				}
	//				log.Printf("Delivered message to topic %s [%d] at offset %v\n",
	//					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	//			}
	//			return
	//
	//		default:
	//			p.logChannels.CountChan <- telemetry.Metric{
	//				Name:        metricCountIgnored,
	//				Value:       1,
	//				ConstLabels: map[string]string{"day": dayKey(time.Now()), "topic": p.config.Topic},
	//			}
	//			log.Printf("Ignored event: %s\n", ev)
	//		}
	//	}
	//}()

	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_8_0_0
	saramaConfig.Net.SASL.User = p.config.Username
	saramaConfig.Net.SASL.Password = p.config.Password
	saramaConfig.Net.SASL.Mechanism = "PLAIN"

	//kConfig := &kafka.ConfigMap{
	//	"bootstrap.servers": c.Broker,
	//	"sasl.username":     c.Username,
	//	"sasl.password":     c.Password,
	//	"sasl.mechanism":    "PLAIN",
	//	"security.protocol": "SASL_SSL",
	//}

	sender, err := kafka_sarama.NewSender([]string{p.config.Broker}, saramaConfig, p.config.Topic)
	if err != nil {
		p.logChannels.ErrorChan <- err
		return
	}
	defer sender.Close(ctx)

	client, err := cloudevents.NewClient(sender, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		p.logChannels.ErrorChan <- err
		return
	}

	for {
		obj := <-ch
		ce, err := p.getCloudEvent(obj)

		if err != nil {
			p.logChannels.ErrorChan <- err
			continue
		}

		result := client.Send(ctx, ce)
		_ = result

		//js, err := ce.MarshalJSON()
		//if err != nil {
		//	p.logChannels.ErrorChan <- err
		//	continue
		//}
		//_ = js

		//kProducer.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &p.config.Topic, Partition: kafka.PartitionAny}, Value: js}
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

	//b, err := serialize(m.Payload, p.schema)
	//if err != nil {
	//	return cloudevents.Event{}, err
	//}

	b, err := json.Marshal(m.Payload)
	if err != nil {
		return ce, err
	}

	//ce.SetData(contentType, b)
	ce.SetData("application/json", b)

	return ce, nil
}

func dayKey(d time.Time) string {
	return d.Format(signaledDateFormat)
}