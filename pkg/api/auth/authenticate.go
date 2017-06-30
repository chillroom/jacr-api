package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/qaisjp/jacr-api/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(username string, password string, c *gin.Context) (userID int, success bool) {
	var u models.User

	db := c.Keys["db"].(*pg.DB)
	_, err := db.QueryOne(&u, "SELECT id, password FROM users WHERE username = ?", username)
	if err != nil {
		if pg.ErrNoRows == err {
			return
		}

		c.JSON(500, gin.H{
			"status":  "error",
			"data":    nil,
			"message": errors.Wrapf(err, "authentication query failed").Error(),
		})

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if (err != nil) && (err != bcrypt.ErrMismatchedHashAndPassword) {
		c.JSON(500, gin.H{
			"status":  "error",
			"data":    nil,
			"message": errors.Wrapf(err, "could not compare hash and password").Error(),
		})

		return
	}

	fmt.Println(err)

	return u.ID, err != bcrypt.ErrMismatchedHashAndPassword
}
