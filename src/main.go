package main

import (
	"log"

	"aura/src/handlers"
	"aura/src/middleware"
	"aura/src/services"

	"github.com/gin-gonic/gin"
)

var fw *services.FileWatcher

func init() {
	fw = services.NewFileWatcher("./privileges.yml")
	go fw.Start()

}

func main() {
	r := gin.Default()
	r.Use(middleware.UserIDMiddleware)

	api := r.Group("/api")
	{
		api.GET("/policies", handlers.CheckPermission)
	}

	log.Fatal(r.Run(":8080"))
}
