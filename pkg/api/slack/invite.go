package slack

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var generalSlackResponses = map[string]string{
	"already_invited": "You have already been invited",
	"already_in_team": "You are already part of our slack group",
	"invalid_email":   "Invalid email address entered",
}

func (i *Impl) Invite(c *gin.Context) {
	email, _ := c.GetPostForm("email")

	// No email? Tell them that.
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrap(err, "could not query slack").Error(),
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Couldn't decode Slack's response. Try to log in anyway.",
		})
		return
	}

	if data.Success {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Invite sent!",
		})
		return
	} else if err := generalSlackResponses[data.Error]; err != "" {
		c.JSON(resp.StatusCode, gin.H{
			"status":  "error",
			"data":    data.Error,
			"message": err,
		})
		return
	}

	c.JSON(resp.StatusCode, gin.H{
		"status":  "error",
		"data":    data.Error,
		"message": data.Error,
	})
}
