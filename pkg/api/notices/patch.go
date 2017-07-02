package notices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	goqu "gopkg.in/doug-martin/goqu.v4"

	"github.com/qaisjp/jacr-api/pkg/models"

	"github.com/gin-gonic/gin"
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

	tx, err := i.GQ.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Could not start transaction",
		})
		return
	}

	err = tx.Wrap(func() error {
		// Creations
		if len(newNotices) > 0 {
			_, err := tx.From("notices").Insert(newNotices).Exec()
			if err != nil {
				return err
			}
		}

		if len(removedNotices) > 0 {
			// Removals: make the []string an []interface{} so that the query method can use it
			removedNoticesExpression := make([]goqu.Expression, len(removedNotices))
			for i, v := range removedNotices {
				removedNoticesExpression[i] = goqu.I("id").Eq(v)
			}

			// Removals: perform actual query
			_, err := tx.From("notices").Where(goqu.Or(removedNoticesExpression...)).Delete().Exec()
			if err != nil {
				return err
			}
		}

		if len(replacedNotices) > 0 {
			fmt.Println(replacedNotices)
			// Replacements: perform query
			for _, notice := range replacedNotices {
				_, err := tx.From("notices").Update(&notice).Exec()
				// _, err := i.DB.Model(&replacedNotices).Column("title", "message").Update()
				if err != nil {
					return err
				}
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
