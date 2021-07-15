package ednaevents

import (
	"context"
	"encoding/json"
	"github.com/3lvia/telemetry-go"
	"github.com/Shopify/sarama"
	"io"
	"sync"
	"testing"

	"github.com/Shopify/sarama/mocks"
)

func Test_producer_message(t *testing.T) {
	// Arrange
	ctx := context.Background()
	logChannels := telemetry.Start(ctx, telemetry.Empty())
	p := &producer{
		config:      &Config{
			Source:            "http://company.com/kis",
			Type:              "com.company.customer.person",
		},
		logChannels: logChannels,
	}
	m := &ProducerEvent{
		EntityID: "hansen",
		Payload:  &person{
			Name: "Joe Hansen",
			Age:  51,
		},
	}

	// Act
	payload, entityID, headers, err := p.message(m)

	// Assert
	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}
	if entityID != "hansen"  {
		t.Errorf("unexpected entityID, got %s", entityID)
	}
	s := string(payload)
	if s != `{"name":"Joe Hansen","age":51}` {
		t.Errorf("unexpected payload, got %s", s)
	}
	if len(headers) != 4 {
		t.Errorf("unexpected number of headers, got %d", len(headers))
	}
}

func Test_producer_start(t *testing.T) {
	// Arrange
	ctx := context.Background()

	errorChan := make(chan string)
	teleChan := make(chan *telemetry.CapturedEvent)
	cpt := &teleCapture{received: teleChan}

	topic := "com.company.hr.person"

	logChannels := telemetry.Start(
		ctx,
		telemetry.Empty(),
		telemetry.WithCapture(cpt))

	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.ClientID = "client-1"
	config.Net.SASL.Enable = true
	config.Net.SASL.User = "user"
	config.Net.SASL.Password = "password"
	config.Net.SASL.Mechanism = "PLAIN"
	config.Net.TLS.Enable = true
	config.Producer.Return.Successes = true

	reporter := &errorReporter{}
	ap := mocks.NewAsyncProducer(reporter, config)

	var receivedMessages []*sarama.ProducerMessage

	wg := &sync.WaitGroup{}
	wg.Add(3)
	checker := func(m *sarama.ProducerMessage) error {
		receivedMessages = append(receivedMessages, m)
		wg.Done()
		return nil
	}
	ap.ExpectInputWithMessageCheckerFunctionAndSucceed(checker)

	ch := make(chan *ProducerEvent)

	p := &producer{
		config:      &Config{
			Topic:    topic,
			ClientID: "client-1",
			Source:   "http://company.com/kis",
			Type:     "com.company.hr.person",
		},
		producer:    ap,
		logChannels: logChannels,
	}

	// Act
	go p.start(ctx, ch)

	go func(mch chan<- *ProducerEvent) {
		mch <- &ProducerEvent{
			ID:       "m-1",
			EntityID: "p-1",
			Payload:  &person{
				Name: "Petter Jensen",
				Age:  33,
			},
		}
	}(ch)

	errorMessage := ""
	var capturedMetrics []*telemetry.CapturedEvent

	go func(tc <-chan *telemetry.CapturedEvent, ec <-chan string, wg *sync.WaitGroup) {
		for {
			select {
			case captured := <-tc:
				capturedMetrics = append(capturedMetrics, captured)
				wg.Done()
			case em := <-ec:
				errorMessage = em
				wg.Done()
			}
		}
	}(teleChan, errorChan, wg)

	wg.Wait()

	// Assert
	if errorMessage != "" {
		t.Errorf("unexpected error: %s", errorMessage)
	}
	if len(receivedMessages) != 1 {
		t.Errorf("expected 1 enqued message, got %d", len(receivedMessages))
	}
	if len(capturedMetrics) != 2 {
		t.Errorf("expected 2 captured metrics, got %d", len(capturedMetrics))
	}
}

type teleCapture struct {
	received chan<- *telemetry.CapturedEvent
}

func (c *teleCapture) Capture(ce *telemetry.CapturedEvent) {
	c.received <- ce
}

type errorReporter struct {
	received chan<- string
}

func (r *errorReporter) Errorf(f string, args ...interface{}) {
	r.received <- f
}

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (p *person) Serialize(w io.Writer) error{
	js, _ := json.Marshal(p)
	w.Write(js)
	return nil
}

type metricsCapture struct {
	ch      chan<- *telemetry.CapturedEvent
	errChan chan<- error
}

func (c *metricsCapture) Capture(e *telemetry.CapturedEvent) {
	if e.Type == "Error" {
		c.errChan <- e.Event.(error)
		return
	}
	c.ch <- e
}