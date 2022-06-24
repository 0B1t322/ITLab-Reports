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

	GrpcAppPort string `default:"8081" envconfig:"ITLAB_REPORTS_APP_GRPC_PORT"`

	ITLabURL string `default:"https://dev.manage.rtuitlab.dev" envconfig:"ITLAB_REPORTS_APP_ITLAB_URL"`
}

type Config struct {
	MongoDB MongoDBConfig
	App     AppConfig
	Auth    AuthConfig
}

type AuthConfig struct {
	KeyURL   string `envconfig:"ITLAB_REPORTS_AUTH_KEY_URL"`
	Audience string `default:"itlab" envconfig:"ITLAB_REPORTS_AUTH_AUDIENCE"`
	Issuer   string `envconfig:"ITLAB_REPORTS_AUTH_ISSUER"`
	Scope    string `envconfig:"ITLAB_REPORTS_AUTH_SCOPE"`
	Roles    Roles
}

type Roles struct {
	User       string `default:"user" envconfig:"ITLAB_REPORTS_ROLE_USER"`
	Admin      string `default:"reports.admin" envconfig:"ITLAB_REPORTS_ROLE_ADMIN"`
	SuperAdmin string `default:"admin" envconfig:"ITLAB_REPORTS_ROLE_SUPER_ADMIN"`
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

	log.WithFields(
		log.Fields{
			"from": "GetConfig",
		},
	).Infof("config: %+v", config)
	return &config
}
