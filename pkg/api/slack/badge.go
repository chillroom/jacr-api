package slack

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (i *Impl) Badge(c *gin.Context) {
	v := url.Values{}
	v.Set("token", i.Config.SlackToken)
	v.Set("presence", "1")
	resp, err := http.PostForm("https://"+i.Config.SlackURL+"/api/users.list", v)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrap(err, "could not connect to slack").Error(),
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Could not decode Slack response",
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
	c.Redirect(http.StatusFound, url)
}
