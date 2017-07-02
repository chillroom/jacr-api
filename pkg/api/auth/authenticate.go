package auth

import (
	"net/http"

	"github.com/qaisjp/jacr-api/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	goqu "gopkg.in/doug-martin/goqu.v4"
)

func (i *Impl) Authenticate(username string, password string, c *gin.Context) (userID int, success bool) {
	var u models.User

	found, err := i.GQ.From("users").Where(goqu.Ex{"username": username}).ScanStruct(&u)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"data":    nil,
			"message": errors.Wrapf(err, "authentication query failed").Error(),
		})

		return
	}

	if !found {
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
