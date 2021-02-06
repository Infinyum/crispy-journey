// Package router provides all the routing logic for the web server.
package router

import (
	"net/http"
)

// Router represents a router with all its routes and associated handlers
type Router struct{}

// NewRouter creates a variable of Router type for external access to the below methods
func NewRouter() *Router {
	return &Router{}
}

// AddRoutes adds a handler for all the routes of the server
func (router *Router) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", router.HandleHome)
}

// HandleHome handles path "/"
func (router *Router) HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello POLYTECH!"))
}
