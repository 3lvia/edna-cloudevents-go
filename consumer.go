package ednaevents

import (
	"context"
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
		c.logChannels.EventChan <- telemetry.Event{
			Name: "no_highwatermarks",
			Data: map[string]string{"topic": topic},
		}
	}

	for _, partition := range partitions {
		var offset int64 = 0
		if hw, ok := hwm[topic]; ok {
			if o, ok := hw[partition]; ok {
				offset = o
			}
		}
		go consumePartition(cnsmr, topic, partition, offset, ch, c.logChannels)
	}
}

func consumePartition(cons sarama.Consumer, topic string, partition int32, offset int64, ch chan<- *ConsumerEvent, logChannels telemetry.LogChannels) {
	logChannels.EventChan <- telemetry.Event{
		Name: "partition_consumer_start",
		Data: map[string]string{
			"partition": fmt.Sprintf("%d", partition),
			"offset": fmt.Sprintf("%d", offset),
		},
	}

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
	saramaConfig.Consumer.Offsets.AutoCommit.Enable = true

	consumer, err := sarama.NewConsumer([]string{c.config.Broker}, saramaConfig)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}