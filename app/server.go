package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/moonwalker/moonbase/app/api"
	"github.com/moonwalker/moonbase/app/pages"
	"github.com/moonwalker/moonbase/pkg/env"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

type Options struct {
	Port int
}

type Server struct {
	*Options
}

type TokenData struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	AccessToken string `json:"accessToken"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

var (
	githubConfig = &oauth2.Config{
		Scopes:   []string{"user:email", "read:org"},
		Endpoint: githuboauth.Endpoint,
	}
	oauthStateString = xid.New().String()
)

func NewServer(options *Options) *Server {
	return &Server{options}
}

func (s *Server) Listen() error {
	mux := http.NewServeMux()

	githubConfig.ClientID = env.GithubClientID
	githubConfig.ClientSecret = env.GithubClientSecret

	// Debug
	mux.HandleFunc("/debug", api.Debug)

	// Home
	pages.Handler(mux)

	// Login route
	mux.HandleFunc("/login/github", githubAuth)

	// Github callback
	mux.HandleFunc("/login/github/callback", githubCallback)

	addr := fmt.Sprintf(":%d", s.Options.Port)
	log.Printf("HTTP Server running at port %s", addr)

	return http.ListenAndServe(addr, mux)
}
