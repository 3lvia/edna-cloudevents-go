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
  
### Configuration only for producers
* **EVENT_SOURCE** the source of events on the form http://elvia.no/[domain]/[system]
* **EVENT_TYPE** the type of the events on the form no.elvia.[domain].[type]

## Options

## Usage

## Metrics
```
kafka_delivered{day="2021-07-11"} 479
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
```