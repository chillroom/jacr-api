package slack

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

var generalSlackResponses = map[string]string{
	"already_invited": "You have already been invited",
	"already_in_team": "You are already part of our slack group",
	"invalid_email":   "Invalid email address entered",
}

func (i *Impl) CheckJACROrigin(c *gin.Context) bool {
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

func (i *Impl) SlackHandler(c *gin.Context) {
	if !i.CheckJACROrigin(c) {
		return
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
	v.Set("token", i.Config.SlackToken)
	v.Set("set_active", "true")
	v.Set("_attempts", "1")
	resp, err := http.PostForm("https://"+i.Config.SlackURL+"/api/users.admin.invite", v)

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

func (i *Impl) SlackImageHandler(c *gin.Context) {
	if !i.CheckJACROrigin(c) {
		return
	}

	v := url.Values{}
	v.Set("token", i.Config.SlackToken)
	v.Set("presence", "1")
	resp, err := http.PostForm("https://"+i.Config.SlackURL+"/api/users.list", v)

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
		Success bool `json:"ok"`
		Members []struct {
			Bot      bool   `json:"is_bot"`
			Deleted  bool   `json:"deleted"`
			Presence string `json:"presence"`
		} `json:"members"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if (err != nil) || !data.Success {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "something_wrong",
			"message": "Something went wrong",
		})
		return
	}

	total := 0
	online := 0
	for _, member := range data.Members {
		if !member.Bot && !member.Deleted {
			total++
			if member.Presence == "active" {
				online++
			}
		}
	}

	// build the url
	url := "https://img.shields.io/badge/slack-" + strconv.Itoa(online) + "%2F" + strconv.Itoa(total) + "%20active-46ccbb.svg?style=flat"
	c.Redirect(302, url)
}
