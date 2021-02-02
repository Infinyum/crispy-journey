// Package router provides all the routing logic for the web server.
package router

import (
	"log"
	"net/http"
	"time"
)

// Router represents a router with embedded logger
type Router struct {
	logger *log.Logger
}

// NewRouter creates a variable of Router type with embedded logger
func NewRouter(logger *log.Logger) *Router {
	return &Router{
		logger: logger,
	}
}

// AddRoutes adds a handler for all the routes of the server
func (router *Router) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", router.Logger(router.HandleHome))
}

// Logger wraps all the handlers to add some logging
func (router *Router) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer router.logger.Printf("Request processed in %s\n", time.Now().Sub(startTime))
		next(w, r)
	}
}

// HandleHome handles home page at path /
func (router *Router) HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world!"))
}
