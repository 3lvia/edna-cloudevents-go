package ednaevents

import (
	"context"
	"github.com/prometheus/common/log"
	"os"
	"os/signal"
	"syscall"
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

func StartGroupedConsumer(ctx context.Context, ch chan<- *ConsumerEvent, opts ...Option) {
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
		log.Warnln("terminating: via signal")
	}
	f()
}

// StartDirectConsumer starts a consumer without support for consumer groups or offsets, meaning that all events
// fra the earliest existing offset are returned.
func StartDirectConsumer(ctx context.Context, ch chan<- *ConsumerEvent, opts ...Option) {
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
