package main

import (
	"github.com/highonsemicolon/aura/config"
	"github.com/highonsemicolon/aura/src/server"
)

func main() {
	config := config.GetConfig()

	srv := server.NewServer(config.Address)
	defer srv.Shutdown()

	srv.StartAndWait()
}
