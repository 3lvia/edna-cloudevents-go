package ednaevents

import "context"

//import "github.com/confluentinc/confluent-kafka-go/kafka"

func StartProducer(ctx context.Context, ch <-chan *Message, opts ...Option) {
	collector := &OptionsCollector{}
	for _, opt := range opts {
		opt(collector)
	}

	c := collector.config

	//sr := &schemaReaderImpl{
	//	endpointAddress: c.SchemaAPIEndpoint,
	//	username:        c.SchemaAPIUsername,
	//	password:        c.SchemaAPIPassword,
	//}


	p := &producer{
		schemaReference: "",
		config:          c,

		logChannels:     collector.logChannels,
	}

	go p.start(ctx, ch)
}

func StartConsumer(ch <-chan Message, opts ...Option) {
	collector := &OptionsCollector{}
	for _, opt := range opts {
		opt(collector)
	}

}
