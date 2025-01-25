package main

import (
	"aura/src/api"
	"aura/src/db"
	"aura/src/services"
	"aura/src/utils"
	"log"

	_ "aura/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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

// @title		Aura API
// @version	0.1
// @description.markdown
// @host			localhost:8080
// @Schemes		http https
// @contact.name	Onkar Chendage
// @contact.email	onkar.chendage@gmail.com
// @license.name	MIT
// @license.url	https://opensource.org/licenses/MIT
func main() {

	r := gin.Default()
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
