package main

import (
	"github.com/gin-gonic/gin"
	r "gopkg.in/dancannon/gorethink.v2"
)

func historyListEndpoint(c *gin.Context) {

	res, err := r.
		Table("history").
		OrderBy(r.OrderByOpts{Index: r.Desc("time")}).
		Without("platformID").
		Merge(map[string]interface{}{
			"user": r.Table("users").Get(r.Row.Field("user")).Pluck("username", "id").Default(nil),
			"song": r.Table("songs").Get(r.Row.Field("song")).Pluck("name", "id").Default(nil),
		}).
		Limit(500).
		Run(rethinkSession)

	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "could not query history",
			"message": err.Error(),
		})
		return
	}
	defer res.Close()

	var list []interface{}
	err = res.All(&list)
	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "could not receive history from cursor",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"code":   "success",
		"data":   list,
	})
}

func historyUserListEndpoint(c *gin.Context) {
	res, err := r.
		Table("history").
		Filter(r.Row.Field("user").Eq(c.Param("user"))).
		OrderBy(r.Desc("time")).
		Without("platformID").
		Merge(map[string]interface{}{
			"song": r.Table("songs").Get(r.Row.Field("song")).Pluck("name", "fkid", "type").Default(nil),
		}).
		Limit(500).
		Run(rethinkSession)

	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "could not query history",
			"message": err.Error(),
		})
		return
	}
	defer res.Close()

	var list []interface{}
	err = res.All(&list)
	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "could not receive history from cursor",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"code":   "success",
		"data":   list,
	})
}
