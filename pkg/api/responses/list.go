package responses

import (
	"net/http"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (i *Impl) List(c *gin.Context) {
	list := []struct {
		Cmds     []byte `db:"cmds" pg:",array"`
		Messages []byte `db:"messages" pg:",array"`
	}{}

	err := i.DB.Select(&list,
		`SELECT array_to_json(array_agg(cmds.name)) as cmds, array_to_json(groups.messages) as messages FROM
			response_commands as cmds,
			response_groups as groups
		WHERE
			cmds.group = groups.id
		GROUP BY groups.messages`,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrap(err, "could not receive responses").Error(),
		})
		return
	}

	outList := make([]struct {
		Cmds     []string
		Messages []string
	}, len(list))

	for i, response := range list {
		out := &outList[i]

		err := json.Unmarshal(response.Cmds, &out.Cmds)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(response.Messages, &out.Messages)
		if err != nil {
			panic(err)
		}
	}

	// for _, r := range list {
	// 	i.Log.Info(string(r.Cmds), string(r.Messages))
	// }

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   outList,
	})
}
