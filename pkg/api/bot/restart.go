package bot

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (i *Impl) Restart(c *gin.Context) {
	cmd := exec.Command("pm2", "restart", "jacr-bot")
	cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=" + os.Getenv("HOME")}

	out, err := cmd.CombinedOutput()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errors.Wrap(err, "failed to restart bot").Error(),
			"data":    out,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
