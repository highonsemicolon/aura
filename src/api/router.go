package api

import (
	"github.com/gin-gonic/gin"
	"github.com/highonsemicolon/aura/config"
	"github.com/highonsemicolon/aura/src/dal"
	"github.com/highonsemicolon/aura/src/service"
)

func NewRouter() *gin.Engine {
	config := config.GetConfig()
	db := dal.NewMySQLDAL(config.MySQL)
	objectRepo := dal.NewObjectRepository(db, config.Tables["objects"])

	api := &API{
		object: service.NewObjectService(objectRepo),
	}
	router := gin.Default()
	router.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	server := NewStrictHandler(api, nil)

	v1 := router.Group("/api/v1")
	{
		RegisterHandlers(v1, server)
	}

	return router
}
