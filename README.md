# edna-cloudevents-go
This is a library for producing and consuming events to Kafka in context of 'Elvia Databav' (EDNA). The library hides
the fact that cloudevents are used internally. In other words, as the user of this library, you provide events to be
queued through a channel, and the library handles the details of serializing the event and wrapping it in a cloud event
with all necessary attributes set and formatted correctly.

## Configuration
The library assumes that the following environment variables are set:
* **CLIENT_ID** is set as the client ID on the Kafka producer/consumer.
* **KAFKA_BROKER** is the address of the Kafka broker, often called bootstrap server.
* **KAFKA_TOPIC** the topic to produce/consumer to/from.
* **KAFKA_USER** the user used to authenticate.
* **KAFKA_PASSWORD** the password of the Kafka user.
* **KAFKA_SCHEMA_ENDPOINT** the address of the Kafka schema registry.
* **KAFKA_SCHEMA_USER** the user used to authenticate against the schema registry.
* **KAFKA_SCHEMA_PASSWORD** the password of the Kafka schema user.
* **KAFKA_SCHEMA_ID** the ID of the schema. Currently, must be an integer.
* **TYPE** the event type. Is used as subject to resolve the schema.  

### Configuration only for producers
* **EVENT_SOURCE** the source of events on the form http://elvia.no/[domain]/[system]
* **EVENT_TYPE** the type of the events on the form no.elvia.[domain].[type]

### Configuration only for direct consumers

### Configuration only for grouped consumers
* **GROUP_ID** the consumer group ID. All simultaneous services that have the same ID will automatically cooperate by load balancing consumptions from partitions

## Options

## Usage
### Producing events synchronously
When producing events synchronously, the client waits until the event is written to the Kafka queue before returning.
```
import ednaevents "github.com/3lvia/edna-cloudevents-go"

// Setup
config, err := ednaevents.ConfigFromEnvVars()
if err != nil {
    log.Fatal(err)
}

producer, err := ednaevents.NewSyncProducer(config)
if err != nil {
    log.Fatal(err)
}

// Producing
obj := getObjectToEnqueueFromSomewhere() // obj must satisfy ednaevents.Serializable
msg := &ednaevents.ProducerEvent{
		ID:       "",
		EntityID: "",
		Payload:  obj,
	}
partition, offset, err := producer.Produce(msg)
if err != nil {
    return err
}

log.Infof("written, partition: %d, offset: %d", partition, offset)
```
### Producing events synchronously
When producing events asynchronously, events to be written are communicated via channel.
```
import (
    context
    ednaevents "github.com/3lvia/edna-cloudevents-go"
)

ctx := context.Background()

config, err := edna.ConfigFromEnvVars()
if err != nil {
	log.Fatal(err)
}

producerChannel := make(chan *ednaevents.ProducerEvent)
// Start edna producer in separate go routine
edna.StartProducer(
	ctx,
	producerChannel,
	ednaevents.WithConfig(config))
	
// The channel producerChannel can now be used to send events to be enqueued.
```
### Consuming events directly
When consuming events "directly", a consumer group is not used. This means that all events that currently available
are received.
```
import (
    context
    ednaevents "github.com/3lvia/edna-cloudevents-go"
)

ctx := context.Background()

config, err := edna.ConfigFromEnvVars()
if err != nil {
	log.Fatal(err)
}

producedEventsChannel := make(chan *ednaevents.ConsumerEvent)
ednaevents.StartDirectConsumer(
    ctx,
    producedEventsChannel,
    ednaevents.WithConfig(config)) 
    
go func() {
    for {
        ev := <- producedEventsChannel
        // handle received event
    }
}()
```
### Consuming events with consumer group
When consuming events with a consumer group, the internal Kafka functionality handles and stores the current high
watermark offset. This means that only events that have not previously been seen are received upon startup.

In cases where the topic has multiple partitions, the internal Kafka functionality is able to balance consuming
automatically if more than one service with the same group ID is started.

```
import (
    context
    ednaevents "github.com/3lvia/edna-cloudevents-go"
)

ctx := context.Background()

config, err := edna.ConfigFromEnvVars()
if err != nil {
	log.Fatal(err)
}

producedEventsChannel := make(chan *ednaevents.ConsumerEvent)
ednaevents.StartGroupedConsumer(
    ctx,
    producedEventsChannel,
    ednaevents.WithConfig(config)) 
    
go func() {
    for {
        ev := <- producedEventsChannel
        // handle received event
    }
}()
```

## Metrics
### Metrics for producer
```
kafka_delivered{day="2021-07-11"} 479
```
### Metrics for consumer
```
# HELP kafka_received 
# TYPE kafka_received counter
kafka_received{day="2021-07-24",partition="0"} 75515
```