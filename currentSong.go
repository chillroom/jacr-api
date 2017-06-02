package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

func currentSongEndpoint(c *gin.Context) {
	var result struct {
		Fkid     string
		Name     string
		Type     string
		Time     time.Time
		Username string
	}

	_, err := db.QueryOne(&result, `SELECT songs.fkid, songs.name, songs.type, history.time, dubtrack_users.username
		FROM history, songs, dubtrack_users
		WHERE
		(history.song = songs.id) and
		(history."user" = dubtrack_users.id)
		ORDER BY history.time DESC LIMIT 1`)
	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "could not get the current song",
			"message": err.Error(),
		})
		return
	}

	var output struct {
		Song struct {
			Fkid string `json:"fkid"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"song"`

		DJ   string    `json:"dj"`
		Time time.Time `json:"time"`
	}

	output.Song.Fkid = result.Fkid
	output.Song.Name = result.Name
	output.Song.Type = result.Type
	output.DJ = result.Username
	output.Time = result.Time

	c.JSON(200, gin.H{
		"status": 200,
		"code":   "success",
		"data":   output,
	})
}
