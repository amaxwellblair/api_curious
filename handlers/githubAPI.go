package handler

import (
	"fmt"
	"net/http"
)

//

// UserGithubAPI gets a user object from Github
func (h *Handler) UserGithubAPI(w http.ResponseWriter, r *http.Request) {
	// Parse URL params and retrieve token
	hash := r.URL.Query().Get("user_id")
	token, err := h.Store.User(hash, h.Secrets["saltSecret"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send request to Github API
	data, err := h.Client.GithubGetUser(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	fmt.Fprint(w, data)
}

// GistGithubAPI gets gists for a specific user from Github
func (h *Handler) GistGithubAPI(w http.ResponseWriter, r *http.Request) {
	// Parse URL params and retrieve token
	hash := r.URL.Query().Get("user_id")
	username := r.URL.Query().Get("username")
	token, err := h.Store.User(hash, h.Secrets["saltSecret"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send request to Github API
	data, err := h.Client.GithubGetGists(token, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	fmt.Fprint(w, data)
}
