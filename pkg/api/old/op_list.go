package old

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

type song struct {
	ID       int       `json:"id"`
	Fkid     string    `json:"fkid"`
	Name     string    `json:"name"`
	LastPlay time.Time `json:"lastPlay"`
	Type     string    `json:"type"`
	Plays    int       `json:"plays"`
}

func OpListEndpoint(c *gin.Context) {
	results := make([]song, 0)
	db := c.Keys["db"].(*pg.DB)
	_, err := db.Query(&results, `SELECT id, fkid, name, last_play, type, total_plays as plays FROM songs WHERE skip_reason = 'op'`)

	if err != nil {
		c.JSON(200, gin.H{
			"status":  500,
			"code":    "could not get the op list",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"code":   "success",
		"data":   results,
	})
}
