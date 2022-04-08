package server

import (
	"ITLabReports/config"
	_ "ITLabReports/migrations"
	"ITLabReports/utils"
	"context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
	swag "github.com/swaggo/http-swagger"
	_ "ITLabReports/docs"
)

type App struct {
	Router *mux.Router
	DB *mongo.Client
}

var collection *mongo.Collection
var cfg *config.Config

func (a *App) Init(config *config.Config) {
	cfg = config
	log.WithField("dburi", cfg.DB.URI).Info("Current database URI: ")
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.DB.URI))
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.NewClient",
			"error"	:	err,
			"db_uri":	cfg.DB.URI,
		},
		).Fatal("Failed to create new MongoDB client")
	}

	// Create db connect
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Connect",
			"error"	:	err},
		).Fatal("Failed to connect to MongoDB")
	}

	// Check the connection
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(ctx, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Ping",
			"error"	:	err},
		).Fatal("Failed to ping MongoDB")
	}
	log.Info("Connected to MongoDB!")

	dbName := utils.GetDbName(cfg.DB.URI)
	dbCollectionName := "reports"
	db := client.Database(dbName)
	migrate.SetDatabase(db)
	if err := migrate.Up(migrate.AllAvailable); err != nil {
		log.WithFields(log.Fields{
			"function" : "migrate.Up",
			"error"	:	err},
		).Fatal("Failed to migrate MongoDB!")
	}
	ver, desc, err := migrate.Version()
	log.WithFields(log.Fields{
		"db_name" : dbName,
		"collection_name" : dbCollectionName,
		"version" : ver,
		"description" : desc,
	}).Info("Database information: ")

	log.WithField("testMode", cfg.App.TestMode).Info("Let's check if test mode is on...")
	collection = client.Database(dbName).Collection(dbCollectionName)
	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) setRouters() {
	private := a.Router.PathPrefix("").Subrouter()
	docs := a.Router.PathPrefix("/api/reports/swagger")
	docs.Handler(
		swag.WrapHandler,
	)

	private.HandleFunc("/api/reports", getAllReportsSorted).Methods("GET").Queries("sorted_by","{var}")
	private.HandleFunc("/api/reports/employee/{employee}", getEmployeeReports).Methods("GET").Queries("dateBegin","{dateBegin}", "dateEnd", "{dateEnd}")
	private.HandleFunc("/api/reports/employee/{employee}", getEmployeeReports).Methods("GET")
	private.HandleFunc("/api/reports", getAllReports).Methods("GET")
	private.HandleFunc("/api/reports/archived", getArchivedReports).Methods("GET")
	private.HandleFunc("/api/reports/{id}", getReport).Methods("GET")
	private.HandleFunc("/api/reports", createReport).Methods("POST").Queries("implementer","{implementer}")
	private.HandleFunc("/api/reports", createReport).Methods("POST")
	private.HandleFunc("/api/reports/{id}", updateReport).Methods("PUT")
	private.HandleFunc("/api/reports/{id}", deleteReport).Methods("DELETE")

	if cfg.App.TestMode {
		private.Use(testAuthMiddleware)
	} else {
		private.Use(authMiddleware)
	}
}

func (a *App) Run(addr string) {
	log.WithField("port", cfg.App.AppPort).Info("Starting server on port:")
	log.Info("\n\nNow handling routes!")

	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.ListenAndServe",
			"error"	:	err},
		).Fatal("Failed to run a server!")
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}