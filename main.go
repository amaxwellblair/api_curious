package main

import (
	"net/http"

	"github.com/amaxwellblair/api_curious/handlers"
)

func main() {
	r := handler.NewRouter()
	// h := handler.NewHandler()
	http.ListenAndServe(":9000", r)
}
