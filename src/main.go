package main

import (
	"context"

	"github.com/highonsemicolon/aura/config"
	"github.com/highonsemicolon/aura/src/app"
	"github.com/highonsemicolon/aura/src/dal"
)

func main() {
	cfg := config.GetConfig()

	db := dal.NewMySQLDAL(cfg.MySQL)
	// permissionRepo := dal.NewRelationshipRepository(db, cfg.Tables["relationships"])
	// handlePermissionOperations(permissionRepo)

	// srv := server.NewServer(cfg.Address)
	// defer srv.Shutdown()

	// srv.StartAndWait()

	ctx := context.Background()

	app := app.NewApp(cfg, db)
	app.Run(ctx)
}
