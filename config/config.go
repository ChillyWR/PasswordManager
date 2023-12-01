package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	API APIConfig
	DB  DBConfig
}

type APIConfig struct {
	Port            uint          `envConfig:"PM_SERVER_PORT"             default:"5000"`
	ShutdownTimeout time.Duration `envConfig:"PM_SERVER_SHUTDOWN_TIMEOUT" default:"15s"`
}

type DBConfig struct {
	Host     string `envConfig:"PM_DB_HOST"     default:"localhost"`
	Port     string `envConfig:"PM_DB_PORT"     default:"5432"`
	DBName   string `envConfig:"PM_DB_NAME"     default:"password_manager"`
	Username string `envConfig:"PM_DB_USERNAME" default:"admin"`
	Password string `envConfig:"PM_DB_PASSWORD" default:"12345"`
	SSLMode  string `envConfig:"PM_DB_SSL_MODE" default:"disable"`
}

func New() (*Config, error) {
	var c Config
	err := envconfig.Process("PM_SERVER", &c.API)
	if err != nil {
		return nil, fmt.Errorf("process: %w", err)
	}

	err = envconfig.Process("PM_DB", &c.DB)
	if err != nil {
		return nil, fmt.Errorf("process: %w", err)
	}

	return &c, nil
}
