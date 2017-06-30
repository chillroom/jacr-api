package old

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

func MotdListEndpoint(c *gin.Context) {
	db := c.Keys["db"].(*pg.DB)

	var messages []struct {
		ID      int
		Message string
		Title   string
	}
	_, err := db.Query(&messages, `SELECT * FROM notices`)

	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "could not get the notices",
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
		"message": &messages,
	})
}

func motdPutEndpoint(c *gin.Context) {
	// var messages []string
	// if c.BindJSON(&messages) != nil {
	// 	c.JSON(500, gin.H{
	// 		"status":  500,
	// 		"code":    "error",
	// 		"message": "invalid json. expected array of strings.",
	// 	})
	// 	return
	// }

	// _, err := r.Table("settings").Get("motd").Update(map[string]interface{}{
	// 	"messages": messages,
	// }).RunWrite(rethinkSession)

	// if err != nil {
	// 	c.JSON(500, gin.H{
	// 		"status":  500,
	// 		"code":    "failed to update settings",
	// 		"message": err.Error(),
	// 	})
	// 	return
	// }

	// c.JSON(200, gin.H{
	// 	"status":  200,
	// 	"code":    "success",
	// 	"message": "message of the day listing updated",
	// })
}
