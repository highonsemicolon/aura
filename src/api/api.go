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
	NewSystemHandler().Register(router)

	server := NewStrictHandler(api, nil)
	RegisterHandlersWithOptions(router, server, GinServerOptions{BaseURL: "/api/v1"})

	return router
}
