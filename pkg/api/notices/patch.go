package notices

import (
	"net/http"

	"encoding/json"

	"fmt"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
	"github.com/qaisjp/jacr-api/pkg/models"
)

type NoticePatch struct {
	Op    string
	Path  string
	Value json.RawMessage
}

func Patch(c *gin.Context) {
	db := c.Keys["db"].(*pg.DB)

	patches := []NoticePatch{}
	c.BindJSON(&patches)

	newNotices := []models.Notice{}
	removedNotices := []string{}
	replacedNotices := []models.Notice{}

	for _, patch := range patches {
		success := false
		if (patch.Op == "add") && (patch.Path == "/-") {
			var notice *models.Notice
			err := json.Unmarshal(patch.Value, notice)

			success = (err == nil) && (notice != nil) && (notice.Message != "") && (notice.Title != "")
			if success {
				newNotices = append(newNotices, models.Notice{
					Title:   slug.Make(notice.Title),
					Message: notice.Message,
				})
			}
		} else if (len(patch.Path) > 1) && (patch.Path[0] == '/') {
			if patch.Op == "remove" {
				removedNotices = append(removedNotices, patch.Path[1:])
				success = true
			} else if patch.Op == "replace" {
				var message string
				err := json.Unmarshal(patch.Value, &message)

				success = (err == nil) && (message != "")
				if success {
					replacedNotices = append(replacedNotices, models.Notice{
						Title:   patch.Path[1:],
						Message: message,
					})
				}
			}
		} else {
			success = false
		}

		if !success {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Could not parse patch",
				"data":    patch,
			})
			return
		}
	}

	removedNoticesInterface := make([]interface{}, len(removedNotices))
	for i, v := range removedNotices {
		fmt.Println(v)
		removedNoticesInterface[i] = v
	}

	query := "DELETE FROM notices WHERE false " + strings.Repeat(" or (title = ?)", len(removedNotices))
	_, err := db.Exec(query, removedNoticesInterface...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrap(err, "could not delete notices").Error(),
		})
		return
	}

	// _, err := db.Query(&, `SELECT * FROM notices`)

	// c.JSON(http.StatusOK, gin.H{
	// 	"status": "success",
	// 	"data":   &notices,
	// })
}
