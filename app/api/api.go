package api

import (
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Moonbase API
// @description Git-based headless CMS API.
// @version 1.0

// @license.name MIT
// @license.url https://github.com/moonwalker/moonbase/blob/main/LICENSE
func Routes() chi.Router {
	r := chi.NewRouter()

	// index, 404, etc.
	r.Mount("/", core())

	// swagger docs
	r.Get("/docs", docsRedirect)
	r.Get("/docs/*", httpSwagger.Handler())

	// github login
	r.Get("/login/github", githubAuth)

	// github login callback
	r.Get("/login/github/callback", githubCallback)

	// github login authenticate
	r.Get("/login/github/authenticate", authenticateHandler)
	r.Get("/login/github/authenticate/{code}", authenticateHandler)

	// api routes which needs authenticated user token
	r.Group(func(r chi.Router) {
		r.Use(withUser)
		r.Get("/list", getRepositories)
		r.Get("/{owner}/{repo}/branches", getBranches)
	})

	return r
}
