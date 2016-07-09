package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	r "gopkg.in/dancannon/gorethink.v2"
	"net/http"
	"net/url"
	"os"
)

var conf = struct {
	SlackURL      string
	SlackToken    string
	SlackChannels string
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

	router := gin.Default()
	router.POST("/invite", slackHandler)
	router.GET("/badge-social.svg", slackImageHandler)

	motd := router.Group("/motd")
	{
		motd.GET("/list", motdListEndpoint)
	}

	router.GET("/api/current-song", currentSongEndpoint)
	router.GET("/api/op", opListEndpoint)
	router.GET("/api/history", historyListEndpoint)
	router.GET("/api/history/:user", historyUserListEndpoint)

	address := fs.Lookup("http_address").Value.String()
	http.ListenAndServe(address, router)
}

func checkJACROrigin(c *gin.Context) bool {
	origin := c.Request.Header.Get("Origin")
	parsedOrigin, err := url.Parse(origin)
	if err != nil {
		return false
	}

	if (parsedOrigin.Host == "just-a-chill-room.net") || (parsedOrigin.Host == "www.just-a-chill-room.net") {
		c.Header("Access-Control-Allow-Origin", origin)
		return true
	}
	return true
}
