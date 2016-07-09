package main

import (
	"github.com/gin-gonic/gin"
	r "gopkg.in/dancannon/gorethink.v2"
)

func opListEndpoint(c *gin.Context) {
	res, err := r.
		Table("songs").
		Filter(
			r.Row.Field("plays").Gt(4),
		).
		Avg("plays").
		Default(100).
		Ceil().
		Run(rethinkSession)

	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "could not calculate the avg number of plays",
			"message": err.Error(),
		})
		return
	}
	defer res.Close()

	// returns a number
	var avg int
	err = res.One(&avg)
	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "could not receive the avg number of plays",
			"message": err.Error(),
		})
		return
	}

	res, err = r.Table("songs").Filter(
		r.Row.Field("plays").
			Gt(8).
			And(r.Row.Field("lastPlay").Gt(r.Now().Add(-1209600))).
			And(r.Row.Field("skipReason").Eq(nil)),
	).
		OrderBy(r.Desc("lastPlay")).
		OrderBy(r.Desc("plays")).
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
