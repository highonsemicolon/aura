package api

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	api := &API{}
	router := gin.Default()
	router.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	server := NewStrictHandler(api, nil)

	v1 := router.Group("/api/v1")
	{
		RegisterHandlers(v1, server)
	}

	return router
}
