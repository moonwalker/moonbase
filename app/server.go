package app

import (
	"fmt"

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

func (o *Options) Addr() string {
	return fmt.Sprintf(":%d", o.Port)
}

func (s *Server) Listen() error {
	r := api.Router()
	return r.Run(s.Options.Addr())
}
