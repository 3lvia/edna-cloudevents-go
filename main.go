package ednaevents

import (
	"context"
	"github.com/prometheus/common/log"
)

func StartProducer(ctx context.Context, ch <-chan *ProducerEvent, opts ...Option) {
	collector := &OptionsCollector{}
	for _, opt := range opts {
		opt(collector)
	}

	c := collector.config
	err := c.loadProducer()
	if err != nil {
		log.Error(err)
		return
	}

	p := &producer{
		config:          c,
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
		log.Error(err)
		return
	}

	c := &consumer{
		config:      conf,
	}

	go c.start(ctx, ch)
}
