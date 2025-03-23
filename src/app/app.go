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
	server server.Server
	db     dal.Database
}

func NewApp(cfg *config.Config, db dal.Database) *App {

	repos := dal.NewDalContainer(db, cfg.Tables)
	services := service.NewServiceContainer(repos)
	api := api.NewAPI(services)

	return &App{
		server: setupServer(cfg, api),
		db:     db,
	}
}

func setupServer(cfg *config.Config, api *api.API) *server.HttpServer {
	return server.NewServer(cfg.Address, api.NewRouter())
}

func (app *App) Run(ctx context.Context) {
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

	app.gracefulShutdown(ctx)
}

func (app *App) gracefulShutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	if err := app.db.Close(); err != nil {
		log.Printf("db close error: %v", err)
	}

	log.Println("server shutdown gracefully")
}
