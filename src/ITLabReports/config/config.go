package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type MongoDBConfig struct {
	URI     string `envconfig:"ITLAB_REPORTS_MONGO_URI"`
	TestURI string `envconfig:"ITLAB_REPORTS_MONGO_TEST_URI"`
}

type AppConfig struct {
	AppPort string `envconfig:"ITLAB_REPORTS_APP_PORT"`

	TestMode bool `envconfig:"ITLAB_REPORTS_APP_TEST_MODE"`
}

type Config struct {
	MongoDB MongoDBConfig
	App     AppConfig
	Auth    AuthConfig
}

type AuthConfig struct {
	KeyURL   string `envconfig:"ITLAB_REPORTS_AUTH_KEY_URL"`
	Audience string `envconfig:"ITLAB_REPORTS_AUTH_AUDIENCE" default:"itlab"`
	Issuer   string `envconfig:"ITLAB_REPORTS_AUTH_ISSUER"`
	Scope    string `envconfig:"ITLAB_REPORTS_AUTH_SCOPE"`
	Roles    Roles
}

type Roles struct {
	User       string `envconfig:"ITLAB_REPORTS_ROLE_USER" default:"user"`
	Admin      string `envconfig:"ITLAB_REPORTS_ROLE_ADMIN" default:"reports.admin"`
	SuperAdmin string `envconfig:"ITLAB_REPORTS_ROLE_SUPER_ADMIN" default:"admin"`
}

func GetConfig() *Config {
	return GetConfigFrom("./.env")
}

func GetConfigFrom(filePath string) *Config {
	var config Config
	if err := godotenv.Load(filePath); err != nil {
		log.Warn("Don't find .env file")
	}

	if err := envconfig.Process("itlab_reports", &config); err != nil {
		log.WithFields(
			log.Fields{
				"function": "envconfig.Process",
				"error":    err,
			},
		).Fatal("Can't read env vars, shutting down...")
	}
	return &config
}
