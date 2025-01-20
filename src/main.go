package main

import (
	services "aura/src/services/role"
	"time"
)

var fw services.FileWatcher
var pc services.PrivilegeChecker

func init() {
	fw = services.NewFileWatcher("./privileges.yml")
	fw.Load()
	pc = services.NewChecker(fw)

	go fw.Watch()
}

func main() {
	// r := gin.Default()

	for {
		println(pc.IsActionAllowed("editor", "read"))
		time.Sleep(1 * time.Second)
	}

	// api := r.Group("/api")
	// {
	// 	api.GET("/policies", handlers.CheckPermission)
	// }

	// log.Fatal(r.Run(":8080"))
}
