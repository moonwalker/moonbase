package api

import (
	"net/http"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/moonwalker/moonbase/docs"
)

func docs() chi.Router {
	r := chi.NewRouter()

	// r.Use(docsRedirect())
	r.Get("/", httpSwagger.Handler())

	return r
}

func docsRedirect() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/docs" || r.URL.Path == "/docs/" {
				http.Redirect(w, r, "/docs/index.html", http.StatusTemporaryRedirect)
			}
			next.ServeHTTP(w, r)
		})
	}
}
