package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/qaisjp/jacr-api/pkg/models"
)

func List(c *gin.Context) {
	list := []models.Response{}

	db := c.Keys["db"].(*pg.DB)
	_, err := db.Query(
		&list,
		`SELECT array_agg(cmds.name) as cmds, groups.messages FROM
			response_commands as cmds,
			response_groups as groups
		WHERE
			cmds.group = groups.id
		GROUP BY groups.messages`,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrap(err, "could not receive responses").Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   list,
	})
}
