package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os/exec"
)

func restartCheatEndpoint(c *gin.Context) {
	if c.PostForm("codeword") != "341771337" {
		c.Redirect(302, "/")
		return
	}

	cmd := exec.Command("pm2", "restart", "jacr-bot")
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(out)
		c.String(501, "Failed to restart: "+string(out))
		return
	}

	c.String(200, "Bot has been restarted!")
}
