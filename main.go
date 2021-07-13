package ednaevents

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

//import "github.com/confluentinc/confluent-kafka-go/kafka"

func StartProducer(ctx context.Context, ch <-chan *Message, opts ...Option) {
	collector := &OptionsCollector{}
	for _, opt := range opts {
		opt(collector)
	}

	c := collector.config
	err := c.loadProducer()
	if err != nil {
		collector.logChannels.ErrorChan <- err
		return
	}

	//sr := &schemaReaderImpl{
	//	endpointAddress: c.SchemaAPIEndpoint,
	//	username:        c.SchemaAPIUsername,
	//	password:        c.SchemaAPIPassword,
	//}


	p := &producer{
		serializer: collector.serializer,
		config:          c,
		logChannels:     collector.logChannels,
	}

	go p.start(ctx, ch)
}

func StartConsumer(ctx context.Context, ch chan<- cloudevents.Event, opts ...Option) {
	collector := &OptionsCollector{}
	for _, opt := range opts {
		opt(collector)
	}

	conf := collector.config
	err := conf.loadConsumer()
	if err != nil {
		collector.logChannels.ErrorChan <- err
		return
	}

	c := &consumer{
		config:      conf,
		logChannels: collector.logChannels,
	}

	go c.start(ctx, ch)
}
