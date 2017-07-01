package api

import (
	"time"

	"github.com/qaisjp/jacr-api/pkg/api/auth"
	"github.com/qaisjp/jacr-api/pkg/api/base"
	"github.com/qaisjp/jacr-api/pkg/api/jwt"
	"github.com/qaisjp/jacr-api/pkg/api/notices"
	"github.com/qaisjp/jacr-api/pkg/api/old"
	"github.com/qaisjp/jacr-api/pkg/api/responses"
	"github.com/qaisjp/jacr-api/pkg/api/slack"
	"github.com/qaisjp/jacr-api/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
)

// NewAPI sets up a new API module.
func NewAPI(
	log *logrus.Logger,
	db *pg.DB,
	router *gin.Engine,
	conf *config.Config,
) *base.API {

	a := &base.API{
		Log:    log,
		DB:     db,
		Gin:    router,
		Config: conf,
	}

	auth := auth.Impl{API: a}

	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "jacr-api",
		Key:        []byte(conf.JWTSecret),
		Timeout:    time.Hour * 24,
		MaxRefresh: time.Hour * 24,

		Authenticator: auth.Authenticate,
		Authorizator:  auth.Authorize,
		Unauthorized:  auth.Unauthorized,

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}

	router.POST("/v2/auth/login", authMiddleware.LoginHandler)
	router.POST("/v2/auth/register", auth.Register)

	verifyAuth := authMiddleware.MiddlewareFunc()

	slack := slack.Impl{API: a}
	router.POST("/invite", slack.SlackHandler)
	router.GET("/badge-social.svg", slack.CheckOrigin, slack.SlackImageHandler)

	notices := notices.Impl{API: a}
	router.GET("/v2/notices/", notices.List)
	router.PATCH("/v2/notices/", verifyAuth, notices.Patch)

	responses := responses.Impl{API: a}
	router.GET("/v2/responses/", responses.List)

	{
		router.Use(func(c *gin.Context) {
			c.Set("db", db)
			c.Next()
		})
		router.GET("/motd/list", old.MotdListEndpoint)

		router.GET("/api/current-song", old.CurrentSongEndpoint)
		router.GET("/api/op", old.OpListEndpoint)
		router.GET("/api/history", old.HistoryListEndpoint)
		router.GET("/api/history/:user", old.HistoryUserListEndpoint)

		router.GET("/user/responses", old.ResponsesListEndpoint)

		router.POST("/_/restart", old.RestartCheatEndpoint)
	}

	return a
}
