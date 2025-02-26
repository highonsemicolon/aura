package app

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/highonsemicolon/aura/config"
	"github.com/highonsemicolon/aura/src/api"
	"github.com/highonsemicolon/aura/src/dal"
	"github.com/highonsemicolon/aura/src/server"
	"github.com/highonsemicolon/aura/src/service"
)

type App struct {
	server *server.Server
}

func NewApp() *App {
	config := config.GetConfig()
	
	db := dal.NewMySQLDAL(config.MySQL)
	repos := dal.NewDalContainer(db, config.Tables)
	services := service.NewServiceContainer(repos)

	api := api.NewAPI(services)

	srv := server.NewServer(config.Address, api.NewRouter())
	return &App{server: srv}
}

func (app *App) Run() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.server.Start(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %s\n", err)
		}
	}()

	<-quit

	log.Println("shutting down server...")
	if err := app.server.Shutdown(); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
}
