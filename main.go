package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

var conf = struct {
	SlackURL      string
	SlackToken    string
	SlackChannels string
}{}

func main() {
	fs := getFlagSet()
	fs.Parse(os.Args[1:])

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

	address := fs.Lookup("http_address").Value.String()

	router := gin.Default()
	router.POST("/invite", slackHandler)

	http.ListenAndServe(address, router)
}
