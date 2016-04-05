package handler

import "net/http"

// IndexHandler handles GET requests to the root
func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	u := h.CurrentUser(r)

	if err := h.templates.ExecuteTemplate(w, "index.html", u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
