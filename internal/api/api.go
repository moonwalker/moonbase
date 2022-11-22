package api

import (
	"github.com/go-chi/chi"
)

// @title Moonbase API
// @description ### Git-based headless CMS API.
// @version 1.0

// @license.name MIT
// @license.url https://github.com/moonwalker/moonbase/blob/main/LICENSE

// @securityDefinitions.apikey bearerToken
// @in header
// @name Authorization
func Routes() chi.Router {
	r := chi.NewRouter()

	// index, 404, etc.
	r.Mount("/", core())

	// swagger docs
	r.Get("/docs", docsRedirect)
	r.Get("/docs/*", docsHandler())

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
		r.Get("/user/repos", getRepositories)
		r.Get("/repos/{owner}/{repo}/branches", getBranches)
		r.Get("/repos/{owner}/{repo}/tree/{branch}", getTree)
		r.Get("/repos/{owner}/{repo}/tree/{branch}/*", getTree)
		r.Get("/repos/{owner}/{repo}/blob/{branch}/*", getBlob)
	})

	return r
}
