package config

import (
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server   *Server   `mapstructure:"server" validate:"required"`
		OAuth2   *OAuth2   `mapstructure:"oauth2" validate:"required"`
		Google   *Google   `mapstructure:"google" validate:"required"`
		Minio    *Minio    `mapstructure:"minio" validate:"required"`
		Scb      *Scb      `mapstructure:"scb" validate:"required"`
		Database *Database `mapstructure:"database" validate:"required"`
	}

	Server struct {
		Port         int           `mapstructure:"port" validate:"required"`
		AllowOrigins []string      `mapstructure:"allowOrigins" validate:"required"`
		BodyLimit    string        `mapstructure:"bodyLimit" validate:"required"`
		TimeOut      time.Duration `mapstructure:"timeout" validate:"required"`
	}

	OAuth2 struct {
		UserRedirectUrl string   `mapstructure:"userRedirectUrl" validate:"required"`
		ClientID        string   `mapstructure:"clientID" validate:"required"`
		ClientSecret    string   `mapstructure:"clientSecret" validate:"required"`
		Endpoints       endpoint `mapstructure:"endpoints" validate:"required"`
		Scopes          []string `mapstructure:"scopes" validate:"required"`
		UserInfoUrl     string   `mapstructure:"userInfoUrl" validate:"required"`
		RevokeUrl       string   `mapstructure:"revokeUrl" validate:"required"`
	}

	Google struct {
		ApiKey string `mapstructure:"apiKey" validate:"required"`
	}

	Minio struct {
		BucketEndpoint string `mapstructure:"bucketEndpoint" validate:"required"`
		BucketName     string `mapstructure:"bucketName" validate:"required"`
	}

	Scb struct {
		ScbUrl       string `mapstructure:"scbUrl" validate:"required"`
		ScbAppKey    string `mapstructure:"scbAppKey" validate:"required"`
		ScbAppSecret string `mapstructure:"scbAppSecret" validate:"required"`
		ScbBillerId  string `mapstructure:"scbBillerId" validate:"required"`
	}

	endpoint struct {
		AuthUrl       string `mapstructure:"authUrl" validate:"required"`
		TokenUrl      string `mapstructure:"tokenUrl" validate:"required"`
		DeviceAuthUrl string `mapstructure:"deviceAuthUrl" validate:"required"`
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
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		if err := viper.Unmarshal(&configInstance); err != nil {
			panic(err)
		}

		validating := validator.New()

		if err := validating.Struct(configInstance); err != nil {
			panic(err)
		}

	})

	return configInstance
}
