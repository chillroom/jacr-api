package base

import (
	"github.com/qaisjp/jacr-api/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
)

// API contains all the dependencies of the API server
type API struct {
	Log    *logrus.Logger
	DB     *pg.DB
	Gin    *gin.Engine
	Config *config.Config
}
