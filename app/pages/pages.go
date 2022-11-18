package pages

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi"
)

//go:embed public
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "public/*.html"))

func Handler() chi.Router {
	r := chi.NewRouter()

	fsys := fs.FS(resources)
	assets, _ := fs.Sub(fsys, "public/assets")
	fs := http.FileServer(http.FS(assets))

	r.HandleFunc("/", index)
	r.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	return r
}

func index(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "index.html", nil)
}
