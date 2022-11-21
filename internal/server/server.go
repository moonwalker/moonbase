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

type Options struct {
	Port int
}

type Server struct {
	*Options
}

func NewServer(options *Options) *Server {
	return &Server{options}
}

func (s *Server) Listen() error {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(flate.DefaultCompression))
	r.Use(cors.Handler(corsOptions))

	r.Mount("/", api.Routes())

	addr := fmt.Sprintf(":%d", s.Options.Port)
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
