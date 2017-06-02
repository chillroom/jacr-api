package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

var conf = struct {
	SlackURL      string
	SlackToken    string
	SlackChannels string

	Address string
}{}

var db *pg.DB

func main() {
	var err error

	fs := getFlagSet()
	fs.Parse(os.Args[1:])

	db = pg.Connect(&pg.Options{
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

	loadRoutes()
}

func loadRoutes() {
	router := gin.Default()

	loadTemplates(router)

	router.POST("/invite", slackHandler)
	router.GET("/badge-social.svg", slackImageHandler)

	/////////LEGACY
	motd_legacy := router.Group("/motd")
	{
		motd_legacy.GET("/list", motdListEndpoint)
	}

	router.GET("/api/current-song", currentSongEndpoint)
	router.GET("/api/op", opListEndpoint)
	router.GET("/api/history", historyListEndpoint)
	router.GET("/api/history/:user", historyUserListEndpoint)
	///////////////

	/////
	user_face := router.Group("/user")
	{
		user_face.GET("/responses", responsesListEndpoint)
	}

	// temporary cheats
	router.POST("/_/restart", restartCheatEndpoint)

	// authMiddleware := &auth.GinJWTMiddleware{
	// 	Key:        []byte("secret key"),
	// 	Timeout:    time.Hour,
	// 	MaxRefresh: time.Hour * 24,
	// 	Rethink:    rethinkSession,
	// }

	// authFunc := authMiddleware.MiddlewareFunc

	v1 := router.Group("/v1")
	{
		// auth := v1.Group("/auth")
		// auth.POST("/login", authMiddleware.LoginHandler)
		// auth.POST("/refresh", authMiddleware.RefreshHandler)

		motd := v1.Group("/motd")
		motd.GET("/", motdListEndpoint)
		// motd.PUT("/", authFunc(), motdPutEndpoint)
	}

	http.ListenAndServe(conf.Address, router)
}

func loadTemplates(g *gin.Engine) {
	g.LoadHTMLFiles("templates/responses.html")
}
