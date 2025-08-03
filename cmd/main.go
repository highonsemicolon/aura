package main

import (
	"log"

	"github.com/highonsemicolon/aura/internal/config"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	log.Printf("Service Name: %s", cfg.ServiceName)

}
