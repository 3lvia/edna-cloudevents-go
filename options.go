package ednaevents

// Option wraps configuration information.
type Option func(collector *OptionsCollector)

// OptionsCollector collects options
type OptionsCollector struct {
	config       *Config
}

// WithConfig sets the configuration information.
func WithConfig(config *Config) Option {
	return func(collector *OptionsCollector) {
		collector.config = config
	}
}