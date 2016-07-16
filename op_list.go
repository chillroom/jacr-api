package main

import (
	"github.com/gin-gonic/gin"
	r "gopkg.in/dancannon/gorethink.v2"
)

func opListEndpoint(c *gin.Context) {
	res, err := r.Table("songs").Filter(
		r.Row.Field("recentPlays").
			Gt(10).
			And(r.Row.Field("lastPlay").Gt(r.Now().Add(-5260000))).
			And(r.Row.Field("skipReason").Eq(nil)),
	).
		OrderBy(r.Desc("lastPlay")).
		OrderBy(r.Desc("recentPlays")).
		Merge(map[string]interface{}{
			"plays": r.Row.Field("recentPlays"), // for compatibility with current code
		}).
		Run(rethinkSession)

	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "could not calculate OP list",
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
			"code":    "could not receive OP list from cursor",
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
