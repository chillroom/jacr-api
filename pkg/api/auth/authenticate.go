package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/qaisjp/jacr-api/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

func (i *Impl) Authenticate(username string, password string, c *gin.Context) (userID int, success bool) {
	var u models.User

	_, err := i.DB.QueryOne(&u, "SELECT id, password FROM users WHERE username = ?", username)
	if err != nil {
		if pg.ErrNoRows == err {
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"data":    nil,
			"message": errors.Wrapf(err, "authentication query failed").Error(),
		})

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if (err != nil) && (err != bcrypt.ErrMismatchedHashAndPassword) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"data":    nil,
			"message": errors.Wrapf(err, "could not compare hash and password").Error(),
		})

		return
	}

	return u.ID, err != bcrypt.ErrMismatchedHashAndPassword
}
