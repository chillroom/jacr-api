package slack

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (i *Impl) CheckOrigin(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	parsedOrigin, err := url.Parse(origin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrapf(err, "failed to verify origin %q", origin),
		})
		return
	}

	if (parsedOrigin.Host == "") || (parsedOrigin.Host == "just-a-chill-room.net") || (parsedOrigin.Host == "www.just-a-chill-room.net") {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Next()
		return
	}

	c.Abort()
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  "error",
		"message": fmt.Sprintf("Invalid origin %q", parsedOrigin.Host),
	})
}
