package handler

import (
	"net/http"
	"net/url"
	"time"
)

// NewSessionHandler redirects users to the Github API
func (h *Handler) NewSessionHandler(w http.ResponseWriter, r *http.Request) {
	// Build URL
	u := new(url.URL)
	u.Scheme = "https"
	u.Host = "github.com"
	u.Path = "/login/oauth/authorize"
	params := u.Query()
	params.Add("client_id", clientID)
	params.Add("redirect_uri", "http://localhost:9000/session/create")
	params.Add("state", "testing 123")
	u.RawQuery = params.Encode()

	// Redirect to github OAuth
	http.Redirect(w, r, u.String(), http.StatusFound)
}

// CreateSessionHandler creates the authentication token
func (h *Handler) CreateSessionHandler(w http.ResponseWriter, r *http.Request) {

	value := "bacon"
	c := http.Cookie{
		Name:     "thing1",
		Value:    value,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		MaxAge:   360000,
		Path:     "/",
	}
	http.SetCookie(w, &c)
	http.Redirect(w, r, "/", http.StatusFound)
}
