package api

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	api := &API{}
	router := gin.Default()

	server := NewStrictHandler(api, nil)

	v1 := router.Group("/api/v1")
	{
		RegisterHandlers(v1, server)
	}

	return router
}
