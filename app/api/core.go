package api

import (
	"embed"
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi"
)

const (
	jsonContentType = "application/json"
)

type indexJson struct {
	Server  string `json:"server"`
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
	response(w, r, "index.html", &indexJson{"Moonbase", "v1.0"})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	response(w, r, "404.html", &Error{Code: http.StatusNotFound, Message: "Not Found"})
}

func response(w http.ResponseWriter, r *http.Request, view string, data any) {
	switch r.Header.Get("Content-type") {
	case jsonContentType:
		w.Header().Set("Content-Type", jsonContentType)
		json.NewEncoder(w).Encode(data)
	default:
		t.ExecuteTemplate(w, view, data)
	}
}
