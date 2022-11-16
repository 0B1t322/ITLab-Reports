package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

var GlobalConfig Config = Config{
	MongoDB: MongoDBConfig{
		URI:     "mongodb://root:root@localhost:27018/itlab-reports?authSource=admin",
		TestURI: "mongodb://root:root@localhost:27018/itlab-reports-test?authSource=admin",
	},
	App: AppConfig{
		AppPort:        "8080",
		TestMode:       true,
		GrpcAppPort:    "",
		ITLabURL:       "",
		SalaryGRPCAddr: "",
	},
	Auth: AuthConfig{
		KeyURL:   "",
		Audience: "",
		Issuer:   "",
		Scope:    "",
		Roles:    Roles{},
	},
	RemoteApi: RemoteApiConfig{
		ClientSecret: "",
		TokenURL:     "",
		ClientID:     "",
	},
}

type MongoDBConfig struct {
	URI     string `envconfig:"ITLAB_REPORTS_MONGO_URI"`
	TestURI string `envconfig:"ITLAB_REPORTS_MONGO_TEST_URI"`
}

type AppConfig struct {
	AppPort string `default:"8080" envconfig:"ITLAB_REPORTS_APP_PORT"`

	TestMode bool `envconfig:"ITLAB_REPORTS_APP_TEST_MODE"`

	GrpcAppPort string `default:"8081" envconfig:"ITLAB_REPORTS_APP_GRPC_PORT"`

	ITLabURL string `default:"https://dev.manage.rtuitlab.dev" envconfig:"ITLAB_REPORTS_APP_ITLAB_URL"`

	SalaryGRPCAddr string `default:"salary:5503" envconfig:"ITLAB_REPORTS_APP_SALARY_GRPC_ADDR"`
}

type Config struct {
	MongoDB   MongoDBConfig
	App       AppConfig
	Auth      AuthConfig
	RemoteApi RemoteApiConfig
}

type RemoteApiConfig struct {
	ClientSecret string `envconfig:"ITLAB_REPORTS_REMOTE_API_CLIENT_SECRET"`
	TokenURL     string `envconfig:"ITLAB_REPORTS_REMOTE_API_TOKEN_URL"`
	ClientID     string `envconfig:"ITLAB_REPORTS_REMOTE_API_CLIENT_ID"`
}

type AuthConfig struct {
	KeyURL   string `envconfig:"ITLAB_REPORTS_AUTH_KEY_URL"`
	Audience string `envconfig:"ITLAB_REPORTS_AUTH_AUDIENCE" default:"itlab"`
	Issuer   string `envconfig:"ITLAB_REPORTS_AUTH_ISSUER"`
	Scope    string `envconfig:"ITLAB_REPORTS_AUTH_SCOPE"`
	Roles    Roles
}

type Roles struct {
	User       string `default:"user"          envconfig:"ITLAB_REPORTS_ROLE_USER"`
	Admin      string `default:"reports.admin" envconfig:"ITLAB_REPORTS_ROLE_ADMIN"`
	SuperAdmin string `default:"admin"         envconfig:"ITLAB_REPORTS_ROLE_SUPER_ADMIN"`
}

func GetConfig() Config {
	return GetConfigFrom("./.env")
}

func GetConfigFrom(filePath string) Config {
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
	return config
}

func InitGlobalConfig() {
	if err := godotenv.Load("./.env"); err != nil {
		log.Warn("Don't find .env file")
	}

	if err := envconfig.Process("itlab_reports", &GlobalConfig); err != nil {
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
	).Infof("config: %+v", GlobalConfig)
}
