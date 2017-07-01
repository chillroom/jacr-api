package auth

import "github.com/gin-gonic/gin"

func (i *Impl) Authorize(userId int, c *gin.Context) bool {
	if userId == 5 {
		return true
	}

	return false
}
