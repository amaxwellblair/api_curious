package handler

import (
	"net/http"
	"text/template"

	"github.com/amaxwellblair/api_curious/stores"
	"github.com/gorilla/handlers"
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
func (h *Handler) NewRouter() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", h.IndexHandler).
		Methods("GET")
	r.HandleFunc("/login", h.NewSessionHandler).
		Methods("GET")
	r.HandleFunc("/logout", h.DestroySessionHandler).
		Methods("DELETE")
	r.HandleFunc("/session/handshake", h.OAuthHandshakeHandler).
		Methods("GET")

	return handlers.HTTPMethodOverrideHandler(r)
}

// User holds a token and existence
type User struct {
	Exists bool
	Token  string
}

// NewUser returns a new instance of User
func (h *Handler) NewUser(exist bool, token string) *User {
	return &User{Exists: exist, Token: token}
}

// CurrentUser returns the token linked to the cookie
func (h *Handler) CurrentUser(r *http.Request) *User {
	c, err := r.Cookie("user_id")
	if err != nil {
		return h.NewUser(false, "")
	}
	return h.NewUser(true, c.Value)
}

// Templates returns the parsed templates for all views
func Templates() *template.Template {
	return template.Must(template.ParseFiles("views/index.html"))
}
