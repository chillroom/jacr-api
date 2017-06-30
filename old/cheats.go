package old

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func RestartCheatEndpoint(c *gin.Context) {
	if c.PostForm("codeword") != "341771337" {
		c.Redirect(302, "/")
		return
	}

	cmd := exec.Command("pm2", "restart", "jacr-bot")
	cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=" + os.Getenv("HOME")}

	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(string(out))
		c.String(501, "Failed to restart: "+string(out))
		return
	}

	c.String(200, "Bot has been restarted!")
}
