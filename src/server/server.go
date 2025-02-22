package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/highonsemicolon/aura/src/api"
)

func NewServer(addr string) *http.Server {
	r := api.NewApp()

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

func HandleShutdown(srv *http.Server) {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
		os.Exit(1)
	}

	log.Println("server exited gracefully")
}
