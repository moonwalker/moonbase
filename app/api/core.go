package api

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
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

var (
	t = template.Must(template.ParseFS(resources, "static/*.html"))
)

func staticAssets(r *gin.Engine) {
	fsys := fs.FS(resources)
	assets, _ := fs.Sub(fsys, "static/assets")
	r.StaticFS("/assets", http.FS(assets))
}

func index(c *gin.Context) {
	response(c, http.StatusOK, "index.html", &indexJson{"Moonbase", "v1.0"})
}

func notFound(c *gin.Context) {
	response(c, http.StatusNotFound, "404.html", &Error{Code: http.StatusNotFound, Message: "Not Found"})
}

func response(c *gin.Context, code int, view string, data any) {
	switch c.Request.Header.Get("Content-type") {
	case jsonContentType:
		c.Header("Content-Type", jsonContentType)
		c.JSON(code, data)
	default:
		c.HTML(code, view, data)
	}
}
