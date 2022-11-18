package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/go-github/v48/github"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/moonwalker/moonbase/docs"
	"github.com/moonwalker/moonbase/pkg/env"
	"github.com/moonwalker/moonbase/pkg/jwt"
)

// @title Moonbase API
// @version 1.0
// @description Git-based headless CMS API.

// @license.name MIT
// @license.url https://github.com/moonwalker/moonbase/blob/main/LICENSE

// @host moonbase.mw.zone
// @BasePath /api
func Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/debug", Debug)

	setAuthConfig() // TODO: find a better way
	// Login route
	r.HandleFunc("/login/github", githubAuth)

	// Github callback
	r.HandleFunc("/login/github/callback", githubCallback)

	// api routes which needs authenticated gh user
	r.Group(func(r chi.Router) {
		r.Use(withUser)
		// ...
	})

	// swagger docs
	r.Get("/docs/*", httpSwagger.Handler())

	return r
}

func withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// get token from authorization header
		bearer := r.Header.Get("Authorization")
		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			tokenString = bearer[7:]
		}
		if len(tokenString) == 0 {
			httpError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), fmt.Errorf("no auth token"))
			return
		}

		ghUser := &github.User{}
		token, err := jwt.VerifyAndDecrypt(env.JweKey, env.JwtKey, tokenString, ghUser)
		if err != nil {
			httpError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), err)
			return
		}

		if !token.Valid {
			httpError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), fmt.Errorf("invalid auth token"))
			return
		}

		// add auth claims to context
		ctx := context.WithValue(r.Context(), jwt.AUTH_CLAIMS_KEY, token.Claims)

		// authenticated, pass it through
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
