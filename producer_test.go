package ednaevents

import (
	"context"
	"encoding/json"
	"github.com/3lvia/telemetry-go"
	"io"
	"testing"
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
	m := &Message{
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

