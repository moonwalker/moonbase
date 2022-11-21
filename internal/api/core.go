package api

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi"

	"github.com/moonwalker/moonbase/internal/runtime"
)

type indexData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

//go:embed static
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "static/*.html"))

func core() chi.Router {
	r := chi.NewRouter()

	fsys := fs.FS(resources)
	assets, _ := fs.Sub(fsys, "static/assets")
	fs := http.FileServer(http.FS(assets))

	r.HandleFunc("/", index)
	r.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	r.NotFound(notFound)

	return r
}

func index(w http.ResponseWriter, r *http.Request) {
	response(w, r, http.StatusOK, "index.html", &indexData{runtime.Name, runtime.ShortRev()})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	response(w, r, http.StatusNotFound, "404.html", &errorData{Code: http.StatusNotFound, Message: http.StatusText(http.StatusNotFound)})
}

func response(w http.ResponseWriter, r *http.Request, statusCode int, view string, data any) {
	ct := r.Header.Get("Accept")
	if strings.Contains(ct, "application/json") {
		jsonResponse(w, statusCode, data)
	} else {
		t.ExecuteTemplate(w, view, data)
	}
}
