package api

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/moonwalker/moonbase/docs"
)

func docsRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/docs" || r.URL.Path == "/docs/" {
		http.Redirect(w, r, "/docs/index.html", http.StatusTemporaryRedirect)
	}
}

func docsHandler() http.HandlerFunc {
	return httpSwagger.Handler(
		httpSwagger.UIConfig(map[string]string{
			"requestInterceptor": `(req) => {
				req.headers['Authorization'] = 'Bearer ' + req.headers['Authorization']
				return req
			}`,
		}),
	)
}
