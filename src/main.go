package main

import (
	"log"
	"net/http"

	"github.com/highonsemicolon/aura/src/api"
)

func main() {
	r := api.NewApp()

	log.Fatal(http.ListenAndServe(":8080", r))
}
