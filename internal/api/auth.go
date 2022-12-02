package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/xid"

	"github.com/moonwalker/moonbase/internal/env"
	"github.com/moonwalker/moonbase/internal/gh"
	"github.com/moonwalker/moonbase/internal/jwt"
)

// test the flow:
// http://localhost:8080/login/github?return_url=/login/github/authenticate

type ctxKey int

const (
	retUrlCodePath            = 0
	retUrlCodeQuery           = 1
	oauthStateSep             = "|"
	codeTokenExpires          = time.Minute
	accessTokenExpires        = time.Hour * 24
	ctxKeyAccessToken  ctxKey = iota
)

var (
	oauthStateSecret = xid.New().String()
)

type User struct {
	Login *string `json:"login"`
	Email *string `json:"email"`
	Image *string `json:"image"`
	Token string  `json:"token"`
}

func githubAuth(w http.ResponseWriter, r *http.Request) {
	state, err := encodeState(r)
	if err != nil {
		errAuthEncState().Log(r, err).Json(w)
		return
	}

	url := gh.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func githubCallback(w http.ResponseWriter, r *http.Request) {
	secret, returnURL := decodeState(r)
	if secret != oauthStateSecret {
		err := fmt.Errorf("expected: %s actual: %s", oauthStateSecret, secret)
		errAuthBadSecret().Log(r, err).Json(w)
		return
	}

	code := r.FormValue("code")
	url, err := returnURLWithCode(returnURL, code, retUrlCodePath)
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

	decoded, err := decryptExchangeCode(code)
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

	et, err := encryptAccessToken(token)
	if err != nil {
		errAuthEncToken().Log(r, err).Json(w)
		return
	}

	usr := &User{
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

func encryptAccessToken(accessToken string) (string, error) {
	te, err := jwt.EncryptAndSign(env.JweKey, env.JwtKey, []byte(accessToken), accessTokenExpires)
	if err != nil {
		return "", err
	}

	return te, nil
}

func encodeState(r *http.Request) (string, error) {
	returnURL := r.FormValue("return_url")

	u, err := url.Parse(returnURL)
	if err != nil {
		return "", err
	}

	if !u.IsAbs() {
		u, err = url.Parse(r.Referer())
		u.Path = returnURL
	}

	state := fmt.Sprintf("%s%s%s", oauthStateSecret, oauthStateSep, u)
	return base64.URLEncoding.EncodeToString([]byte(state)), nil
}

func decodeState(r *http.Request) (string, string) {
	state, _ := base64.URLEncoding.DecodeString(r.FormValue("state"))
	parts := strings.Split(string(state), oauthStateSep)
	return parts[0], parts[1]
}

func returnURLWithCode(returnURL, code string, m int) (string, error) {
	u, err := url.Parse(returnURL)
	if err != nil {
		return "", err
	}

	codeToken, err := jwt.EncryptAndSign(env.JweKey, env.JwtKey, []byte(code), codeTokenExpires)
	if err != nil {
		return "", err
	}

	switch {
	case m == retUrlCodeQuery:
		u.RawQuery = url.Values{
			"code": {codeToken},
		}.Encode()
	case m == retUrlCodePath:
		u.Path = path.Join(u.Path, codeToken)
	}

	return u.String(), nil
}

func decryptExchangeCode(code string) (string, error) {
	token, err := jwt.VerifyAndDecrypt(env.JweKey, env.JwtKey, code)
	if err != nil {
		return "", err
	}

	return string(token.Claims.(*jwt.AuthClaims).Data), nil
}

func withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// get token from authorization header
		bearer := r.Header.Get("Authorization")
		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			tokenString = bearer[7:]
		}
		if len(tokenString) == 0 {
			errAuthNoToken().Json(w)
			return
		}

		token, err := jwt.VerifyAndDecrypt(env.JweKey, env.JwtKey, tokenString)
		if err != nil {
			errAuthBadToken().Json(w)
			return
		}

		authClaims, ok := token.Claims.(*jwt.AuthClaims)
		if !ok {
			errAuthBadClaims().Json(w)
			return
		}

		// add auth claims to context
		ctx := context.WithValue(r.Context(), ctxKeyAccessToken, string(authClaims.Data))

		// authenticated, pass it through
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func accessTokenFromContext(ctx context.Context) string {
	return ctx.Value(ctxKeyAccessToken).(string)
}
