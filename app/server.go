package app

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/moonwalker/moonbase/app/api"
	"github.com/moonwalker/moonbase/app/pages"
)

type Options struct {
	Port int
}

type Server struct {
	*Options
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewServer(options *Options) *Server {
	return &Server{options}
}

func (s *Server) Listen() error {
	r := chi.NewRouter()

	r.Mount("/", pages.Handler())
	r.Mount("/api", api.Routes())

	addr := fmt.Sprintf(":%d", s.Options.Port)
	log.Printf("HTTP Server running at port %s", addr)

	return http.ListenAndServe(addr, r)
}
