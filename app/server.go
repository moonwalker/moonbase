package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/moonwalker/moonbase/app/api"
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

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Mount("/", api.Routes())

	addr := fmt.Sprintf(":%d", s.Options.Port)
	log.Printf("HTTP Server running at port %s", addr)

	return http.ListenAndServe(addr, r)
}
