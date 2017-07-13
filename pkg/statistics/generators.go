package statistics

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Generator struct {
	Name      string
	Query     string
	Frequency time.Duration
}

func (g *Generator) Spawn(queue chan *Generator) {
	for {
		queue <- g
		time.Sleep(g.Frequency)
	}
}

func (g *Generator) Run(db *sqlx.DB) error {
	query := fmt.Sprintf(`
		with result as (%s)
		insert into statistics (name, value)
		select $1, to_json(value) from result
		ON CONFLICT (name) DO UPDATE SET value = excluded.value
	`, g.Query)

	_, err := db.Exec(query, g.Name)
	return err
}

func GetGenerators() []*Generator {
	return []*Generator{
		{
			Name:      "user-count",
			Query:     "select count(distinct history.user) as value from history",
			Frequency: time.Hour,
		},

		{
			Name:      "history-count",
			Query:     "select count(id) as value from history",
			Frequency: time.Hour,
		},

		{
			Name:      "total-upvotes",
			Query:     "select sum(history.score_up) as value from history",
			Frequency: time.Hour * 6,
		},

		{
			Name:      "total-downvotes",
			Query:     "select sum(history.score_down) as value from history",
			Frequency: time.Hour * 6,
		},

		{
			Name:      "total-grabs",
			Query:     "select sum(history.score_grab) as value from history",
			Frequency: time.Hour * 6,
		},

		{
			Name:      "total-songs",
			Query:     "select count(id) as value from songs",
			Frequency: time.Hour,
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
			Frequency: time.Hour * 12,
		},

		{
			Name: "top-karma",
			Query: `
				select array_to_json(array_agg(t.value)) as value from (
					select json_build_object(
						'ID', id,
						'Username', username,
						'Karma', karma
					) as value
					from dubtrack_users
					order by karma desc
					limit 3
				) as t
			`,
			Frequency: time.Hour * 24,
		},

		{
			Name: "newest-song",
			Query: `
				select json_build_object('ID', songs.id, 'Fkid', songs.fkid, 'Name', songs.name) as value
				from songs
				where (total_plays = 1) order by last_play desc limit 1
			`,
			Frequency: time.Minute * 5,
		},

		{
			Name: "most-upvoted-song",
			Query: `
				with score as (select max(score_up) from history)
				select
				json_build_object(
					'ScoreUp', score.max,
					'Name', songs.name,
					'Fkid', songs.fkid,
					'ID', songs.id
				) as value
				from score, history, songs
				where (history.score_up = score.max) and (history.song = songs.id)
				order by history.time desc
				limit 1
			`,
			Frequency: time.Hour * 24,
		},

		{
			Name: "most-grabbed-song",
			Query: `
				with score as (select max(score_grab) from history)
				select
				json_build_object(
					'ScoreGrab', score.max,
					'Name', songs.name,
					'Fkid', songs.fkid,
					'ID', songs.id
				) as value
				from score, history, songs
				where (history.score_grab = score.max) and (history.song = songs.id)
				order by history.time desc
				limit 1
				`,
			Frequency: time.Hour * 24,
		},

		{
			Name: "songs-played-once",
			Query: `
				select count(id) as value from songs where total_plays = 1
			`,
			Frequency: time.Minute * 5,
		},

		{
			Name: "user-playing-most-tracks",
			Query: `
				select json_build_object(
				'ID', id,
				'Username', username,
				'Count', history.count
				) as value
				from dubtrack_users, (
				SELECT
					history.user,
					count(history)
				FROM history
				GROUP BY history.user
				) as history
				where history.user = dubtrack_users.id
				order by count desc
				limit 1
			`,
			Frequency: time.Hour * 24,

			// Alternate, slower, version
			/*
				Query: `
					SELECT
					json_build_object(
						'ID', users.id,
						'Username', username,
						'Count', history.count
					) as value
					FROM history, dubtrack_users as users
					where users.id = history.user
					GROUP BY history.user, users.username, users.id
					order by count(history) desc
					limit 1
				`,
			*/
		},
	}
}
