package main

import (
	"time"

	"aura/src/services"
)

var fw *services.FileWatcher

func init() {
	fw = services.NewFileWatcher("./privileges.yml").Load()
	go fw.Watch()
}

func main() {
	// r := gin.Default()
	// r.Use(middleware.UserIDMiddleware)

	// time.Sleep(2 * time.Second)

	var mp = fw.GetEffectivePrivilegesCache()

	for {
		println(services.IsActionAllowed("editor", "read", mp))
		time.Sleep(2 * time.Second)
	}

	// for {
	// 	mp.Range(func(key, value interface{}) bool {
	// 		fmt.Printf("Role: %s: %v\n", key, value)
	// 		return true
	// 	})

	// 	fmt.Println(&mp)
	// 	time.Sleep(2 * time.Second)

	// 	println()
	// }

	// api := r.Group("/api")
	// {
	// 	api.GET("/policies", handlers.CheckPermission)
	// }

	// log.Fatal(r.Run(":8080"))
}
