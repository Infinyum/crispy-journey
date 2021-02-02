package main

import (
	"crispy-journey/router"
	"crispy-journey/server"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestStackHome(t *testing.T) {
	is := is.New(t)
	logger := log.New(os.Stdout, "[TESTING] ", log.Ltime|log.Lshortfile)

	// Set up the server
	router := router.NewRouter(logger)
	mux := http.NewServeMux()
	router.AddRoutes(mux)
	s := server.New(mux, serverAddr)

	// Create and send the request
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s.Handler.ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK) // Expect status code 200
}

func TestHandleHome(t *testing.T) {
	is := is.New(t)

	// Set up the server
	router := router.NewRouter(nil)

	// Create and send the request
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.HandleHome(w, r)

	is.Equal(w.Code, http.StatusOK) // Expect status code 200
}
