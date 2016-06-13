package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

var generalSlackResponses = map[string]string{
	"already_invited": "You have already been invited",
	"already_in_team": "You are already part of our slack group",
	"invalid_email":   "Invalid email address entered",
}

func slackHandler(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	parsedOrigin, err := url.Parse(origin)
	if err != nil {
		return
	}

	if (parsedOrigin.Host == "just-a-chill-room.net") || (parsedOrigin.Host == "www.just-a-chill-room.net") {
		c.Header("Access-Control-Allow-Origin", origin)
	}

	email, _ := c.GetPostForm("email")

	// No email? Tell them that.
	if email == "" {
		c.JSON(200, gin.H{
			"status":  400,
			"code":    "email_required",
			"message": "Email address required",
		})
		return
	}

	v := url.Values{}
	v.Set("email", email)
	v.Set("token", conf.SlackToken)
	v.Set("set_active", "true")
	v.Set("_attempts", "1")
	resp, err := http.PostForm("https://"+conf.SlackURL+"/api/users.admin.invite", v)

	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "something_wrong",
			"message": err.Error(),
		})
		return
	}

	defer resp.Body.Close()

	data := struct {
		Success bool   `json:"ok"`
		Error   string `json:"error"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "something_wrong",
			"message": "Couldn't decode Slack's response. Try to log in anyway.",
		})
		return
	}

	if data.Success {
		c.JSON(200, gin.H{
			"status":  200,
			"code":    "invite_sent",
			"message": "Invite sent!",
		})
		return
	} else if err := generalSlackResponses[data.Error]; err != "" {
		c.JSON(200, gin.H{
			"status":  200,
			"code":    data.Error,
			"message": err,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  200,
		"code":    "slack_error",
		"message": data.Error,
	})
}
