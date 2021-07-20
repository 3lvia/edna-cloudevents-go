package ednaevents

import (
	"bytes"
	"context"
	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/prometheus/common/log"
	"time"
)

const (
	metricCountDelivered = "kafka_delivered"
	metricCountAttempted = "kafka_attempted"

	signaledDateFormat = "2006-01-02"
)

func NewSyncProducer(config *Config) (ProducerService, error) {
	p, err := syncProducer(config)
	if err != nil {
		return nil, err
	}

	pp := &producer{
		config:       config,
		syncProducer: p,
	}

	return pp, nil
}

type producer struct {
	config       *Config
	producer     sarama.AsyncProducer
	syncProducer sarama.SyncProducer
}

func (p *producer) Produce(e *ProducerEvent) (partition int32, offset int64, err error) {
	mBytes, key, headers, err := p.message(e)
	if err != nil {
		return 0, 0, err
	}

	m := &sarama.ProducerMessage{
		Topic:   p.config.Topic,
		Key:     sarama.StringEncoder(key),
		Value:   sarama.ByteEncoder(mBytes),
		Headers: headers,
	}

	return p.syncProducer.SendMessage(m)
}

func (p *producer) start(ctx context.Context, ch <-chan *ProducerEvent) {
	producer, err := p.asyncProducer()
	if err != nil {
		log.Error(err)
		return
	}
	defer producer.Close()

	go func(successes <-chan *sarama.ProducerMessage){
		for {
			success := <- successes
			_ = success
			metrics.incCounter(metricCountDelivered, map[string]string{"day": dayKey(time.Now().UTC())})
		}
	}(producer.Successes())

	input := producer.Input()
	for {
		obj := <-ch

		mBytes, key, headers, err := p.message(obj)
		if err != nil {
			log.Error(err)
			continue
		}

		m := &sarama.ProducerMessage{
			Topic:   p.config.Topic,
			Key:     sarama.StringEncoder(key),
			Value:   sarama.ByteEncoder(mBytes),
			Headers: headers,
		}

		input <- m

		metrics.incCounter(metricCountAttempted, map[string]string{"day": dayKey(time.Now().UTC())})
	}
}

func syncProducer(config *Config) (sarama.SyncProducer, error) {
	saramaConfig := kafkaConfig(config)
	saramaConfig.Producer.Return.Successes = true
	p, err := sarama.NewSyncProducer([]string{config.Broker}, saramaConfig)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *producer) asyncProducer() (sarama.AsyncProducer, error) {
	if p.producer != nil {
		return p.producer, nil
	}

	saramaConfig := kafkaConfig(p.config)
	saramaConfig.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer([]string{p.config.Broker}, saramaConfig)
	if err != nil {
		log.Error()
		return nil, err
	}

	return producer, nil
}

func (p *producer) message(m *ProducerEvent) ([]byte, string, []sarama.RecordHeader, error) {
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