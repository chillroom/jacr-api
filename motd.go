package main

import (
	"github.com/gin-gonic/gin"
	r "gopkg.in/dancannon/gorethink.v2"
)

func motdListEndpoint(c *gin.Context) {
	res, err := r.Table("settings").Get("motd").Pluck("messages").Field("messages").Run(rethinkSession)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  500,
			"code":    "could not receive MOTD messages",
			"message": err.Error(),
		})
		return
	}
	defer res.Close()

	var messages []interface{}
	err = res.All(&messages)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  500,
			"code":    "could not receive select MOTD message from response",
			"message": err.Error(),
		})
		return
	}

	out := c.JSON
	if c.Query("mode") == "pretty" {
		out = c.IndentedJSON
	}
	out(200, gin.H{
		"status":  200,
		"code":    "success",
		"message": messages,
	})
}

func motdPutEndpoint(c *gin.Context) {
	var messages []string
	if c.BindJSON(&messages) != nil {
		c.JSON(500, gin.H{
			"status":  500,
			"code":    "error",
			"message": "nothing to save",
		})
		return
	}

	_, err := r.Table("settings").Get("motd").Update(map[string]interface{}{
		"messages": messages,
	}).RunWrite(rethinkSession)

	if err != nil {
		c.JSON(500, gin.H{
			"status":  500,
			"code":    "failed to update settings",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  200,
		"code":    "success",
		"message": "message of the day listing updated",
	})
}
