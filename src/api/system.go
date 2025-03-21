package api

import (
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type systemHandler struct{}

func NewSystemHandler() *systemHandler {
	return &systemHandler{}
}

func (h *systemHandler) Register(router *gin.Engine) {
	router.GET("/readyz", readyHandler)
	router.GET("/livez", liveHandler)
	router.GET("/infoz", infoHandler)
}

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
		"name": appName,
		"env":  os.Getenv("APP_ENV"),

		"version": version,
		"build": map[string]any{
			"commit":     commitHash,
			"build_time": buildTime,
			"build_host": buildHost,
		},
		"deployment": map[string]any{
			"go_version":    runtime.Version(),
			"arch":          runtime.GOARCH,
			"start_time":    startTime.Format(time.RFC3339),
			"uptime":        uptime,
			"hostname":      hostname,
			"num_cpu":       runtime.NumCPU(),
			"num_goroutine": runtime.NumGoroutine(),
			"num_cgo_call":  runtime.NumCgoCall(),
		},
	}
	c.JSON(200, info)
}

func readyHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func liveHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
