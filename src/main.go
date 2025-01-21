package main

import (
	"aura/src/api"
	"aura/src/db"
	"aura/src/middleware"
	"aura/src/services"
	"log"

	"github.com/gin-gonic/gin"

	role "aura/src/services/role"
)

var fw role.FileWatcher
var pc role.PrivilegeChecker

func init() {
	fw = role.NewFileWatcher("./privileges.yml").Load()
	pc = role.NewChecker(fw)

	go fw.Watch()
}

func main() {
	r := gin.Default()

	// for {
	// 	println(pc.IsActionAllowed("editor", "read"))
	// 	time.Sleep(1 * time.Second)
	// }

	r.Use(middleware.UserIDMiddleware)

	repo := db.NewDB(nil)
	defer repo.Close()

	service := services.NewPrivilegeService(pc, repo)
	handler := api.NewPrivilegeHandler(service)

	api.Register(r, handler)

	log.Fatal(r.Run(":8080"))
}
