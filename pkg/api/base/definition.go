package base

import (
	"github.com/qaisjp/jacr-api/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// API contains all the dependencies of the API server
type API struct {
	Config *config.Config
	Log    *logrus.Logger
	DB     *sqlx.DB
	Gin    *gin.Engine
}
