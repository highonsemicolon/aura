package main

import (
	"log"

	"aura/src/handlers"
	"aura/src/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(middleware.UserIDMiddleware)

	api := r.Group("/api")
	{
		api.GET("/policies", handlers.CheckPermission)
	}

	log.Fatal(r.Run(":8080"))
}
