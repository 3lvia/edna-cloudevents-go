package ednaevents

import "context"

//import "github.com/confluentinc/confluent-kafka-go/kafka"

func StartProducer(ctx context.Context, ch <-chan *Message, opts ...Option) {
	collector := &OptionsCollector{}
	for _, opt := range opts {
		opt(collector)
	}

	c := collector.config

	sr := &schemaReaderImpl{
		endpointAddress: c.SchemaAPIEndpoint,
		username:        c.SchemaAPIUsername,
		password:        c.SchemaAPIPassword,
	}

	schemaReference, schema, err := sr.getSchema(c.SchemaID)
	if err != nil {
		collector.logChannels.ErrorChan <- err
		return
	}

	//kConfig := &kafka.ConfigMap{
	//	"bootstrap.servers": c.Broker,
	//	"sasl.username":     c.Username,
	//	"sasl.password":     c.Password,
	//	"sasl.mechanism":    "PLAIN",
	//	"security.protocol": "SASL_SSL",
	//}

	p := &producer{
		//kConfig:         kConfig,
		schema:          schema,
		schemaReference: schemaReference,
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
