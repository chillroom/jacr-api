package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/qaisjp/jacr-api/pkg/models"
)

func (i *Impl) Authorize(userId int, c *gin.Context) bool {
	var u models.User

	err := i.DB.Model(&u).Where("id = ?", userId).Select()
	if err != nil {
		return false
	}

	c.Set("user", u)

	return (u.Level > 1) && (!u.Banned) && (u.Activated)
}
