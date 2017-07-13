package statistics

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/qaisjp/jacr-api/pkg/config"
	"github.com/sirupsen/logrus"
)

// Statistics contains all the dependencies of the Statistics Generation server
type Statistics struct {
	Config *config.Config
	Log    *logrus.Logger
	DB     *sqlx.DB

	Queue chan *Generator
}

// NewStatistics sets up a new Statistics module
func NewStatistics(
	conf *config.Config,
	log *logrus.Logger,
	db *sqlx.DB,
) *Statistics {

	s := &Statistics{
		Config: conf,
		Log:    log,
		DB:     db,
	}

	s.Queue = make(chan *Generator, 2) // only process two at a time I guess?

	// Initialise each generator by running them
	for _, gen := range GetGenerators() {
		go gen.Spawn(s.Queue)
	}

	return s
}

// Start begins handling all of the statistic generators in the queue.
func (a *Statistics) Start() {
	for gen := range a.Queue {
		err := gen.Run(a.DB)

		if err != nil {
			a.Log.WithFields(logrus.Fields{
				"module":    "statistics",
				"error":     err.Error(),
				"generator": gen.Name,
			}).Warn("Generator failed to run")
		} else {
			a.Log.WithFields(logrus.Fields{
				"module":    "statistics",
				"generator": gen.Name,
				"frequency": gen.Frequency,
			}).Info("Generator succeeded")
		}
	}
}

// Shutdown shuts down the Statistics server
func (a *Statistics) Shutdown(ctx context.Context) error {
	// if err := a.Server.Shutdown(ctx); err != nil {
	// 	return err
	// }

	return nil
}
