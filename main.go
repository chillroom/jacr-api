package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/qaisjp/jacr-api/pkg/api"
	"github.com/qaisjp/jacr-api/pkg/config"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
)

func main() {
	var err error

	fs := getFlagSet()
	fs.Parse(os.Args[1:])

	conf := &config.Config{}

	conf.SlackURL = fs.Lookup("slack_url").Value.String()
	if conf.SlackURL == "" {
		fmt.Println("slack_url is empty")
		return
	}
	conf.SlackToken = fs.Lookup("slack_token").Value.String()
	if conf.SlackToken == "" {
		fmt.Println("slack_token is empty")
		return
	}
	conf.SlackChannels = fs.Lookup("slack_channels").Value.String()
	if conf.SlackChannels == "" {
		fmt.Println("slack_channels is empty")
		return
	}

	conf.Address = fs.Lookup("http_address").Value.String()

	conf.JWTSecret = fs.Lookup("jwt_secret").Value.String()
	if conf.JWTSecret == "" {
		fmt.Println("jwt_secret is empty")
		return
	}

	// var (
	// 	db *sqlx.DB
	// 	gq *goqu.Database
	// )

	db := pg.Connect(&pg.Options{
		Addr:     fs.Lookup("postgres_addr").Value.String(),
		User:     fs.Lookup("postgres_user").Value.String(),
		Database: fs.Lookup("postgres_database").Value.String(),
		Password: fs.Lookup("postgres_password").Value.String(),
	})

	_, err = db.Exec("SELECT 1")
	if err != nil {
		log.Fatal("Postgres connection error!")
		return
	}

	log.Infoln("Connected to PostgreSQL (OLD).")

	db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		query, err := event.FormattedQuery()
		if err != nil {
			panic(err)
		}

		log.Printf("%s %s", time.Since(event.StartTime), query)
	})

	loadRoutes(db, conf)
}

func loadRoutes(db *pg.DB, conf *config.Config) {
	defer db.Close()

	logger := log.StandardLogger()
	logger.Level = log.DebugLevel

	api := api.NewAPI(
		logger,
		db,
		conf,
	)

	http.ListenAndServe(conf.Address, api.Gin)
}
