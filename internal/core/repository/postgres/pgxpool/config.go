package pgxpool

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// config holds the environment-based settings required to connect
// to the PostgreSQL database via pgx and define operation timeouts.
type config struct {
	Host     string        `envconfig:"HOST" required:"true"`
	Port     string        `envconfig:"PORT" default:"5432"`
	User     string        `envconfig:"USER" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Database string        `envconfig:"DB" required:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" required:"true"`
}

// newConfig parses the system environment variables with the "POSTGRES_" prefix
// into the config struct. It returns an error if required variables are missed or malformed.
func newConfig() (config, error) {
	var cfg config

	if err := envconfig.Process("POSTGRES", &cfg); err != nil {
		return config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return cfg, nil
}

// NewConfigMust builds the database configuration.
// If there are errors, it panics. Panic is allowed:
// the application cannot function without a valid database connection.
func NewConfigMust() config {
	cfg, err := newConfig()
	if err != nil {
		err = fmt.Errorf("get Postgres DB config: %w", err)
		panic(err)
	}

	return cfg
}
