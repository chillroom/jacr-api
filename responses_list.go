package main

import (
	"github.com/gin-gonic/gin"
	r "gopkg.in/dancannon/gorethink.v2"
)

type Response struct {
	Name      string   `gorethink:"name"`
	Aliases   []string `gorethink:"aliases"`
	Responses []string `gorethink:"responses"`
}

func responsesListEndpoint(c *gin.Context) {

	res, err := r.
		Table("responses").
		OrderBy(r.OrderByOpts{Index: r.Asc("name")}).
		Pluck("name", "responses", "aliases").
		Run(rethinkSession)

	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "could not query responses",
			"message": err.Error(),
		})
		return
	}
	defer res.Close()

	var list []Response
	err = res.All(&list)
	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "could not receive responses from cursor",
			"message": err.Error(),
		})
		return
	}

	// c.JSON(200, gin.H{
	// 	"status": 200,
	// 	"code":   "success",
	// 	"data":   list,
	// })
	// c.JSON(200, list)
	c.HTML(200, "responses.html", list)
}
