package server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// config holds the environment-based settings required to configure
// server, specifically the server path and timeout for the graceful shutdown.
type config struct {
	Addr            string        `envconfig:"ADDR" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
}

// newConfig parses the system environment variables with the "HTTP_" prefix
// into the config struct. It returns an error if required variables are missed or malformed.
func newConfig() (config, error) {
	var cfg config

	if err := envconfig.Process("HTTP", &cfg); err != nil {
		return config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return cfg, nil
}

// NewConfigMust builds the server configuration.
// If there are errors, it panics.
// Panic is allowed: server does not function properly without appropriate settings.
func NewConfigMust() config {
	cfg, err := newConfig()
	if err != nil {
		err = fmt.Errorf("get HTTP server config: %w", err)
		panic(err)
	}

	return cfg
}
