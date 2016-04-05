package handler

import (
	"text/template"

	"github.com/amaxwellblair/api_curious/stores"
	"github.com/gorilla/mux"
)

// Handler contains variables that are shared across routes
type Handler struct {
	Secrets   map[string]string
	Store     *store.Store
	templates *template.Template
}

// NewHandler returns a new instance of Handler
func NewHandler(secrets map[string]string) *Handler {
	s := store.NewStore("db/dev_db.db")
	if err := s.Open(); err != nil {
		panic(err)
	}
	return &Handler{
		Secrets:   secrets,
		Store:     s,
		templates: Templates(),
	}
}

// NewRouter creates a new router
func (h *Handler) NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", h.IndexHandler).
		Methods("GET")
	r.HandleFunc("/login", h.NewSessionHandler).
		Methods("GET")
	r.HandleFunc("/session/handshake", h.OAuthHandshakeHandler).
		Methods("GET")
	r.HandleFunc("/session/create", h.CreateSessionHandler).
		Methods("GET")

	return r
}

// Templates returns the parsed templates for all views
func Templates() *template.Template {
	return template.Must(template.ParseFiles("views/index.html"))
}
