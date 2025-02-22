package server

import (
	"context"
	"errors"
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
	router := api.NewRouter()

	return &Server{
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           router,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       15 * time.Second,
			MaxHeaderBytes:    1 << 20,
		},
	}
}

func (s *Server) StartAndWait() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

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
