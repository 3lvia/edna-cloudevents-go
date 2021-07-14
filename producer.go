package ednaevents

import (
	"bytes"
	"context"
	"github.com/3lvia/telemetry-go"
	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"time"
)

const (
	metricCountDelivered = "kafka_delivered"
	metricCountReceived = "kafka_received"

	signaledDateFormat = "2006-01-02"
)

type producer struct {
	config      *Config
	producer    sarama.AsyncProducer
	logChannels telemetry.LogChannels
}

func (p *producer) start(ctx context.Context, ch <-chan *Message) {
	producer, err := p.asyncProducer()
	if err != nil {
		p.logChannels.ErrorChan <- err
		return
	}
	defer producer.Close()

	go func(successes <-chan *sarama.ProducerMessage){
		for {
			success := <- successes
			_ = success
			p.logChannels.CountChan <- telemetry.Metric{
				Name:  metricCountDelivered,
				Value: 1,
				ConstLabels: map[string]string{
					"day": dayKey(time.Now().UTC()),
				},
			}
		}
	}(producer.Successes())

	input := producer.Input()
	for {
		obj := <-ch

		mBytes, key, headers, err := p.message(obj)
		if err != nil {
			p.logChannels.ErrorChan <- err
			continue
		}

		m := &sarama.ProducerMessage{
			Topic:   p.config.Topic,
			Key:     sarama.StringEncoder(key),
			Value:   sarama.ByteEncoder(mBytes),
			Headers: headers,
		}

		input <- m

		p.logChannels.CountChan <- telemetry.Metric{
			Name:  metricCountReceived,
			Value: 1,
			ConstLabels: map[string]string{
				"day": dayKey(time.Now().UTC()),
			},
		}
	}
}

func (p *producer) asyncProducer() (sarama.AsyncProducer, error) {
	if p.producer != nil {
		return p.producer, nil
	}

	saramaConfig := kafkaConfig(p.config)
	saramaConfig.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer([]string{p.config.Broker}, saramaConfig)
	if err != nil {
		p.logChannels.ErrorChan <- err
		return nil, err
	}

	return producer, nil
}

func (p *producer) message(m *Message) ([]byte, string, []sarama.RecordHeader, error) {
	var buf bytes.Buffer
	err := m.Payload.Serialize(&buf)
	if err != nil {
		return nil, "", nil, err
	}

	id := m.ID
	if id == "" {
		id = uuid.New().String()
	}

	var headers []sarama.RecordHeader
	headers = append(headers, sarama.RecordHeader{Key: []byte("id"), Value: []byte(id)})
	headers = append(headers, sarama.RecordHeader{Key: []byte("source"), Value: []byte(p.config.Source)})
	headers = append(headers, sarama.RecordHeader{Key: []byte("type"), Value: []byte(p.config.Type)})
	if m.EntityID != "" {
		headers = append(headers, sarama.RecordHeader{Key: []byte("subject"), Value: []byte(m.EntityID)})
	}

	return buf.Bytes(), m.EntityID, headers, nil
}

func dayKey(d time.Time) string {
	return d.Format(signaledDateFormat)
}

func kafkaConfig(config *Config) *sarama.Config {
	saramaConfig := sarama.NewConfig()

	saramaConfig.Version = sarama.V2_8_0_0
	saramaConfig.ClientID = config.ClientID

	saramaConfig.Net.SASL.Enable = true
	saramaConfig.Net.SASL.User = config.Username
	saramaConfig.Net.SASL.Password = config.Password
	saramaConfig.Net.SASL.Mechanism = "PLAIN"

	saramaConfig.Net.TLS.Enable = true

	//sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	return saramaConfig
}