package main

import (
	"crispy-journey/router"
	"crispy-journey/server"
	"log"
	"net/http"
)

const serverAddr = ":8080"

func main() {
	log.SetPrefix("[CRISPY-JOURNEY] ")
	log.SetFlags(log.Ltime)

	// Use custom "router" package to have all the routes in the same file
	router := router.NewRouter()
	mux := http.NewServeMux()
	router.AddRoutes(mux)

	// Use custom "server" package to configure http.Server with TLS
	s := server.New(mux, serverAddr)

	log.Printf("Server starting on %s\n", serverAddr)
	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
