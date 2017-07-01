package old

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/qaisjp/jacr-api/pkg/models"
)

func ResponsesListEndpoint(c *gin.Context) {
	list := make([]models.Response, 0)

	db := c.Keys["db"].(*pg.DB)
	_, err := db.Query(&list, `
		SELECT array_agg(cmds.name) as cmds, groups.messages FROM
			response_commands as cmds,
			response_groups as groups
		WHERE
			cmds.group = groups.id
		GROUP BY groups.messages`)

	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "could not get the responses",
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
