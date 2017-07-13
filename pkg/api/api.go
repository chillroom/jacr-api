package api

import (
	"time"

	"github.com/qaisjp/jacr-api/pkg/api/auth"
	"github.com/qaisjp/jacr-api/pkg/api/base"
	"github.com/qaisjp/jacr-api/pkg/api/bot"
	"github.com/qaisjp/jacr-api/pkg/api/history"
	"github.com/qaisjp/jacr-api/pkg/api/jwt"
	"github.com/qaisjp/jacr-api/pkg/api/notices"
	"github.com/qaisjp/jacr-api/pkg/api/responses"
	"github.com/qaisjp/jacr-api/pkg/api/slack"
	"github.com/qaisjp/jacr-api/pkg/api/statistics"
	"github.com/qaisjp/jacr-api/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// NewAPI sets up a new API module.
func NewAPI(
	conf *config.Config,
	log *logrus.Logger,
	db *sqlx.DB,
) *base.API {

	router := gin.Default()

	a := &base.API{
		Config: conf,
		Log:    log,
		DB:     db,
		Gin:    router,
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
	router.POST("/v2/slack/invite", slack.Invite)
	router.GET("/v2/slack/badge.svg", slack.CheckOrigin, slack.Badge)

	notices := notices.Impl{API: a}
	router.GET("/v2/notices/", notices.List)
	router.PATCH("/v2/notices/", verifyAuth, notices.Patch)

	responses := responses.Impl{API: a}
	router.GET("/v2/responses/", responses.List)

	history := history.Impl{API: a}
	router.GET("/v2/history/", history.List)

	statistics := statistics.Impl{API: a}
	router.GET("/v2/statistics", statistics.List)

	bot := bot.Impl{API: a}
	router.POST("/v2/bot/restart", verifyAuth, bot.Restart)

	return a
}
