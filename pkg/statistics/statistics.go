package statistics

import (
	"context"

	"time"

	"github.com/jmoiron/sqlx"
	"github.com/qaisjp/jacr-api/pkg/config"
	"github.com/sirupsen/logrus"
)

type Generator struct {
	time.Duration
	Next      <-chan time.Time
	Generator func(s *Statistics) error
}

// Statistics contains all the dependencies of the Statistics Generation server
type Statistics struct {
	Config *config.Config
	Log    *logrus.Logger
	DB     *sqlx.DB

	Generators []*Generator
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

	for _, gen := range []*Generator{DJCountGenerator, TotalTrackVotes} {
		gen.Next = time.After(0)
		s.AddGenerator(gen)
	}

	return s
}

func (a *Statistics) AddGenerator(gen *Generator) {
	a.Generators = append(a.Generators, gen)
}

// Start begins handling all of the statistic generators in the queue.
func (a *Statistics) Start() error {
	for {
		for _, stat := range a.Generators {
			select {
			case <-stat.Next:
				if err := stat.Generator(a); err != nil {
					a.Log.WithFields(logrus.Fields{
						"module": "statistics",
						"error":  err.Error(),
					}).Warnf("Generator failed to run")
				}

				stat.Next = time.After(stat.Duration)
			default:
			}
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
