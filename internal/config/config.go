// Package config provides a singleton instance of the configuration.
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// EnvType is the type of the environment.
type EnvType string

const (
	test EnvType = "test"
	prod EnvType = "prod"
	dev  EnvType = "dev"
)

type (
	// Config is the configuration for the application.
	Config struct {
		Environment EnvType `envconfig:"ENVIRONMENT" required:"true"`
		HTTP        HTTP
		Postgres    Postgres
		Logger      Logger
		Token       Token
		CORS        CORS
	}

	// HTTP is the configuration for the HTTP server.
	HTTP struct {
		Host           string        `envconfig:"HTTP_HOST"             required:"true"`
		Port           string        `envconfig:"HTTP_PORT"             required:"true"`
		MaxHeaderBytes int           `envconfig:"HTTP_MAX_HEADER_BYTES"                 default:"1"`
		ReadTimeout    time.Duration `envconfig:"HTTP_READ_TIMEOUT"                     default:"10s"`
		WriteTimeout   time.Duration `envconfig:"HTTP_WRITE_TIMEOUT"                    default:"10s"`
	}

	// Postgres is the configuration for the Postgres database.
	Postgres struct {
		Host     string `envconfig:"POSTGRES_HOST"     required:"true"`
		Port     string `envconfig:"POSTGRES_PORT"     required:"true"`
		DBName   string `envconfig:"POSTGRES_DBNAME"   required:"true"`
		User     string `envconfig:"POSTGRES_USER"     required:"true"`
		Password string `envconfig:"POSTGRES_PASSWORD" required:"true" json:"-"`
		SSLMode  string `envconfig:"POSTGRES_SSLMODE"                           default:"disable"`
	}

	// Logger is the configuration for the logger.
	Logger struct {
		Level string `envconfig:"LOGGER_LEVEL" default:"info"`
	}

	// Token is the configuration for the token.
	Token struct {
		SecretKey string        `envconfig:"TOKEN_SECRET_KEY" required:"true" json:"-"`
		Expired   time.Duration `envconfig:"TOKEN_EXPIRED"                             default:"15m"`
	}

	// CORS is the configuration for the CORS.
	CORS struct {
		AllowOrigins []string `envconfig:"CORS_ALLOW_ORIGINS" required:"true"`
	}
)

// IsDev check that environment is dev.
func (c *Config) IsDev() bool {
	return c.Environment == dev
}

// IsTest check that environment is test.
func (c *Config) IsTest() bool {
	return c.Environment == test
}

// IsProd check that environment is prod.
func (c *Config) IsProd() bool {
	return c.Environment == prod
}

var (
	instance Config
	once     sync.Once
)

// Get returns the singleton instance of the configuration.
func Get() *Config {
	once.Do(func() {
		if err := envconfig.Process("", &instance); err != nil {
			log.Fatal(err)
		}

		switch instance.Environment {
		case test, prod, dev:
		default:
			log.Fatal("config environment should be test, prod or dev")
		}
		if instance.IsDev() {
			configBytes, err := json.MarshalIndent(instance, "", " ")
			if err != nil {
				log.Fatal(fmt.Errorf("error marshaling indent config: %w", err))
			}
			fmt.Println("Configuration:", string(configBytes))
		}
	})

	return &instance
}
