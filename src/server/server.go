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

type Server struct {
	httpServer *http.Server
}

func NewServer(addr string) *Server {
	r := api.NewApp()

	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: r,
		},
	}
}

func (s *Server) ListenAndServe() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func (s *Server) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

func (s *Server) Shutdown() {

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
		os.Exit(1)
	}

	log.Println("server exited gracefully")
}
