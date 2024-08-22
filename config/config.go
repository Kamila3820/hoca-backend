package config

import (
	"fmt"
	"sync"
	"time"
)

type (
	Config struct {
		Server   *Server   `mapstructure:"server" validate:"required"`
		State    *State    `mapstructure:"state" validate:"required"`
		Database *Database `mapstructure:"database" validate:"required"`
	}

	Server struct {
		Port         int           `mapstructure:"port" validate:"required"`
		AllowOrigins []string      `mapstructure:"allowOrigins" validate:"required"`
		BodyLimit    string        `mapstructure:"bodyLimit" validate:"required"`
		TimeOut      time.Duration `mapstructure:"timeout" validate:"required"`
	}

	OAuth2 struct {
	}

	endpoint struct {
	}

	State struct {
		Secret    string        `mapstructure:"secret" validate:"required"`
		ExpiresAt time.Duration `mapstructure:"expiresAt" validate:"required"`
		Issuer    string        `mapstructure:"issuer" validate:"required"`
	}

	Database struct {
		Host     string `mapstructure:"host" validate:"required"`
		Port     int    `mapstructure:"port" validate:"required"`
		User     string `mapstructure:"user" validate:"required"`
		Password string `mapstructure:"password" validate:"required"`
		DBName   string `mapstructure:"dbname" validate:"required"`
		SSLMode  string `mapstructure:"sslmode" validate:"required"`
		Schema   string `mapstructure:"schema" validate:"required"`
	}
)

var (
	once           sync.Once
	configInstance *Config
)

func ConfigGetting() *Config {
	once.Do(func() {

	})
	fmt.Print("hello")
	return nil
}
