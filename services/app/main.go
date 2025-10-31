package main

import (
	"github.com/highonsemicolon/aura/services/app/internal/app"
)

var (
	Version   = "dev"
	Commit    = ""
	BuildTime = ""
	BuiltBy   = ""
)

func main() {
	app.RunApp(Version, Commit, BuildTime, BuiltBy)
}
