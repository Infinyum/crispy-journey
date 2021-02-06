package main

import (
	"crispy-journey/router"
	"crispy-journey/server"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestStackHome(t *testing.T) {
	is := is.New(t)

	// Use custom "router" package to have all the routes in the same file
	router := router.NewRouter()
	mux := http.NewServeMux()
	router.AddRoutes(mux)

	// Use custom "server" package to configure http.Server with TLS
	s := server.New(mux, serverAddr)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s.Handler.ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK) // Expect status code 200 (OK)
}

func TestHandleHome(t *testing.T) {
	is := is.New(t)

	router := router.NewRouter()

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.HandleHome(w, r)

	is.Equal(w.Code, http.StatusOK) // Expect status code 200 (OK)
}
