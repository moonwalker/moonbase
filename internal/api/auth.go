package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/xid"

	gh "github.com/moonwalker/moonbase/pkg/github"
)

// test the flow:
// http://localhost:8080/login/github?return_url=/login/github/authenticate

var (
	oauthStateSecret = xid.New().String()
)

func githubAuth(w http.ResponseWriter, r *http.Request) {
	state, err := gh.EncodeState(r, oauthStateSecret)
	if err != nil {
		errAuthEncState().Log(r, err).Json(w)
		return
	}

	url := gh.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func githubCallback(w http.ResponseWriter, r *http.Request) {
	secret, returnURL := gh.DecodeState(r)
	if secret != oauthStateSecret {
		err := fmt.Errorf("expected: %s actual: %s", oauthStateSecret, secret)
		errAuthBadSecret().Log(r, err).Json(w)
		return
	}

	code := r.FormValue("code")
	url, err := gh.ReturnURLWithCode(returnURL, code, gh.RetUrlCodePath)
	if err != nil {
		errAuthEncRetURL().Log(r, err).Json(w)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	if code == "" {
		code = chi.URLParam(r, "code")
		if code == "" {
			errAuthCodeMissing().Log(r, nil).Json(w)
			return
		}
	}

	decoded, err := gh.DecryptExchangeCode(code)
	if err != nil {
		errAuthDecOAuth().Log(r, err).Json(w)
		return
	}

	token, err := gh.Exchange(decoded)
	if err != nil {
		errAuthExchange().Log(r, err).Json(w)
		return
	}

	ghUser, err := gh.GetUser(r.Context(), token)
	if err != nil {
		errAuthGetUser().Log(r, err).Json(w)
		return
	}

	et, err := gh.EncryptAccessToken(token)
	if err != nil {
		errAuthEncToken().Log(r, err).Json(w)
		return
	}

	usr := &gh.User{
		Login: ghUser.Login,
		Email: ghUser.Email,
		Image: ghUser.AvatarURL,
		Token: et,
	}

	if usr.Email == nil {
		e := fmt.Sprintf("%d+%s@users.noreply.github.com", *ghUser.ID, *ghUser.Login)
		usr.Email = &e
	}

	jsonResponse(w, http.StatusOK, usr)
}
