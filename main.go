package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/qaisjp/jacr-api/pkg/api"
	"github.com/qaisjp/jacr-api/pkg/api/old"
	"github.com/qaisjp/jacr-api/pkg/config"

	"github.com/gin-gonic/gin"
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

	db := pg.Connect(&pg.Options{
		Addr:     fs.Lookup("postgres_addr").Value.String(),
		User:     fs.Lookup("postgres_user").Value.String(),
		Database: fs.Lookup("postgres_database").Value.String(),
		Password: fs.Lookup("postgres_password").Value.String(),
	})

	_, err = db.Exec("SELECT 1")
	if err != nil {
		log.Print("Postgres connection error!\n")
		panic(err)
	}
	log.Print("Connected to PostgreSQL.\n")

	db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		query, err := event.FormattedQuery()
		if err != nil {
			panic(err)
		}

		log.Printf("%s %s", time.Since(event.StartTime), query)
	})

	loadRoutes(db, conf)
}

func oldRoutes(router *gin.Engine) {
	/////////LEGACY
	legacy := router.Group("/motd")
	{
		legacy.GET("/list", old.MotdListEndpoint)
	}

	router.GET("/api/current-song", old.CurrentSongEndpoint)
	router.GET("/api/op", old.OpListEndpoint)
	router.GET("/api/history", old.HistoryListEndpoint)
	router.GET("/api/history/:user", old.HistoryUserListEndpoint)
	///////////////

	/////
	user := router.Group("/user")
	{
		user.GET("/responses", old.ResponsesListEndpoint)
	}

	// temporary cheats
	router.POST("/_/restart", old.RestartCheatEndpoint)
}

func loadTemplates(g *gin.Engine) {
	g.LoadHTMLFiles("templates/responses.html")
}

func loadRoutes(db *pg.DB, conf *config.Config) {
	router := gin.Default()

	// just for the old routes
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	loadTemplates(router)

	oldRoutes(router)

	logger := log.StandardLogger()
	logger.Level = log.DebugLevel

	api.NewAPI(
		logger,
		db,
		router,
		conf,
	)

	http.ListenAndServe(conf.Address, router)
	db.Close()
}
