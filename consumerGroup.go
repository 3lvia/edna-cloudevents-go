package ednaevents

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
	"sync"
)

type groupedConsumer struct {
	config *Config
	ready  chan bool
	ch     chan<- *ConsumerEvent
}

func (p *groupedConsumer) init(ctx context.Context) (func(), error) {
	log.Print("kafka init...")

	err := p.config.loadConsumerGroup()
	if err != nil {
		return nil, err
	}

	config := kafkaConfig(p.config)

	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange // partition allocation strategy
	config.Consumer.Offsets.Initial = -2                                   // Where to start consumption when no group consumption displacement is found

	ctx, cancel := context.WithCancel(ctx)
	client, err := sarama.NewConsumerGroup([]string{p.config.Broker}, p.config.GroupID, config)
	if err != nil {
		return nil, err
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			//util.HandlePanic("client.Consume panic", log.StandardLogger())
		}()
		for {
			if err := client.Consume(ctx, []string{p.config.Topic}, p); err != nil {
				log.Fatalf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				log.Printf("kafka consumer stopped, %v", ctx.Err())
				return
			}
			p.ready = make(chan bool)
		}
	}()
	<-p.ready
	log.Printf("Sarama consumer up and running!...")
	// Ensure that the message in the channel is consumed when the system exits
	f := func() {
		log.Printf("kafka close")
		cancel()
		wg.Wait()
		if err = client.Close(); err != nil {
			log.Printf("kafka close error, %v", err)
		}
	}

	return f, nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (p *groupedConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(p.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (p *groupedConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (p *groupedConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	// Specific consumption news
	for message := range claim.Messages() {
		ce := consumerEvent(message)
		p.ch <- ce
		session.MarkMessage(message, "")

		//partition := claim.Partition()
		//labels := map[string]string{"day": dayKey(time.Now().UTC()), "partition": fmt.Sprintf("%d", partition)}
		//metrics.incCounter(metricCountReceived, labels)
	}
	return nil
}
