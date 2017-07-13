package statistics

import "time"

type Generator struct {
	Name     string
	Query    string
	Duration time.Duration
	Next     <-chan time.Time
}

func GetGenerators() []*Generator {
	return []*Generator{
		{
			Name:     "dj-count",
			Query:    "select count(distinct history.user) as value from history",
			Duration: time.Hour,
		},

		{
			Name:     "history-count",
			Query:    "select count(id) as value from history",
			Duration: time.Hour,
		},

		{
			Name:     "total-upvotes",
			Query:    "select sum(history.score_up) as value from history",
			Duration: time.Hour * 6,
		},

		{
			Name:     "total-downvotes",
			Query:    "select sum(history.score_down) as value from history",
			Duration: time.Hour * 6,
		},

		{
			Name:     "total-grabs",
			Query:    "select sum(history.score_grab) as value from history",
			Duration: time.Hour * 6,
		},

		{
			Name:     "total-songs",
			Query:    "select count(id) as value from songs",
			Duration: time.Hour,
		},

		{
			Name: "one-time-djs",
			Query: `
				select count(*) as value from (
					select 1
					from history
					GROUP BY history.user
					HAVING count(history.user) = 1
				) as count
			`,
			Duration: time.Hour * 12,
		},
	}
}
