package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/qaisjp/jacr-api/pkg/models"
	goqu "gopkg.in/doug-martin/goqu.v4"
)

func (i *Impl) Authorize(userId int, c *gin.Context) bool {
	var u models.User

	found, err := i.GQ.From("users").Where(goqu.Ex{"id": userId}).ScanStruct(&u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrapf(err, "could not authorize user").Error(),
		})
	}
	if !found {
		return false
	}

	c.Set("user", u)

	return (u.Level > 1) && (!u.Banned) && (u.Activated)
}
