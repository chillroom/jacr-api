package main

import (
	"net/http"
	"os"

	"github.com/qaisjp/jacr-api/pkg/api"
	"github.com/qaisjp/jacr-api/pkg/config"
	"github.com/qaisjp/jacr-api/pkg/database"
	"github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
	"github.com/koding/multiconfig"
)

func main() {
	var err error

	m := multiconfig.NewWithPath(os.Getenv("config"))
	cfg := &config.Config{}
	m.MustLoad(cfg)

	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	logger := logrus.StandardLogger()
	logger.Level = logLevel

	logger.WithFields(logrus.Fields{
		"module": "init",
	}).Info("Starting up the application")

	// Initialize the database
	var db *sqlx.DB

	db, err = database.NewPostgres(cfg.Postgres)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"module": "init",
			"error":  err.Error(),
			"cstr":   cfg.Postgres.ConnectionString,
		}).Fatal("Unable to connect to the Postgres server")
		return
	}

	logger.WithFields(logrus.Fields{
		"module": "init",
		"cstr":   cfg.Postgres.ConnectionString,
	}).Info("Connected to a Postgres server")

	api := api.NewAPI(
		cfg,
		logger,
		db,
	)

	{
		router := api.Gin
		router.LoadHTMLFiles("templates/responses.html")

	}

	http.ListenAndServe(cfg.Address, api.Gin)
}
