package notices

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/qaisjp/jacr-api/pkg/models"
)

func List(c *gin.Context) {
	db := c.Keys["db"].(*pg.DB)

	var notices []models.Notice
	_, err := db.Query(&notices, `SELECT * FROM notices`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrap(err, "could not receive notices").Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   &notices,
	})
}
