package api

import (
	"github.com/gin-gonic/gin"
)

// @title Moonbase API
// @description Git-based headless CMS API.

// @license.name MIT
// @license.url https://github.com/moonwalker/moonbase/blob/main/LICENSE
func Router() *gin.Engine {
	r := gin.Default()

	r.SetTrustedProxies(nil)
	r.SetHTMLTemplate(t)

	// 404
	r.NoRoute(notFound)

	// index
	staticAssets(r)
	r.GET("/", index)

	// swagger docs
	r.GET("/docs/*any", docs())

	// github login
	r.GET("/login/github", githubAuth)

	// github login callback
	r.GET("/login/github/callback", githubCallback)

	// api routes with authorized user token
	authorized := r.Group("/")
	authorized.Use(withUser())
	{
		// ...
	}

	return r
}
