package main

import (
	"aura/src/api"
	"aura/src/db"
	"aura/src/services"
	"aura/src/utils"
	"log"

	"github.com/gin-gonic/gin"

	role "aura/src/services/role"
)

var fw role.FileWatcher
var pc role.PrivilegeChecker

var cfg *utils.Config

func init() {

	cfg = utils.LoadConfig("./config.yml")
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

	repo := db.NewSqlDB(cfg.MySQL.DSN, cfg.MySQL.CACertPath)
	defer repo.Close()

	service := services.NewPrivilegeService(pc, repo)
	handler := api.NewPrivilegeHandler(service)

	api.Register(r, handler)

	log.Fatal(r.Run(":8080"))
}
