package api

import (
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	d "github.com/moonwalker/moonbase/docs"
	"github.com/moonwalker/moonbase/internal/runtime"
)

const repoURL = "https://github.com/moonwalker/moonbase"

func docsRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/docs" || r.URL.Path == "/docs/" {
		http.Redirect(w, r, "/docs/index.html", http.StatusTemporaryRedirect)
	}
}

const (
	loginBtnHtml1 = `<button id="loginbtn" class="btn" style="margin-right: 1em; width: 10em">Login</button>`
	loginBtnHtml  = `<button id="loginbtn" class="btn" style="margin-right: 1em; width: 10em"><svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-brand-github" width="16" height="16" viewBox="0 0 20 20" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round"><path stroke="none" d="M0 0h24v24H0z" fill="none"></path><path d="M9 19c-4.3 1.4 -4.3 -2.5 -6 -3m12 5v-3.5c0 -1 .1 -1.4 -.5 -2c2.8 -.3 5.5 -1.4 5.5 -6a4.6 4.6 0 0 0 -1.3 -3.2a4.2 4.2 0 0 0 -.1 -3.2s-1.1 -.3 -3.5 1.3a12.3 12.3 0 0 0 -6.2 0c-2.4 -1.6 -3.5 -1.3 -3.5 -1.3a4.2 4.2 0 0 0 -.1 3.2a4.6 4.6 0 0 0 -1.3 3.2c0 4.6 2.7 5.7 5.5 6c-.6 .6 -.6 1.2 -.5 2v3.5"></path></svg> Login</button>`
)

const (
	jsPopupFunc = "function popupWindow(url, windowName, win, w, h) { const y = win.top.outerHeight / 2 + win.top.screenY - ( h / 2); const x = win.top.outerWidth / 2 + win.top.screenX - ( w / 2); return win.open(url, windowName, `toolbar=no, location=no, directories=no, status=no, menubar=no, scrollbars=no, resizable=no, copyhistory=no, width=${w}, height=${h}, top=${y}, left=${x}`);}"
)

func docsHandler() http.HandlerFunc {
	d.SwaggerInfo.Description += fmt.Sprintf("\n%s", revisionMarkdown())
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
					const authwin = popupWindow("/login/github?return_url=/login/github/authenticate", "swagger-gh-auth", window, 640, 480)
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

func revisionMarkdown() string {
	if runtime.IsDev() {
		return fmt.Sprintf("Revision: ```%s```", runtime.ShortRev())
	}
	return fmt.Sprintf("Revision: [```%s```](%s/commit/%[1]s)", runtime.ShortRev(), repoURL)
}
