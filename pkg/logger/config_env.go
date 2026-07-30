package logger

import "os"

// Environment variable names read by ConfigFromEnv.
const (
	EnvLogLevel  = "LOG_LEVEL"
	EnvLogFormat = "LOG_FORMAT"
)

// ConfigFromEnv returns a Config seeded with the mandatory identity fields and
// with Level and Format read from the LOG_LEVEL and LOG_FORMAT environment
// variables. Unset or invalid values fall back to the production defaults
// (level info, format json). The returned Config can be further customized
// before being passed to New.
func ConfigFromEnv(service, environment, version string) Config {
	cfg := Config{
		Service:     service,
		Environment: environment,
		Version:     version,
		Level:       LevelInfo,
		Format:      FormatJSON,
	}

	if raw := os.Getenv(EnvLogLevel); raw != "" {
		if lvl, err := ParseLevel(raw); err == nil {
			cfg.Level = lvl
		}
	}

	switch Format(os.Getenv(EnvLogFormat)) {
	case FormatJSON:
		cfg.Format = FormatJSON
	case FormatLogfmt:
		cfg.Format = FormatLogfmt
	case FormatOTel:
		cfg.Format = FormatOTel
	}

	return cfg
}
