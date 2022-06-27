package mux

import (
	"net/http"
	"strings"
)

// Route represents a route.
type Route struct {
	Path string
	Func func(http.ResponseWriter, *http.Request)
}

// Mux is a multiplexer for HTTP handlers.
type Mux struct {
	Routes []Route
}

// HandleFunc adds a handler function for the given path.
func (m *Mux) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	m.Routes = append(m.Routes, Route{Path: path, Func: f})
}

// ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range m.Routes {
		if strings.HasPrefix(r.URL.Path, route.Path) {
			route.Func(w, r)
			return
		}
	}
}
