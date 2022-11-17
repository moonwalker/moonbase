package app

import (
	"fmt"
	"net/http"

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
	mux := http.NewServeMux()

	mux.HandleFunc("/debug", api.Debug)
	// ...

	addr := fmt.Sprintf(":%d", s.Options.Port)
	return http.ListenAndServe(addr, mux)
}
