package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/highonsemicolon/aura/utils/logger"
)

func NewRouter() *gin.Engine {
	api := &API{}
	router := gin.New()
	router.Use(gin.Recovery())

	logger := logger.InitLogger(logger.Std)

	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
		)
	})

	server := NewStrictHandler(api, nil)

	v1 := router.Group("/api/v1")
	{
		RegisterHandlers(v1, server)
	}

	return router
}
