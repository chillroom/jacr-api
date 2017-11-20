package songs

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

	var songs = make([]struct {
		ID           int               `db:"id"`
		Fkid         string            `db:"fkid"`
		Name         string            `db:"name"`
		LastPlay     time.Time         `db:"last_play"`
		SkipReason   models.SkipReason `db:"skip_reason" json:",omitempty"`
		RecentPlays  int               `db:"recent_plays"`
		TotalPlays   int               `db:"total_plays"`
		SongType     models.SongType   `db:"type"`
		Retagged     bool              `db:"retagged"`
		AutoRetagged bool              `db:"autoretagged"`
	}, 0)

	query := `
		SELECT *
		FROM songs
	`

	if c.Query("filter_op") == "1" {
		query += `
			WHERE is_op(songs.last_play, songs.recent_plays)
		`
	} else {
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
			ORDER BY name
			LIMIT %d
			OFFSET %d
		`, count, offset)
	}

	err := i.DB.Select(&songs, query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrap(err, "could not receive songs").Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   songs,
	})
}
