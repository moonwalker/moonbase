package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/moonwalker/moonbase/pkg/env"
	"github.com/moonwalker/moonbase/pkg/jwt"
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
		// ...
	})

	return r
}

const USER_CTX_KEY = "user-token"

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

		token, err := jwt.VerifyAndDecrypt(env.JweKey, env.JwtKey, tokenString)
		if err != nil {
			httpError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), err)
			return
		}

		authClaims, ok := token.Claims.(*jwt.AuthClaims)
		if !ok {
			httpError(w, http.StatusInternalServerError, "invalid auth claims type", nil)
			return
		}

		tokenData := &TokenData{}
		err = json.Unmarshal(authClaims.Data, tokenData)
		if err != nil {
			httpError(w, http.StatusInternalServerError, "failed to unmarshal auth claims data", err)
			return
		}

		// add auth claims to context
		ctx := context.WithValue(r.Context(), USER_CTX_KEY, tokenData)

		// authenticated, pass it through
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
