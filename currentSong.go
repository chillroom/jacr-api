package main

import (
	"github.com/gin-gonic/gin"
	r "gopkg.in/dancannon/gorethink.v2"
)

func currentSongEndpoint(c *gin.Context) {
	res, err := r.Table("history").
		OrderBy(r.OrderByOpts{Index: r.Desc("time")}).
		Limit(1).
		Pluck("song", "user", "time").
		Merge(map[string]interface{}{
			"dj":   r.Table("users").Get(r.Row.Field("user")).Field("username"),
			"song": r.Table("songs").Get(r.Row.Field("song")).Field("name"),
		},
		).
		Without("user").
		Run(rethinkSession)

	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "could not get the current song",
			"message": err.Error(),
		})
		return
	}
	defer res.Close()

	var messages interface{}
	err = res.One(&messages)
	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "could not receive current song from response",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"code":   "success",
		"data":   messages,
	})
}
