package models

import "time"

type History struct {
	ID        int       `db:"id" goqu:"skipinsert"`
	DubID     string    `db:"dub_id"`
	ScoreDown int       `db:"score_down"`
	ScoreGrab int       `db:"score_grab"`
	ScoreUp   int       `db:"score_up"`
	Song      int       `db:"song"`
	User      int       `db:"user"`
	Time      time.Time `db:"time"`
}
