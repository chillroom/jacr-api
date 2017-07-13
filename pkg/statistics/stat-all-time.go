package statistics

import "time"

var DJCountGenerator = &Generator{
	Duration: time.Hour * 24,
	Generator: func(s *Statistics) error {
		_, err := s.DB.Exec("select count(distinct history.user) from history")
		if err != nil {
			return err
		}

		return nil
	},
}
var TotalTracksPlayed = &Generator{
	Duration: time.Hour * 12,
	Generator: func(s *Statistics) error {
		_, err := s.DB.Exec("SELECT count(id) FROM history")
		if err != nil {
			return err
		}

		return nil
	},
}

// select sum(history.score_up) from history
// select sum(history.score_down) from history
// select sum(history.score_grab) from history
var TotalTrackVotes = &Generator{
	Duration: time.Second,
	Generator: func(s *Statistics) error {
		_, err := s.DB.Exec("SELECT count(*) FROM songs")
		if err != nil {
			return err
		}

		return nil
	},
}

// one time djs
/*
select count(*) from (
    select 1
    from history
    GROUP BY history.user
    HAVING count(history.user) = 1
) as count
*/
