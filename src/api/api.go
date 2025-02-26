package api

import (
	"github.com/gin-gonic/gin"
	"github.com/highonsemicolon/aura/src/service"
)

func NewAPI(services *service.ServiceContainer) *API {
	return &API{svc: services}
}

func (api *API) NewRouter() *gin.Engine {

	router := gin.Default()
	router.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	server := NewStrictHandler(api, nil)

	v1 := router.Group("/api/v1")
	{
		RegisterHandlers(v1, server)
	}

	return router
}
