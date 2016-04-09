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
	params.Add("scope", "gist")
	u.RawQuery = params.Encode()

	// Redirect to github OAuth
	http.Redirect(w, r, u.String(), http.StatusFound)
}

// DestroySessionHandler deletes a session
func (h *Handler) DestroySessionHandler(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:     "user_id",
		Value:    "deleted",
		Expires:  time.Now(),
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
	}
	http.SetCookie(w, &c)

	// Redirect to root
	http.Redirect(w, r, "/", http.StatusFound)
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
	params.Add("redirect_uri", "http://localhost:9000/session/")
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
	hashedToken, err := h.Store.DigestToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build URL
	u = new(url.URL)
	u.Scheme = "http"
	u.Host = "localhost:2015"
	u.Path = "/"
	params = u.Query()
	params.Add("user_id", hashedToken)
	u.RawQuery = params.Encode()

	// Redirect to landing page
	http.Redirect(w, r, u.String(), http.StatusFound)
}
