package server

import (
	"compress/flate"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/moonwalker/moonbase/internal/api"
)

func Listen(port int) error {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(flate.DefaultCompression))
	r.Use(cors.Handler(corsOptions))
	r.Use(middleware.Recoverer)

	r.Mount("/", api.Routes())

	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, r)
}

var corsOptions = cors.Options{
	// AllowedOrigins: []string{"*"},
	AllowOriginFunc: func(r *http.Request, origin string) bool {
		return true
	},
	AllowedMethods: []string{
		http.MethodHead,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	},
	AllowedHeaders:   []string{"*"},
	AllowCredentials: true,
}
