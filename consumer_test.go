package ednaevents

import (
	"context"
	"github.com/Shopify/sarama"
	"os"
	"sync"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	ctx := context.Background()
	setEnv()

	c, err := ConfigFromEnvVars()
	if err != nil {
		t.Fatal(err)
	}

	events, err := StartDirectConsumer(ctx, WithConfig(c))

	go func(ch <-chan *ConsumerEvent) {
		for {
			event := <-ch
			t.Log(event)
		}
	}(events)

	<-time.After(1000 * time.Second)
}

func TestBasic(t *testing.T) {
	ctx := context.Background()
	_ = ctx
	setEnv()

	c, err := ConfigFromEnvVars()
	if err != nil {
		t.Fatal(err)
	}

	saramaConfig := kafkaConfig(c)

	cn, err := sarama.NewConsumer([]string{c.Broker}, saramaConfig)
	if err != nil {
		t.Fatal(err)
	}

	pc, err := cn.ConsumePartition(c.Topic, 0, sarama.OffsetOldest)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(messages <-chan *sarama.ConsumerMessage, errors <-chan *sarama.ConsumerError) {
		defer wg.Done()
		select {
		case msg := <-messages:
			t.Log(msg)
		case e := <-errors:
			t.Fatal(e)
		}
	}(pc.Messages(), pc.Errors())

	wg.Wait()
}

func setEnv() {
	os.Setenv("CLIENT_ID", "edna-searchindexer-hes-topology")
	os.Setenv("KAFKA_BROKER", "pkc-lq8gm.westeurope.azure.confluent.cloud:9092")
	os.Setenv("KAFKA_TOPIC", "public.msi.msim.meteringpoint_versions_latest_includeinvalid")
	os.Setenv("KAFKA_USER", "HD7H5I7PBPLM35W2")
	os.Setenv("KAFKA_PASSWORD", "VhHDh2m0lPBWa7oxLCgm5aSUcyLbu1Dz021ZJ+HJDApdQbSDBbvsx8FqR0rcqOex")
	os.Setenv("KAFKA_SCHEMA_ENDPOINT", "https://psrc-j39np.westeurope.azure.confluent.cloud")
	os.Setenv("KAFKA_SCHEMA_USER", "PYLKNYM44JYVIA65")
	os.Setenv("KAFKA_SCHEMA_PASSWORD", "6XPukppIAZIuL9zLb+zFUfdJtH3B9Y+MDRBdWa+DlO0E9N67oA+A5VoCDr/xEWny")
	os.Setenv("KAFKA_SCHEMA_ID", "public.msi.msim.meteringpoint_versions_latest_includeinvalid-value")
	os.Setenv("TYPE", "")
}
