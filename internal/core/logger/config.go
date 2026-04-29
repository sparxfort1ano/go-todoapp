package logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// config holds the environment-based settings required to configure
// the Zap logger, specifically the minimum log level 
// and the output directory for log files.
type config struct {
	Level  string `envconfig:"LEVEL"  required:"true"`
	Folder string `envconfig:"FOLDER" required:"true"`
}

// newConfig parses the system environment variables with the "LOGGER_" prefix.
// into the config struct. It returns an error if required variables are missing or malformed.
func newConfig() (config, error) {
	var cfg config

	if err := envconfig.Process("LOGGER", &cfg); err != nil {
		return config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return cfg, nil
}

// NewConfigMust builds the logger configuration.
// If there are errors, it panics. Panic is allowed:
// the logger middleware of the application cannot function
// without a running logger.
func NewConfigMust() config {
	config, err := newConfig()
	if err != nil {
		err = fmt.Errorf("get Logger config: %w", err)
		panic(err)
	}

	return config
}
