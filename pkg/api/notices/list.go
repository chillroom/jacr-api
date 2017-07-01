package notices

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/qaisjp/jacr-api/pkg/models"
)

// List allows you to list all the notices in the system
func (i *Impl) List(c *gin.Context) {

	var notices []models.Notice
	_, err := i.DB.Query(&notices, `SELECT * FROM notices`)

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
