package ednaevents

import (
	"context"
	"github.com/3lvia/telemetry-go"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type consumer struct {
	config      *Config
	logChannels telemetry.LogChannels
}

func (c *consumer) start(ctx context.Context, ch chan<- cloudevents.Event) {
	saramaConfig := kafkaConfig(c.config)
	consumer, err := kafka_sarama.NewConsumer([]string{c.config.Broker}, saramaConfig, "", c.config.Topic)
	if err != nil {
		c.logChannels.ErrorChan <- err
		return
	}

	client, err := cloudevents.NewClient(consumer, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		c.logChannels.ErrorChan <- err
		return
	}

	receive := func(ctx context.Context, event cloudevents.Event) {
		ch <- event
	}

	client.StartReceiver(ctx, receive)
}