package server

import (
	"context"
	"net/http"
	"time"
)

type HttpServer struct {
	server Server
}

func NewServer(addr string, handler http.Handler) *HttpServer {

	return &HttpServer{
		server: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       15 * time.Second,
			MaxHeaderBytes:    1 << 20,
		},
	}
}

func (s *HttpServer) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
