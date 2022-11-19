package api

import (
	"net/http"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/moonwalker/moonbase/docs"
)

func docs() chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusTemporaryRedirect)
	})

	r.Get("/*", httpSwagger.Handler())

	return r
}
