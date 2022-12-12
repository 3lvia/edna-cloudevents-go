// Package ednaevents provides a simple interface for producing and consuming events from Kafka.
package ednaevents

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// StartProducer starts a producer that can be used to send events to Kafka.
func StartProducer(ctx context.Context, ch <-chan *ProducerEvent, opts ...Option) {
	collector := &OptionsCollector{}
	for _, opt := range opts {
		opt(collector)
	}

	c := collector.config
	err := c.loadProducer()
	if err != nil {
		log.Printf("error loading producer: %v", err)
		return
	}

	p := &producer{
		config: c,
	}

	go p.start(ctx, ch)
}

// StartGroupedConsumer starts a consumer that supports consumer groups and offsets, meaning that the central Kafka
// cluster will keep track of which events have been consumed and which have not.
func StartGroupedConsumer(ctx context.Context, opts ...Option) <-chan *ConsumerEvent {
	channel := make(chan *ConsumerEvent)
	go func(ch chan<- *ConsumerEvent) {
		collector := &OptionsCollector{}
		for _, opt := range opts {
			opt(collector)
		}

		k := groupedConsumer{
			config: collector.config,
			ready:  make(chan bool),
			ch:     ch,
		}
		f, err := k.init(ctx)
		if err != nil {
			log.Fatal(err)
		}

		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-sigterm:
			log.Print("terminating: via signal")
		}

		f()
	}(channel)

	return channel
}

// StartDirectConsumer starts a consumer without support for consumer groups or offsets, meaning that all events
// from the earliest existing offset are returned.
func StartDirectConsumer(ctx context.Context, opts ...Option) (<-chan *ConsumerEvent, error) {
	collector := &OptionsCollector{}
	for _, opt := range opts {
		opt(collector)
	}

	conf := collector.config
	err := conf.loadConsumer()
	if err != nil {
		return nil, err
	}

	c := &consumer{
		config: conf,
	}

	ch := make(chan *ConsumerEvent)

	go c.start(ctx, ch)

	return ch, nil
}
