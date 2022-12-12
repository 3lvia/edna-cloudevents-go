package ednaevents

import (
	"errors"
	"os"
)

func ConfigFromEnvVars() (*Config, error) {
	c := &Config{}
	if err := c.load(); err != nil {
		return nil, err
	}
	return c, nil
}

// Config contains information needed to configure both the cloud events and the connection to Kafka.
type Config struct {
	// ClientID is the unique ID of the client. This value should match the name of the deployment on Kubernetes
	// (if running on Kubernetes).
	ClientID string

	// Broker the Kafka broker, often also called "bootstrap servers"
	Broker string

	// Topic to produce or consume from
	Topic string

	// GroupID used by the consumer
	GroupID string

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

	// Verbose enables verbose logging for the underlying Kafka client.
	Verbose bool
}

// load loads this instance with values from environment variables.
func (c *Config) load() error {
	c.ClientID = os.Getenv("CLIENT_ID")
	if c.ClientID == "" {
		return errors.New("missing env var CLIENT_ID")
	}
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

	return nil
}

func (c *Config) loadProducer() error {
	c.Source = os.Getenv("EVENT_SOURCE")
	if c.Source == "" {
		return errors.New("missing env var EVENT_SOURCE")
	}
	c.Type = os.Getenv("EVENT_TYPE")
	if c.Type == "" {
		return errors.New("missing env var EVENT_TYPE")
	}
	return nil
}

func (c *Config) loadConsumer() error {
	return nil
}

func (c *Config) loadConsumerGroup() error {
	c.GroupID = os.Getenv("GROUP_ID")
	if c.GroupID == "" {
		return errors.New("missing env var GROUP_ID")
	}
	return nil
}
