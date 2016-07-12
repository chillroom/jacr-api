package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qaisjp/jacr-api/auth"
	r "gopkg.in/dancannon/gorethink.v2"
	"net/http"
	"os"
	"time"
)

var conf = struct {
	SlackURL      string
	SlackToken    string
	SlackChannels string

	Address string
}{}

var rethinkSession *r.Session

func main() {
	var err error

	fs := getFlagSet()
	fs.Parse(os.Args[1:])

	rethinkSession, err = r.Connect(r.ConnectOpts{
		Address:  fs.Lookup("rethinkdb_address").Value.String(),
		Database: fs.Lookup("rethinkdb_database").Value.String(),
		Username: fs.Lookup("rethinkdb_username").Value.String(),
		Password: fs.Lookup("rethinkdb_password").Value.String(),
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

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

	authMiddleware := &auth.GinJWTMiddleware{
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour * 24,
		Rethink:    rethinkSession,
	}

	authFunc := authMiddleware.MiddlewareFunc

	v1 := router.Group("/v1")
	{
		auth := v1.Group("/auth")
		auth.POST("/login", authMiddleware.LoginHandler)
		auth.POST("/refresh", authMiddleware.RefreshHandler)

		motd := v1.Group("/motd")
		motd.GET("/", motdListEndpoint)
		motd.PUT("/", authFunc(), motdPutEndpoint)
	}

	http.ListenAndServe(conf.Address, router)
}
