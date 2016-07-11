package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/appleboy/gin-jwt.v2"
	r "gopkg.in/dancannon/gorethink.v2"
	"time"
)

func getAuthMiddleware() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:         "jacr-admin",
		Key:           []byte("secret key"),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour * 24,
		Authenticator: authAuthenticator,
		Authorizator:  authAuthorizator,
		Unauthorized:  authOnFail,
	}
}

type AdminUser struct {
	ID       string `gorethink:"id"`
	Level    int    `gorethink:"level"`
	Username string `gorethink:"username"`
	Password string `gorethink:"password"`
}

func authAuthenticator(user string, password string, c *gin.Context) (username string, success bool) {
	res, err := r.Table("admins").GetAllByIndex("username", user).Pluck("id", "password").Run(rethinkSession)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"status":  "internal error checking database",
			"message": err.Error(),
		})
		c.Set("written", true)
		return
	}
	defer res.Close()

	var doc map[string]string
	err = res.One(&doc)

	if (err == nil) && (bcrypt.CompareHashAndPassword([]byte(doc["password"]), []byte(password)) == nil) {
		return doc["id"], true
	}

	c.JSON(500, gin.H{
		"code":    401,
		"message": "incorrect username/password",
	})
	c.Set("written", true)
	return
}

func authAuthorizator(uid string, c *gin.Context) bool {
	res, err := r.Table("admins").Get(uid).Run(rethinkSession)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"status":  "could not authenticate",
			"message": err.Error(),
		})
		c.Set("written", true)
		return false
	}
	defer res.Close()

	var user AdminUser
	err = res.One(&user)

	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"status":  "could not authenticate",
			"message": "could not find user",
		})
		c.Set("written", true)
		return false
	}

	c.Set("user", user)
	return true
}

func authOnFail(c *gin.Context, code int, message string) {
	written, _ := c.Get("written")
	if written == true {
		return
	}

	c.JSON(500, gin.H{
		"code":    401,
		"message": "not authorised",
	})
}
