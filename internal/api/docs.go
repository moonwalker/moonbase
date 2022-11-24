package api

import (
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	d "github.com/moonwalker/moonbase/docs"
	"github.com/moonwalker/moonbase/internal/runtime"
)

func docsRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/docs" || r.URL.Path == "/docs/" {
		http.Redirect(w, r, "/docs/index.html", http.StatusTemporaryRedirect)
	}
}

const (
	loginBtnHtml = `<button id="loginbtn" class="btn" style="margin-right: 1em">Login</button>`
)

const (
	jsPopupFunc = "function popupWindow(url, windowName, win, w, h) { const y = win.top.outerHeight / 2 + win.top.screenY - ( h / 2); const x = win.top.outerWidth / 2 + win.top.screenX - ( w / 2); return win.open(url, windowName, `toolbar=no, location=no, directories=no, status=no, menubar=no, scrollbars=no, resizable=no, copyhistory=no, width=${w}, height=${h}, top=${y}, left=${x}`);}"
)

func docsHandler() http.HandlerFunc {
	d.SwaggerInfo.Description += fmt.Sprintf("\nRevision: ```%s```", runtime.ShortRev())
	return httpSwagger.Handler(
		httpSwagger.BeforeScript(jsPopupFunc),
		httpSwagger.UIConfig(map[string]string{
			"persistAuthorization": "true",
			"requestInterceptor": `(req) => {
				if (req.headers['Authorization']) {
					if (!req.headers['Authorization'].startsWith('Bearer')) {
						req.headers['Authorization'] = 'Bearer ' + req.headers['Authorization']
					}
				}
				return req
			}`,
			"onComplete": fmt.Sprintf(`() => {
				document.querySelector(".auth-wrapper").insertAdjacentHTML("afterbegin", '%s')
				document.querySelector("#loginbtn").addEventListener('click', () => {
					const authwin = popupWindow("/login/github?returnURL=/login/github/authenticate", "swagger-gh-auth", window, 640, 480)
					authwin.onload = () => {
						const res = JSON.parse(authwin.document.querySelector("pre").innerText)
						window.localStorage.setItem('authorized', JSON.stringify({"bearerToken":{"name":"bearerToken","schema":{"type":"apiKey","name":"Authorization","in":"header"},"value":res.token}}))
						authwin.close()
						window.location.reload()
					}
				})
			}`, loginBtnHtml),
		}),
	)
}
