package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/highonsemicolon/aura/internal/middleware"
)

func main() {
	r := gin.Default()
	r.Use(middleware.UserIDMiddleware)

	log.Fatal(r.Run(":8080"))
}
