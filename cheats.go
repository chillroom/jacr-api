package main

import (
	"github.com/gin-gonic/gin"
	"os/exec"
)

func restartCheatEndpoint(c *gin.Context) {
	if c.PostForm("codeword") != "341771337" {
		c.Redirect(302, "/")
		return
	}

	cmd := exec.Command("pm2", "restart", "jacr-api")
	err := cmd.Run()

	if err != nil {
		c.String(501, "Failed to restart :(")
		return
	}

	c.String(200, "Bot has been restarted!")
}
