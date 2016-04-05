package handler

import (
	"io/ioutil"
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
	params.Add("client_id", h.Secrets["clientID"])
	params.Add("redirect_uri", "http://localhost:9000/session/handshake")
	params.Add("state", "testing 123")
	u.RawQuery = params.Encode()

	// Redirect to github OAuth
	http.Redirect(w, r, u.String(), http.StatusFound)
}

// OAuthHandshakeHandler is an intermediary step between OAuth and the Server
func (h *Handler) OAuthHandshakeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse params
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	// Build URL
	u := new(url.URL)
	u.Scheme = "https"
	u.Host = "github.com"
	u.Path = "/login/oauth/access_token"
	params := u.Query()
	params.Add("client_id", h.Secrets["clientID"])
	params.Add("client_secret", h.Secrets["clientSecret"])
	params.Add("code", code)
	params.Add("redirect_uri", "http://localhost:9000/session/create")
	params.Add("state", state)
	u.RawQuery = params.Encode()

	// Create request
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Request access token - do not follow redirect
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil && resp.StatusCode != 302 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	params, err = url.ParseQuery(string(body))
	token := params.Get("access_token")

	// Create user with given token
	if err := h.Store.CreateUser(token, h.Secrets["saltSecret"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set cookie
	hashedToken, err := h.Store.DigestToken(token, h.Secrets["saltSecret"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c := http.Cookie{
		Name:     "user_id",
		Value:    hashedToken,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		MaxAge:   360000,
		Path:     "/",
	}
	http.SetCookie(w, &c)

	// Redirect to landing page
	http.Redirect(w, r, "/", http.StatusFound)
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
