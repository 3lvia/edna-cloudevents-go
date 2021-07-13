package ednaevents

import "github.com/3lvia/telemetry-go"

// Option wraps configuration information.
type Option func(collector *OptionsCollector)

// OptionsCollector collects options
type OptionsCollector struct {
	config       *Config
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