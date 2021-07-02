package ednaevents

// Serializable is an object that is meant to be serialized using the Avro format. In order for this to work the
// object must be able to represent itself as a map.
type Serializable interface {
	// ToMap returns a map-based representation og the object
	ToMap() map[string]interface{}
}

// Deserializer is able to take a Avro deserialized map as input and convert it to an object.
type Deserializer func (m interface{}) (Serializable, error)

// Config contains information needed to configure both the cloud events and the connection to Kafka.
type Config struct {
	// Broker the Kafka broker, often also called "bootstrap servers"
	Broker string

	// Topic to produce or consume from
	Topic string

	// Username is the username used in order to authenticate against Confluent Kafka.
	Username string

	// Password is the password for the Confluent Kafka user.
	Password string

	// Source is the source domain that produced the message. This value is used to  populate the setting 'source' in
	// the cloudevents event. This value must be a valid URI. The general recommended pattern for Elvia is
	// http://elvia.no/[DOMAIN]/[SYSTEM]. Example: http://elvia.no/msi/msim
	Source string

	// Type of the entity of the message. This value is used to  populate the setting 'type' in
	// the cloudevents event. The general recommended pattern for Elvia is no.elvia.[DOMAIN].[TYPE]. Example:
	// no.elvia.msi.meteringpointversion
	Type string

	// SchemaAPIEndpoint is the http endpoint from which schema information can be fetched.
	SchemaAPIEndpoint string

	// SchemaAPIUsername is the username to be used when authenticating against the schema API.
	SchemaAPIUsername string

	// SchemaAPIPassword is the password for the user of the schema API endpoint.
	SchemaAPIPassword string
}

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
