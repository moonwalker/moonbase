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
				if (req.headers['Authorization']) {
					req.headers['Authorization'] = 'Bearer ' + req.headers['Authorization']
				}
				return req
			}`,
			"onComplete": `() => {
				const btn = document.createElement("button")
				btn.innerText = "Login"
				btn.className = "btn"
				btn.style.marginRight = "1em"
				btn.onclick = () => {
					window.open("/login/github?returnURL=/login/github/authenticate")
				}
				const authWrapper = document.querySelector(".auth-wrapper")
				// authWrapper.style.cssText += 'justify-content: space-between'
				authWrapper.insertBefore(btn, authWrapper.firstChild)
			}`,
		}),
	)
}
