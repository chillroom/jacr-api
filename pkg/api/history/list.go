package history

import (
	"net/http"
	"time"

	"github.com/qaisjp/jacr-api/pkg/models"

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
		SkipReason models.SkipReason `db:"skip_reason" json:",omitempty"`
	}, 0)

	err := i.DB.Select(&history, `
		SELECT
			history.time,
			songs.fkid,
			songs.name,
			users.username,
			users.id as user_id,
			songs.id as song_id,
			songs.type as song_type
		FROM history, songs, dubtrack_users as users
		WHERE (history.user = users.id) and (history.song = songs.id)
		LIMIT 10
	`)

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
