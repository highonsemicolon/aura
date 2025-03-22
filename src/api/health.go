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

func (h *API) RegisterHealthHandler(router *gin.Engine) {
	router.GET("/readyz", h.readyHandler)
	router.GET("/livez", h.liveHandler)
	router.GET("/infoz", h.infoHandler)
}

func (h *API) infoHandler(c *gin.Context) {
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

func (h *API) readyHandler(c *gin.Context) {
	ready := h.svc.HealthService.Readiness(c.Request.Context())
	c.JSON(200, ready)
}

func (h *API) liveHandler(c *gin.Context) {
	if err := h.svc.HealthService.Liveness(c.Request.Context()); err != nil {
		c.JSON(500, gin.H{"status": "error"})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}
