package api

import "github.com/gin-gonic/gin"

func NewApp() *gin.Engine {
	api := &API{}
	r := gin.Default()

	server := NewStrictHandler(api, nil)

	v1 := r.Group("/api/v1")
	{
		RegisterHandlers(v1, server)
	}

	return r

}
