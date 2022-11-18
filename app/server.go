package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/moonwalker/moonbase/app/api"
	"github.com/moonwalker/moonbase/app/pages"
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
	mux := http.NewServeMux()

	mux.HandleFunc("/", pages.Index)
	mux.HandleFunc("/debug", api.Debug)
	// ...

	addr := fmt.Sprintf(":%d", s.Options.Port)
	log.Printf("HTTP Server running at port %s", addr)

	return http.ListenAndServe(addr, mux)
}
