package auth

import (
	"net/http"

	goqu "gopkg.in/doug-martin/goqu.v4"

	"github.com/pkg/errors"
	"github.com/qaisjp/jacr-api/pkg/models"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"golang.org/x/crypto/bcrypt"
)

func (i *Impl) Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	u := models.User{
		Username: username,
		Password: password,
		Email:    email,
		Slug:     slug.Make(username),
	}

	success, err := govalidator.ValidateStruct(&u)
	if !success {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	count, err := i.GQ.From("users").Where(goqu.Or(
		goqu.I("username").Eq(u.Username),
		goqu.I("slug").Eq(u.Slug),
		goqu.I("email").Eq(u.Email),
	)).Count()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrapf(err, "could not check existence").Error(),
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "error",
			"message": "Account already exists with that username, slug, or email",
		})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	u.Password = string(hashedPassword)

	s := i.GQ.From("users").Insert(&u).Sql
	i.Log.Println(s)
	_, err = i.GQ.Exec(s)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrap(err, "could not insert").Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
	})
}
