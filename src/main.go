package main

import "github.com/highonsemicolon/aura/src/app"

func main() {
	// config := config.GetConfig()

	// db := dal.NewMySQLDAL(config.MySQL)
	// defer db.Close()

	// permissionRepo := dal.NewRelationshipRepository(db, config.Tables["relationships"])
	// handlePermissionOperations(permissionRepo)

	// srv := server.NewServer(config.Address)
	// defer srv.Shutdown()

	// srv.StartAndWait()

	app := app.NewApp()
	app.Run()
}
