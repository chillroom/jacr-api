package notices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/qaisjp/jacr-api/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
)

type NoticePatch struct {
	Op    string
	Path  string
	Value json.RawMessage
}

// Patch allows you to add, remove, or replace notices
func (i *Impl) Patch(c *gin.Context) {

	patches := []NoticePatch{}
	if err := c.BindJSON(&patches); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid JSON input",
		})
		return
	}

	var newNotices []models.Notice
	var removedNotices []int
	var replacedNotices []models.Notice

	for _, patch := range patches {
		success := false
		if (patch.Op == "add") && (patch.Path == "/-") {
			var notice models.Notice
			err := json.Unmarshal(patch.Value, &notice)

			if err != nil {
				i.Log.Warnf(err.Error())
			}

			success = (err == nil) && (notice.Message != "") && (notice.Title != "")
			if success {
				newNotices = append(newNotices, models.Notice{
					Title:   slug.Make(notice.Title),
					Message: notice.Message,
				})
			}
		} else if (len(patch.Path) > 1) && (patch.Path[0] == '/') {
			if patch.Op == "remove" {
				i, err := strconv.Atoi(patch.Path[1:])
				if err == nil {
					removedNotices = append(removedNotices, i)
					success = true
				}
			} else if patch.Op == "replace" {
				var notice models.Notice

				err := json.Unmarshal(patch.Value, &notice)

				success = (err == nil) && (notice.ID != 0) && (notice.Title != "") && (notice.Message != "")
				if success {
					notice.Title = slug.Make(notice.Title)
					replacedNotices = append(replacedNotices, notice)
					success = true
				}
			}
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

	err := i.DB.RunInTransaction(func(tx *pg.Tx) error {
		// Creations
		if len(newNotices) > 0 {
			err := i.DB.Insert(&newNotices)
			if err != nil {
				return err
			}
		}

		if len(removedNotices) > 0 {
			// Removals: make the []string an []interface{} so that the query method can use it
			removedNoticesInterface := make([]interface{}, len(removedNotices))
			for i, v := range removedNotices {
				fmt.Println(v)
				removedNoticesInterface[i] = v
			}

			// Removals: perform actual query
			query := "DELETE FROM notices WHERE false " + strings.Repeat(" or (id = ?)", len(removedNotices))
			_, err := i.DB.Exec(query, removedNoticesInterface...)
			if err != nil {
				return err
			}
		}

		if len(replacedNotices) > 0 {
			fmt.Println(replacedNotices)
			// Replacements: perform query
			_, err := i.DB.Model(&replacedNotices).Column("title", "message").Update()
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrapf(err, "Encountered errors in PATCH").Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
