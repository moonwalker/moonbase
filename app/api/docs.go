package api

import (
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/moonwalker/moonbase/docs"
)

func docsRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/docs" || r.URL.Path == "/docs/" {
		http.Redirect(w, r, "/docs/index.html", http.StatusTemporaryRedirect)
	}
}

const (
	loginBtnHtml = `<button id="loginbtn" class="btn" style="margin-right: 1em">Login</button>`
)

func docsHandler() http.HandlerFunc {
	return httpSwagger.Handler(
		httpSwagger.UIConfig(map[string]string{
			"requestInterceptor": `(req) => {
				if (req.headers['Authorization']) {
					req.headers['Authorization'] = 'Bearer ' + req.headers['Authorization']
				}
				return req
			}`,
			"onComplete": fmt.Sprintf(`() => {
				document.querySelector(".auth-wrapper").insertAdjacentHTML("afterbegin", '%s')
				document.querySelector("#loginbtn").addEventListener('click', () => {
					window.open("/login/github?returnURL=/login/github/authenticate")
				})
			}`, loginBtnHtml),
		}),
	)
}
