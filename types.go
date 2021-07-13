package ednaevents

import "io"

// Serializable is an object that is meant to be serialized using the Avro format. In order for this to work the
// object must be able to represent itself as a map.
type Serializable interface {
	// Serialize serializes the object as bytes onto the writer
	Serialize(w io.Writer) error
}

type Serializer interface {
	SetSchema(config *SchemaConfig) error

	ContentType() string

	Serialize(input interface{}) (interface{}, error)
}

// Deserializer is able to take a Avro deserialized map as input and convert it to an object.
type Deserializer func (m interface{}) (Serializable, error)

// Message wraps an object to be put on a topic or consumed from it.
type Message struct {
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
