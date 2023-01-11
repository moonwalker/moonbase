package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func printRouteParams(w http.ResponseWriter, r *http.Request) {
	rctx := chi.RouteContext(r.Context())
	log.Println("pattern:", rctx.RoutePattern())
	for i, k := range rctx.URLParams.Keys {
		v := rctx.URLParams.Values[i]
		if k == "*" && v == "" {
			continue
		}
		log.Println(k, "=>", v)
	}
}
