package api

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/moonwalker/moonbase/internal/runtime"
)

const (
	jsonContentType = "application/json"
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
	response(w, r, "index.html", &indexData{runtime.Name, runtime.ShortRev()})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	response(w, r, "404.html", &Error{Code: http.StatusNotFound, Message: "Not Found"})
}

func response(w http.ResponseWriter, r *http.Request, view string, data any) {
	switch r.Header.Get("Content-type") {
	case jsonContentType:
		jsonEncode(w, data)
	default:
		t.ExecuteTemplate(w, view, data)
	}
}
