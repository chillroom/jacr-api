package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

var stmt *pg.Stmt

func historyListEndpoint(c *gin.Context) {

	var results []struct {
		Fkid      string    `json:"-"`
		Name      string    `json:"-"`
		Type      string    `json:"-"`
		Time      time.Time `json:"time"`
		Username  string    `json:"-"`
		UserID    int       `json:"-"`
		ScoreUp   int       `json:"-"`
		ScoreGrab int       `json:"-"`
		ScoreDown int       `json:"-"`

		Song struct {
			Name string `json:"name"`
		} `json:"song"`
		Score struct {
			Down int `json:"down"`
			Grab int `json:"grab"`
			Up   int `json:"up"`
		} `json:"score"`
		User struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
	}

	_, err := db.Query(&results, `
		SELECT
			songs.fkid, songs.name, songs.type,
			dubtrack_users.id as user_id, dubtrack_users.username,
			history.score_up, history.score_down, history.score_grab, history.time
		FROM history, songs, dubtrack_users
		WHERE
		(history.song = songs.id) and
		(history."user" = dubtrack_users.id)
		ORDER BY history.time DESC LIMIT 100
		`)

	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "could not get the history",
			"message": err.Error(),
		})
		return
	}

	for i, result := range results {
		results[i].Song.Name = result.Name

		results[i].Score.Up = result.ScoreUp
		results[i].Score.Grab = result.ScoreGrab
		results[i].Score.Down = result.ScoreDown

		results[i].User.ID = result.UserID
		results[i].User.Username = result.Username
	}

	c.JSON(200, gin.H{
		"status": 200,
		"code":   "success",
		"data":   &results,
	})
}

func historyUserListEndpoint(c *gin.Context) {
	// res, err := r.
	// 	Table("history").
	// 	Filter(r.Row.Field("user").Eq(c.Param("user"))).
	// 	OrderBy(r.Desc("time")).
	// 	Without("platformID").
	// 	Merge(map[string]interface{}{
	// 		"song": r.Table("songs").Get(r.Row.Field("song")).Pluck("name", "fkid", "type").Default(nil),
	// 	}).
	// 	Limit(500).
	// 	Run(rethinkSession)

	// if err != nil {
	// 	c.JSON(200, gin.H{
	// 		"status":  500,
	// 		"code":    "could not query history",
	// 		"message": err.Error(),
	// 	})
	// 	return
	// }
	// defer res.Close()

	// var list []interface{}
	// err = res.All(&list)
	// if err != nil {
	// 	c.JSON(200, gin.H{
	// 		"status":  500,
	// 		"code":    "could not receive history from cursor",
	// 		"message": err.Error(),
	// 	})
	// 	return
	// }

	// c.JSON(200, gin.H{
	// 	"status": 200,
	// 	"code":   "success",
	// 	"data":   list,
	// })
}
