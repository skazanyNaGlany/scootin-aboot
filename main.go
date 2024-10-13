package main

import (
	"log"
	"net/http"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

// Initializes the logger with the desired flags.
func initLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// Entry point of the application.
// It initializes the logger, initializes the API, adds scooters, and starts the server.
func main() {
	initLogger()

	_, router := InitAPI()

	log.Println("Starting server on :80")
	http.ListenAndServe(":80", router)
}
