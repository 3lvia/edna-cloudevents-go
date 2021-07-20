package ednaevents

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/prometheus/common/log"
	"time"
)

const metricCountReceived = "kafka_received"

type consumer struct {
	config      *Config
	consumer    sarama.Consumer
}

func (c *consumer) start(ctx context.Context, ch chan<- *ConsumerEvent) {
	cnsmr, err := c.getConsumer()
	if err != nil {
		log.Error(err)
		return
	}

	topic := c.config.Topic
	partitions, err := cnsmr.Partitions(topic)
	if err != nil {
		log.Error(err)
		return
	}

	hwm := cnsmr.HighWaterMarks()
	if _, ok := hwm[topic]; !ok {
		log.Infof("no_highwatermarks, topic: %s", topic)
	}

	for _, partition := range partitions {
		var offset int64 = 0
		if hw, ok := hwm[topic]; ok {
			if o, ok := hw[partition]; ok {
				offset = o
			}
		}
		go consumePartition(cnsmr, topic, partition, offset, ch)
	}
}

func consumePartition(cons sarama.Consumer, topic string, partition int32, offset int64, ch chan<- *ConsumerEvent) {
	log.Infof("partition_consumer_start, parition: %d, offset: %d", partition, offset)

	pConsumer, err := cons.ConsumePartition(topic, partition, offset)
	if err != nil {
		log.Error(err)
		return
	}

	messages := pConsumer.Messages()
	for {
		m := <- messages
		ch <- consumerEvent(m)

		labels := map[string]string{"day": dayKey(time.Now().UTC()), "partition": fmt.Sprintf("%d", partition)}
		metrics.incCounter(metricCountReceived, labels)
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