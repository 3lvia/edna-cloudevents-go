package ednaevents

import (
	"io"
	"time"
)

// Serializable is an object that is meant to be serialized using the Avro format. In order for this to work the
// object must be able to represent itself as a map.
type Serializable interface {
	// Serialize serializes the object as bytes onto the writer
	Serialize(w io.Writer) error
}

// ProducerEvent wraps an object to be put on a topic or consumed from it.
type ProducerEvent struct {
	// ID is the id of the message. This value is used to populate the property 'id' of the cloudevents event. If this
	// value is not set a GUID-based value will be generated and set  internally. This value should only be used in the
	// (rare) case where the producer needs control over the actual id
	// message ID.
	ID string

	// EntityID is the ID of the entity of the message. This value can be omitted for time series events, but must be
	// included for entity events. This value is used both as the property 'subject' of the cloudevents event as well
	// as the key of the Kafka event.
	EntityID string

	// Payload is the actual entity- or time series event to be sent. This object will be serialized to using Avro and
	// wrapped in a cloudevent before being queued.
	Payload Serializable
}

// ConsumerEvent wraps an object that was consumed from a topic.
type ConsumerEvent struct {
	// Value is the payload of the actual message as it was put onto Kafka by the producer.
	Value []byte

	// Headers are the headers of the message as it was put onto Kafka by the producer.
	Headers map[string]string

	// Metadata contains Kafka-related metadata of the message.
	Metadata KafkaMetadata
}

// KafkaMetadata contains Kafka-related metadata of a message.
type KafkaMetadata struct {
	// Key is the key of the message as it was put onto Kafka by the producer. If the topic is compacted, messages are
	// grouped by key and only the last message for a given key is kept. This value is used as the property 'subject' of
	// the cloudevents event.
	Key []byte

	// Topic is the topic the message was dequeued from.
	Topic string

	// Partition is the partition the message was dequeued from.
	Partition int32

	// Offset is the offset of the message within the partition.
	Offset int64

	// Timestamp is the time when the message was put onto Kafka by the producer.
	Timestamp time.Time

	// BlockTimestamp is the time when the message was dequeued from Kafka.
	BlockTimestamp time.Time
}

// ProducerService is the interface for a service that can produce events.
type ProducerService interface {
	// Produce produces an event on the given topic.
	Produce(m *ProducerEvent) (partition int32, offset int64, err error)
}
