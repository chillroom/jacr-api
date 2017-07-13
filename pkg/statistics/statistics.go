package statistics

import (
	"context"
	"fmt"

	"time"

	"github.com/jmoiron/sqlx"
	"github.com/qaisjp/jacr-api/pkg/config"
	"github.com/sirupsen/logrus"
)

type Generator struct {
	Name     string
	Query    string
	Duration time.Duration
	Next     <-chan time.Time
}

// Statistics contains all the dependencies of the Statistics Generation server
type Statistics struct {
	Config *config.Config
	Log    *logrus.Logger
	DB     *sqlx.DB

	Generators []*Generator
	// Queue      chan *Generator
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

	s.AddGenerators()

	// Initialise each generator by running them
	for _, gen := range s.Generators {
		gen.Next = time.After(time.Second * 5)
	}

	return s
}

// Start begins handling all of the statistic generators in the queue.
func (a *Statistics) Start() error {
	for {
		for _, stat := range a.Generators {
			select {
			case <-stat.Next:
				query := fmt.Sprintf(`
					with result as (%s)
					insert into statistics (name, value)
					select $1, to_json(value) from result
					ON CONFLICT (name) DO UPDATE SET value = excluded.value
				`, stat.Query)

				_, err := a.DB.Exec(query, stat.Name)

				if err != nil {
					a.Log.WithFields(logrus.Fields{
						"module":    "statistics",
						"error":     err.Error(),
						"generator": stat.Name,
					}).Warn("Generator failed to run")
				} else {
					a.Log.WithFields(logrus.Fields{
						"module":    "statistics",
						"generator": stat.Name,
					}).Info("Generator succeeded")
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
