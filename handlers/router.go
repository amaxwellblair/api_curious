package handler

import (
	"text/template"

	"github.com/gorilla/mux"
)

// Handler contains variables that are shared across routes
type Handler struct {
	templates *template.Template
}

// NewHandler returns a new instance of Handler
func NewHandler() *Handler {
	return &Handler{
		templates: Templates(),
	}
}

// NewRouter creates a new router
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	h := NewHandler()
	r.HandleFunc("/", h.IndexHandler).
		Methods("GET")
	r.HandleFunc("/login", h.NewSessionHandler).
		Methods("GET")
	r.HandleFunc("/session/create", h.CreateSessionHandler).
		Methods("GET")

	return r
}

// Templates returns the parsed templates for all views
func Templates() *template.Template {
	return template.Must(template.ParseFiles("views/index.html"))
}
