package ednaevents

import "github.com/3lvia/telemetry-go"

// Option wraps configuration information.
type Option func(collector *OptionsCollector)

// OptionsCollector collects options
type OptionsCollector struct {
	config       *Config
	serializer   Serializer
	deserializer Deserializer
	logChannels  telemetry.LogChannels
}

// WithConfig sets the configuration information.
func WithConfig(config *Config) Option {
	return func(collector *OptionsCollector) {
		collector.config = config
	}
}

// WithLogChannels sets the logging infrastructure
func WithLogChannels(lc telemetry.LogChannels) Option {
	return func(collector *OptionsCollector) {
		collector.logChannels = lc
	}
}

// WithDeserializer sets
func WithDeserializer(d Deserializer) Option {
	return func(collector *OptionsCollector) {
		collector.deserializer = d
	}
}

func WithProtobuf() Option {
	return func(collector *OptionsCollector) {
		collector.serializer = &protoSerializer{}
	}
}

func WithAvro() Option {
	return func(collector *OptionsCollector) {
		collector.serializer = &avroSerializer{}
	}
}

func WithJson() Option {
	return func(collector *OptionsCollector) {
		collector.serializer = &jsonSerializer{}
	}
}

func WithSerializer(s Serializer) Option {
	return func(collector *OptionsCollector) {
		collector.serializer = s
	}
}