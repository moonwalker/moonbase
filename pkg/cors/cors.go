package cors

import (
	"net/http"

	corsHandler "github.com/rs/cors"
)

func AllowAll() *corsHandler.Cors {
	return corsHandler.New(corsHandler.Options{
		// AllowedOrigins: []string{"*"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
}
