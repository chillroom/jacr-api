package auth

import (
	"net/http"

	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/qaisjp/jacr-api/pkg/models"
)

func (i *Impl) Authorize(userId int, c *gin.Context) bool {
	var u models.User

	err := i.DB.Get(&u, "SELECT * FROM accounts WHERE id = $1", userId)

	if err == sql.ErrNoRows {
		return false
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrapf(err, "could not authorize user").Error(),
		})
	}

	c.Set("user", u)

	return (u.Level > 1) && (!u.Banned) && (u.Activated)
}
