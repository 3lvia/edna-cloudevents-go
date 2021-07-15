package ednaevents

import (
	"context"
)

func StartProducer(ctx context.Context, ch <-chan *ProducerEvent, opts ...Option) {
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

	p := &producer{
		config:          c,
		logChannels:     collector.logChannels,
	}

	go p.start(ctx, ch)
}

func StartConsumer(ctx context.Context, ch chan<- *ConsumerEvent, opts ...Option) {
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
