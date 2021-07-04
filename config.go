package ednaevents

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func ConfigFromEnvVars() (*Config, error) {
	c := &Config{}
	err := c.load()
	if err != nil {
		return nil, err
	}
	return c, nil
}

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

	// SchemaID the id of the Avro schema
	SchemaID int
}

// load loads this instance with values from environment variables.
func (c *Config) load() error {
	c.Broker = os.Getenv("KAFKA_BROKER")
	if c.Broker == "" {
		return errors.New("missing env var KAFKA_BROKER")
	}
	c.Topic = os.Getenv("KAFKA_TOPIC")
	if c.Topic == "" {
		return errors.New("missing env var KAFKA_TOPIC")
	}
	c.Username = os.Getenv("KAFKA_USER")
	if c.Username == "" {
		return errors.New("missing env var KAFKA_USER")
	}
	c.Password = os.Getenv("KAFKA_PASSWORD")
	if c.Password == "" {
		return errors.New("missing env var KAFKA_PASSWORD")
	}
	c.Source = os.Getenv("EVENT_SOURCE")
	if c.Source == "" {
		return errors.New("missing env var EVENT_SOURCE")
	}
	c.Type = os.Getenv("EVENT_TYPE")
	if c.Type == "" {
		return errors.New("missing env var EVENT_TYPE")
	}
	c.SchemaAPIEndpoint = os.Getenv("KAFKA_SCHEMA_ENDPOINT")
	if c.SchemaAPIEndpoint == "" {
		return errors.New("missing env var KAFKA_SCHEMA_ENDPOINT")
	}
	c.SchemaAPIUsername = os.Getenv("KAFKA_SCHEMA_USER")
	if c.SchemaAPIUsername == "" {
		return errors.New("missing env var KAFKA_SCHEMA_USER")
	}
	c.SchemaAPIPassword = os.Getenv("KAFKA_SCHEMA_PASSWORD")
	if c.SchemaAPIPassword == "" {
		return errors.New("missing env var KAFKA_SCHEMA_PASSWORD")
	}
	sid := os.Getenv("KAFKA_SCHEMA_ID")
	if sid == "" {
		 return errors.New("missing env var KAFKA_SCHEMA_ID")
	}
	siid, err := strconv.Atoi(sid)
	if err != nil {
		return errors.New(fmt.Sprintf("kafka schema id %s could not be converted to an int", sid))
	}
	c.SchemaID = siid
	return nil
}
