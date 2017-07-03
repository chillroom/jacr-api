package notices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/qaisjp/jacr-api/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/qaisjp/jacr-api/pkg/database"
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
					Title:   notice.Title,
					Message: notice.Message,
				})
			}
		} else if (len(patch.Path) > 1) && (patch.Path[0] == '/') {
			id, err := strconv.Atoi(patch.Path[1:])
			if err == nil {
				if patch.Op == "remove" {
					removedNotices = append(removedNotices, id)
					success = true
				} else if patch.Op == "replace" {
					var notice models.Notice

					err := json.Unmarshal(patch.Value, &notice)

					success = (err == nil) && (notice.Title != "") && (notice.Message != "")
					if success {
						notice.ID = id
						replacedNotices = append(replacedNotices, notice)
					}
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

	err := database.RunInTransaction(i.DB, func(tx *sqlx.Tx) error {
		// Replacements
		if len(replacedNotices) > 0 {
			sqlStr := `
				update notices as n
				set
					title=c.title,
					message=c.message
				from (values`

			vals := make([]interface{}, len(replacedNotices)*2)

			for i, row := range replacedNotices {
				num := (i * 2)

				sqlStr += fmt.Sprintf(" (%#v, $%d, $%d),", row.ID, num+1, num+2)

				vals[num] = row.Title
				vals[num+1] = row.Message
				// vals = append(vals, row.ID, row.Title, row.Message)
			}

			// trim the last ,
			sqlStr = sqlStr[0:len(sqlStr)-1] + `
				) as c(id, title, message)
				where c.id = n.id;
			`

			tx.MustExec(sqlStr, vals...)

			// tx.MustExec(sqlStr,
			fmt.Println("aok")
		}

		if len(removedNotices) > 0 {
			// Removals: make the []int an []interface{} so that the query method can use it
			vals := make([]interface{}, len(removedNotices))

			for i, v := range removedNotices {
				vals[i] = v
			}

			// Removals: perform actual query
			query := "DELETE FROM notices WHERE false "
			for i := range removedNotices {
				query += fmt.Sprintf(" or (id = $%d)", i+1)
			}

			tx.MustExec(query, vals...)
		}

		// Creations
		if len(newNotices) > 0 {
			sqlStr := "INSERT INTO notices(title, message) VALUES "
			vals := make([]interface{}, len(newNotices)*2)

			for i, row := range newNotices {
				num := (i * 2)

				sqlStr += fmt.Sprintf("($%d, $%d),", num+1, num+2)
				vals[num] = row.Title
				vals[num+1] = row.Message
			}

			// trim the last ,
			sqlStr = sqlStr[0 : len(sqlStr)-1]

			fmt.Println(sqlStr)

			tx.MustExec(sqlStr, vals...)
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
