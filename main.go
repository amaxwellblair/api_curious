package main

import (
	"net/http"

	"github.com/amaxwellblair/api_curious/handlers"
)

func main() {
	h := handler.NewHandler(secrets())
	r := h.NewRouter()
	// h := handler.NewHandler()
	http.ListenAndServe(":9000", r)
}
