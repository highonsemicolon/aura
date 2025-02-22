package main

import (
	"github.com/highonsemicolon/aura/src/server"
)

func main() {
	srv := server.NewServer(":8080")
	defer srv.Shutdown()

	srv.ListenAndServe()
	srv.WaitForShutdown()
}
