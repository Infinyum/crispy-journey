package main

import (
	"crispy-journey/routes"
	"crispy-journey/server"
	"log"
	"net/http"
	"os"
)

const serverAddr = "localhost:8080"

func main() {
	logger := log.New(os.Stdout, "[CRISPY-JOURNEY] ", log.Ltime|log.Lshortfile)

	// Set up the router
	router := routes.NewRouter(logger)
	mux := http.NewServeMux()
	router.AddRoutes(mux)

	s := server.New(mux, serverAddr)

	logger.Printf("Server starting at %s\n", serverAddr)
	err := s.ListenAndServe()
	if err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
