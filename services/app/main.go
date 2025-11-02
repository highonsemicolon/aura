package main

import (
	"log"

	"github.com/highonsemicolon/aura/services/app/internal/app"
)

var (
	Version   = "dev"
	Commit    = ""
	BuildTime = ""
	BuiltBy   = ""
)

func main() {
	if err := app.RunApp(Version, Commit, BuildTime, BuiltBy); err != nil {
		log.Fatal(err)
	}
}
