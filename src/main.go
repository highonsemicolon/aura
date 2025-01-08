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
	if err := fw.Start(); err != nil {
		log.Fatal("Error starting file watcher:", err)
		return
	}
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
