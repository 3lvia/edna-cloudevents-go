package ednaevents

func StartProducer(ch <-chan *Message, opts ...Option) {
	collector := &OptionsCollector{}
	for _, opt := range opts {
		opt(collector)
	}
}

func StartConsumer(ch <-chan Message, schemaID int, opts ...Option) {
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

	schemaReference, schema, err := sr.getSchema(schemaID)
	if err != nil {
		collector.logChannels.ErrorChan <- err
		return
	}

	h := &handler{
		schema:          schema,
		schemaReference: schemaReference,
		config:          c,
	}

	go h.start(ch)
}