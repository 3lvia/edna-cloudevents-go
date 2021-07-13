package ednaevents

import (
	"context"
	"github.com/3lvia/telemetry-go"
	"sync"
	"testing"
)

func Test_handler_getCloudEvent(t *testing.T) {
	//// Arrange
	//config := &Config{
	//	Source:            "http://elvia.no/kunde/salesforce",
	//	Type:              "no.elvia.kunde.customer",
	//}
	//
	//h := &producer{
	//	schema:          schemaCustomer,
	//	schemaReference: "http://schema.repo.com/kunde/customer",
	//	config:          config,
	//}
	//
	//c := &customer{
	//	ID:          "123",
	//	Surname:     "Hansen",
	//	DisplayName: "Sigvald Hansen",
	//}
	//
	//m := Message{
	//	EntityID: "123",
	//	Payload:  c,
	//}
	//
	//// Act
	//e, err := h.getCloudEvent(m)
	//
	//// Assert
	//if err != nil {
	//	t.Errorf("unexpected error, %v", err)
	//}
	//if len(e.ID()) != 36 {
	//	t.Errorf("expected guid ID, got %s", e.ID())
	//}
	//if len(e.Data()) != 36 {
	//	t.Errorf("unexpected data, got bytes of length %d", len(e.Data()))
	//}
	//if e.Source() != "http://elvia.no/kunde/salesforce" {
	//	t.Errorf("unexpected source, got %s", e.Source())
	//}
	//if e.Type() != "no.elvia.kunde.customer" {
	//	t.Errorf("unexpected type, got %s", e.Type())
	//}
	//if e.Subject() != "123" {
	//	t.Errorf("unexpected subject, got %s", e.Subject())
	//}
	//if e.DataContentType() != "application/avro" {
	//	t.Errorf("unexpected data content type, got %s", e.DataContentType())
	//}
	//if e.DataSchema() != "http://schema.repo.com/kunde/customer" {
	//	t.Errorf("unexpected data schema, got %s", e.DataSchema())
	//}
}

func Test_producer_start(t *testing.T) {
	// Arrange
	ctx := context.Background()

	cch := make(chan *telemetry.CapturedEvent)
	errChan := make(chan error)
	capture := &metricsCapture{ch: cch, errChan: errChan}

	logChannels := telemetry.Start(ctx, telemetry.Empty(), telemetry.WithCapture(capture))

	config := &Config{
		ClientID:          "core.vault-audit-processor",
		Broker:            "pkc-lq8gm.westeurope.azure.confluent.cloud:9092",
		Topic:             "heitmann.test",
		Username:          "USER",
		Password:          "PASSWORD",
		Source:            "http://elvia.no/dp/vault",
		Type:              "no.elvia.dp.auditlogentry",
		SchemaAPIEndpoint: "",
		SchemaAPIUsername: "",
		SchemaAPIPassword: "",
	}

	p := &producer{config: config, logChannels: logChannels}

	ch := make(chan *Message)

	var capturedError error

	// Act
	go p.start(ctx, ch)

	go func(channel chan<- *Message){
		pp := &person{
			Name: "Jens",
			Age:  36,
		}
		msg := &Message{
			ID:       "",
			EntityID: "",
			Payload:  pp,
		}
		channel <- msg
	}(ch)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		select {
		case ee := <-errChan:
			capturedError = ee
			wg.Done()
		case <-cch:
			wg.Done()
		}
	}()

	wg.Wait()

	// Arrange
	if capturedError != nil {
		t.Error(capturedError)
	}
}

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (p *person) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name": p.Name,
		"age": p.Age,
	}
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