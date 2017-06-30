package auth

import (
	"net/http"

	"github.com/qaisjp/jacr-api/pkg/models"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/gosimple/slug"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	u := &models.User{
		Username: username,
		Password: password,
		Email:    email,
		Slug:     slug.Make(username),
	}

	success, err := govalidator.ValidateStruct(u)
	if !success {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	count := 0
	db := c.Keys["db"].(*pg.DB)
	_, err = db.Query(&count, "SELECT COUNT(id) from users WHERE (username = ?) or (slug = ?) or (email = ?)", u.Username, u.Slug, u.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
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

	err = db.Insert(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
	})
}
