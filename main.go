package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/qaisjp/jacr-api/pkg/api/auth"
	"github.com/qaisjp/jacr-api/pkg/api/jwt"
	"github.com/qaisjp/jacr-api/pkg/api/old"
	"github.com/qaisjp/jacr-api/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"golang.org/x/crypto/bcrypt"
)

var conf = struct {
	SlackURL      string
	SlackToken    string
	SlackChannels string

	JWTSecret string
	Address   string
}{}

func main() {
	var err error

	fs := getFlagSet()
	fs.Parse(os.Args[1:])

	conf.SlackURL = fs.Lookup("slack_url").Value.String()
	if conf.SlackURL == "" {
		fmt.Println("slack_url is empty")
		return
	}
	conf.SlackToken = fs.Lookup("slack_token").Value.String()
	if conf.SlackToken == "" {
		fmt.Println("slack_token is empty")
		return
	}
	conf.SlackChannels = fs.Lookup("slack_channels").Value.String()
	if conf.SlackChannels == "" {
		fmt.Println("slack_channels is empty")
		return
	}

	conf.Address = fs.Lookup("http_address").Value.String()

	conf.JWTSecret = fs.Lookup("jwt_secret").Value.String()
	if conf.JWTSecret == "" {
		fmt.Println("jwt_secret is empty")
		return
	}

	db := pg.Connect(&pg.Options{
		Addr:     fs.Lookup("postgres_addr").Value.String(),
		User:     fs.Lookup("postgres_user").Value.String(),
		Database: fs.Lookup("postgres_database").Value.String(),
		Password: fs.Lookup("postgres_password").Value.String(),
	})

	_, err = db.Exec("SELECT 1")
	if err != nil {
		log.Print("Postgres connection error!\n")
		panic(err)
	}
	log.Print("Connected to PostgreSQL.\n")

	loadRoutes(db)
}

func oldRoutes(router *gin.Engine) {
	/////////LEGACY
	motd_legacy := router.Group("/motd")
	{
		motd_legacy.GET("/list", old.MotdListEndpoint)
	}

	router.GET("/api/current-song", old.CurrentSongEndpoint)
	router.GET("/api/op", old.OpListEndpoint)
	router.GET("/api/history", old.HistoryListEndpoint)
	router.GET("/api/history/:user", old.HistoryUserListEndpoint)
	///////////////

	/////
	user_face := router.Group("/user")
	{
		user_face.GET("/responses", old.ResponsesListEndpoint)
	}

	// temporary cheats
	router.POST("/_/restart", old.RestartCheatEndpoint)
}

func getJWTMiddleware(db *pg.DB) *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:      "jacr-api",
		Key:        []byte(conf.JWTSecret),
		Timeout:    time.Hour * 24,
		MaxRefresh: time.Hour * 24,
		Authenticator: func(username string, password string, c *gin.Context) (userID int, success bool) {
			var u models.User
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
		},

		Authorizator: func(userId int, c *gin.Context) bool {
			if userId == 5 {
				return true
			}

			return false
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			if c.Writer.Written() {
				return
			}

			c.JSON(code, gin.H{
				"status":  "error",
				"message": message,
			})
		},

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}
}

func getDatabaseMiddleware(db *pg.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func loadRoutes(db *pg.DB) {
	router := gin.Default()
	router.Use(getDatabaseMiddleware(db))

	loadTemplates(router)

	router.POST("/invite", slackHandler)
	router.GET("/badge-social.svg", slackImageHandler)

	oldRoutes(router)

	authMiddleware := getJWTMiddleware(db)

	v2 := router.Group("/v2")

	authGroup := v2.Group("/auth")
	{
		authGroup.POST("/login", authMiddleware.LoginHandler)
		authGroup.POST("/register", auth.Register)
	}

	rootGroup := v2.Group("/")
	rootGroup.Use(authMiddleware.MiddlewareFunc())
	{
		notices := rootGroup.Group("/notices")
		{
			notices.GET("/", old.MotdListEndpoint)
		}
	}

	http.ListenAndServe(conf.Address, router)
}

func loadTemplates(g *gin.Engine) {
	g.LoadHTMLFiles("templates/responses.html")
}
