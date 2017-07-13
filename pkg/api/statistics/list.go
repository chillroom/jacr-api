package statistics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/qaisjp/jacr-api/pkg/models"
)

func (i *Impl) List(c *gin.Context) {
	stats := make([]models.Statistic, 0)

	err := i.DB.Select(&stats, "select * from statistics")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrap(err, "could not receive stats").Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   stats,
	})
}
