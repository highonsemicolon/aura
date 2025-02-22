package main

import (
	"log"
	"net/http"

	"github.com/highonsemicolon/aura/src/server"
)

func main() {
	srv := server.NewServer(":8080")
	defer server.HandleShutdown(srv)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

}
