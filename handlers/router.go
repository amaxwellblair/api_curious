package handler

import (
	"net/http"
	"text/template"

	"github.com/amaxwellblair/api_curious/clients"
	"github.com/amaxwellblair/api_curious/stores"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Handler contains variables that are shared across routes
type Handler struct {
	Secrets   map[string]string
	Store     *store.Store
	Client    *client.Client
	templates *template.Template
}

// NewHandler returns a new instance of Handler
func NewHandler(secrets map[string]string) *Handler {
	s := store.NewStore("db/dev_db.db")
	if err := s.Open(); err != nil {
		panic(err)
	}
	c := client.NewClient("api.github.com")
	return &Handler{
		Secrets:   secrets,
		Store:     s,
		Client:    c,
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
	r.HandleFunc("/github/user", h.UserGithubAPI).
		Methods("GET")
	r.HandleFunc("/github/user/gists", h.GistGithubAPI).
		Methods("GET")

	handler := handlers.HTTPMethodOverrideHandler(r)

	o := handlers.AllowedOrigins([]string{"*"})

	handler = handlers.CORS(o)(r)

	return handler
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
	token, err := h.Store.User(c.Value, h.Secrets["saltSecret"])
	if err != nil {
		return h.NewUser(false, "")
	}

	return h.NewUser(true, token)
}

// Templates returns the parsed templates for all views
func Templates() *template.Template {
	return template.Must(template.ParseFiles("views/index.html"))
}
