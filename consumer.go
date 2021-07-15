package ednaevents

import (
	"context"
	"errors"
	"fmt"
	"github.com/3lvia/telemetry-go"
	"github.com/Shopify/sarama"
	"time"
)

const metricCountReceived = "kafka_received"

type consumer struct {
	config      *Config
	consumer    sarama.Consumer
	logChannels telemetry.LogChannels
}

func (c *consumer) start(ctx context.Context, ch chan<- *ConsumerEvent) {
	cnsmr, err := c.getConsumer()
	if err != nil {
		c.logChannels.ErrorChan <- err
		return
	}

	topic := c.config.Topic
	partitions, err := cnsmr.Partitions(topic)
	if err != nil {
		c.logChannels.ErrorChan <- err
		return
	}

	hwm := cnsmr.HighWaterMarks()
	if _, ok := hwm[topic]; !ok {
		c.logChannels.ErrorChan <- errors.New(fmt.Sprintf("no high water marks for topic %s", topic))
		return
	}

	for _, partition := range partitions {
		go consumePartition(cnsmr, topic, partition, hwm[topic][partition], ch, c.logChannels)
	}
}

func consumePartition(cons sarama.Consumer, topic string, partition int32, offset int64, ch chan<- *ConsumerEvent, logChannels telemetry.LogChannels) {
	pConsumer, err := cons.ConsumePartition(topic, partition, offset)
	if err != nil {
		logChannels.ErrorChan <- err
		return
	}

	messages := pConsumer.Messages()
	for {
		m := <- messages
		ch <- consumerEvent(m)

		logChannels.CountChan <- telemetry.Metric{
			Name:        metricCountReceived,
			Value:       1,
			ConstLabels: map[string]string{"day": dayKey(time.Now().UTC()), "partition": fmt.Sprintf("%d", partition)},
		}
	}
}

func consumerEvent(m *sarama.ConsumerMessage) *ConsumerEvent {
	headers := map[string]string{}
	for _, header := range m.Headers {
		headers[string(header.Key)] = string(header.Value)
	}

	return &ConsumerEvent{
		Value:        m.Value,
		Headers:      headers,
		Metadata: KafkaMetadata{
			Key:            m.Key,
			Topic:          m.Topic,
			Partition:      m.Partition,
			Offset:         m.Offset,
			Timestamp:      m.Timestamp,
			BlockTimestamp: m.BlockTimestamp,
		},
	}
}

func (c *consumer) getConsumer() (sarama.Consumer, error) {
	if c.consumer != nil {
		return c.consumer, nil
	}

	saramaConfig := kafkaConfig(c.config)
	consumer, err := sarama.NewConsumer([]string{c.config.Broker}, saramaConfig)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}