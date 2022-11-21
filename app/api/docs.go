package api

import (
	"net/http"

	_ "github.com/moonwalker/moonbase/docs"
)

func docsRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/docs" || r.URL.Path == "/docs/" {
		http.Redirect(w, r, "/docs/index.html", http.StatusTemporaryRedirect)
	}
}
