package api

import (
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	appName    = "unknown"
	version    = "unknown"
	commitHash = "unknown"
	buildTime  = "unknown"
	buildHost  = "unknown"
	startTime  = time.Now()
)

func infoHandler(c *gin.Context) {
	hostname, _ := os.Hostname()
	uptime := time.Since(startTime).String()

	info := gin.H{
		"app_name":      appName,
		"app_env":       os.Getenv("APP_ENV"),
		"version":       version,
		"commit":        commitHash,
		"build_time":    buildTime,
		"build_host":    buildHost,
		"start_time":    startTime.Format(time.RFC3339),
		"uptime":        uptime,
		"hostname":      hostname,
		"go_version":    runtime.Version(),
		"num_cpu":       runtime.NumCPU(),
		"num_goroutine": runtime.NumGoroutine(),
		"num_cgo_call":  runtime.NumCgoCall(),
	}
	c.JSON(200, info)
}
