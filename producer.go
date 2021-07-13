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
	specVersion = "1.0"

	metricCountDelivered = "kafka_delivered"
	metricCountReceived = "kafka_received"

	signaledDateFormat = "2006-01-02"
)

type producer struct {
	config          *Config
	logChannels     telemetry.LogChannels
}



func (p *producer) start(ctx context.Context, ch <-chan *Message) {
	saramaConfig := kafkaConfig(p.config)

	producer, err := sarama.NewAsyncProducer([]string{p.config.Broker}, saramaConfig)
	if err != nil {
		p.logChannels.ErrorChan <- err
		return
	}
	defer producer.Close()

	input := producer.Input()

	go func(successes <-chan *sarama.ProducerMessage){
		for {
			success := <- successes
			_ = success
			p.logChannels.CountChan <- telemetry.Metric{
				Name:  metricCountDelivered,
				Value: 1,
				ConstLabels: map[string]string{
					"day": dayKey(time.Now()),
				},
			}
		}
	}(producer.Successes())

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
				"day": dayKey(time.Now()),
			},
		}
	}
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