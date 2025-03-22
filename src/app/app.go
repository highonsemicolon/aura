package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/highonsemicolon/aura/config"
	"github.com/highonsemicolon/aura/src/api"
	"github.com/highonsemicolon/aura/src/dal"
	"github.com/highonsemicolon/aura/src/server"
	"github.com/highonsemicolon/aura/src/service"
)

type App struct {
	server *server.Server
	repos  *dal.DalContainer
}

func NewApp() *App {
	config := config.GetConfig()

	db := dal.NewMySQLDAL(config.MySQL)
	repos := dal.NewDalContainer(db, config.Tables)
	services := service.NewServiceContainer(repos)
	api := api.NewAPI(services)

	return &App{
		server: setupServer(config, api),
		repos:  repos,
	}
}

func setupServer(cfg *config.Config, api *api.API) *server.Server {
	return server.NewServer(cfg.Address, api.NewRouter())
}

func (app *App) Run() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	serverErr := make(chan error, 1)

	go func() {
		if err := app.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		log.Printf("server error: %v\n", err)
	case <-quit:
		log.Println("received signal to shutdown server")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}

	if err := app.repos.Close(); err != nil {
		log.Fatalf("db close error: %v", err)
	}

	log.Println("server shutdown gracefully")
}
