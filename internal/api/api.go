package api

import (
	"github.com/go-chi/chi/v5"
)

// @title Moonbase
// @description ### Git-based headless CMS API
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
		// low level github apis
		r.Group(func(r chi.Router) {
			r.Get("/repos", getRepos)
			r.Get("/repos/{owner}/{repo}/branches", getBranches)
			r.Get("/repos/{owner}/{repo}/tree/{ref}", getTree)
			r.Get("/repos/{owner}/{repo}/tree/{ref}/*", getTree)
			r.Get("/repos/{owner}/{repo}/blob/{ref}/*", getBlob)
			r.Post("/repos/{owner}/{repo}/blob/{ref}/*", postBlob)
			r.Delete("/repos/{owner}/{repo}/blob/{ref}/*", delBlob)
		})
		// higher level cms apis
		r.Group(func(r chi.Router) {
			// dashboard
			r.Get("/cms/{owner}/{repo}/{ref}", getInfo)
			// collections
			r.Get("/cms/{owner}/{repo}/{ref}/collections", getCollections)
			r.Post("/cms/{owner}/{repo}/{ref}/collections", postCollection)
			r.Delete("/cms/{owner}/{repo}/{ref}/collections/{collection}", delCollection)
			// entries
			r.Get("/cms/{owner}/{repo}/{ref}/collections/{collection}", getEntries)
			r.Get("/cms/{owner}/{repo}/{ref}/collections/{collection}/{entry}", getEntry)
			r.Post("/cms/{owner}/{repo}/{ref}/collections/{collection}", postEntry)
			r.Put("/cms/{owner}/{repo}/{ref}/collections/{collection}/{entry}", putEntry)
			r.Delete("/cms/{owner}/{repo}/{ref}/collections/{collection}/{entry}", delEntry)

			r.Post("/cms/{owner}/{repo}/{ref}/images", postImage)

			r.Get("/cms/{owner}/{repo}/{ref}/settings", getSettings)
			r.Get("/cms/{owner}/{repo}/{ref}/settings/{setting}", getSetting)
		})
	})

	return r
}
