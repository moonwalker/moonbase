package pages

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed public
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "public/*.html"))

func Handler(mux *http.ServeMux) {
	fsys := fs.FS(resources)
	assets, _ := fs.Sub(fsys, "public/assets")
	fs := http.FileServer(http.FS(assets))

	mux.HandleFunc("/", index)
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
}

func index(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "index.html", nil)
}
