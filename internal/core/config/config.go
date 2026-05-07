// Package config provides application-wide core configuration settings.
package config

import (
	"fmt"
	"os"
	"time"
)

// Config holds the globally shared application settings
type Config struct {
	TimeZone *time.Location
}

// newConfig reads the required environment variables into the Config struct.
// It returns error if a required variable is malformed or missed.
func newConfig() (*Config, error) {
	tz := os.Getenv("TIME_ZONE")

	zone, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf(
			"load time zone: %s: %w",
			tz,
			err,
		)
	}

	return &Config{
		TimeZone: zone,
	}, nil
}

// NewConfigMust builds the core configuration.
// If there are errors, it panics.
// Panic is allowed: the application cannot function properly
// without the required server's initialization settings.
func NewConfigMust() *Config {
	config, err := newConfig()
	if err != nil {
		err = fmt.Errorf("get core config: %w", err)
		panic(err)
	}

	return config
}
