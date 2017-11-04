package history

import (
	"net/http"
	"time"

	"github.com/qaisjp/jacr-api/pkg/models"

	"strconv"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (i *Impl) List(c *gin.Context) {

	var history = make([]struct {
		Time       time.Time         `db:"time"`
		Name       string            `db:"name"`
		Username   string            `db:"username"`
		Fkid       string            `db:"fkid"`
		UserID     int               `db:"user_id"`
		SongID     int               `db:"song_id"`
		SongType   models.SongType   `db:"song_type"`
		ScoreUp    int               `db:"score_up"`
		ScoreDown  int               `db:"score_down"`
		ScoreGrab  int               `db:"score_grab"`
		SkipReason models.SkipReason `db:"skip_reason" json:",omitempty"`
	}, 0)

	query := `
		SELECT
			history.time,
			history.score_up as score_up,
			history.score_down as score_down,
			history.score_grab as score_grab,
			songs.fkid,
			songs.name,
			users.username,
			users.id as user_id,
			songs.id as song_id,
			songs.type as song_type
		FROM history, songs, dubtrack_users as users
		WHERE
			(history.user = users.id) and (history.song = songs.id)
	`

	if userStr := c.Query("user"); userStr != "" {
		uid, err := strconv.Atoi(userStr)
		if (err != nil) || (uid < 1) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid user provided: must be a positive integer",
			})
			return
		}
		query += fmt.Sprintf(`
			and (users.id = %d)
		`, uid)
	}

	if c.Query("filter_op") == "1" {
		query += `
			and is_op(songs.last_play, songs.recent_plays)
		`
	}

	countStr := c.DefaultQuery("count", "100")
	count, err := strconv.Atoi(countStr)
	if (err != nil) || (count < 1) || (count > 100) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid count provided: must be an integer between 1 and 100",
		})
		return
	}

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if (err != nil) || (offset < 0) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid offset provided: must be a positive integer",
		})
		return
	}

	query += fmt.Sprintf(`
		ORDER BY time DESC
		LIMIT %d
		OFFSET %d
	`, count, offset)

	err = i.DB.Select(&history, query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrap(err, "could not receive history").Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   history,
	})
}
