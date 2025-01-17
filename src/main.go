package main

import services "aura/src/services/privilege"

var fw *services.FileWatcher

func init() {
	fw = services.NewFileWatcher("./privileges.yml").Load()
	go fw.Watch()
}

func main() {
	// r := gin.Default()
	// r.Use(middleware.UserIDMiddleware)

	for {
		println(services.IsActionAllowed("editor", "read", fw))
		// time.Sleep(1 * time.Second)
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
