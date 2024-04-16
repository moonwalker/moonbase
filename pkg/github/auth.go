package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/moonwalker/moonbase/internal/env"
	"github.com/moonwalker/moonbase/internal/jwt"
)

type ctxKey int

const (
	ctxKeyAccessToken ctxKey = iota
	ctxKeyUser        ctxKey = iota

	RetUrlCodePath     = 0
	RetUrlCodeQuery    = 1
	oauthStateSep      = "|"
	codeTokenExpires   = time.Minute
	accessTokenExpires = time.Hour * 24
)

type User struct {
	Login *string `json:"login"`
	Email *string `json:"email"`
	Image *string `json:"image"`
	Token string  `json:"token"`
}

func WithUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// get token from authorization header
		bearer := r.Header.Get("Authorization")
		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			tokenString = bearer[7:]
		}
		if len(tokenString) == 0 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		token, err := jwt.VerifyAndDecrypt(env.JweKey, env.JwtKey, tokenString)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		authClaims, ok := token.Claims.(*jwt.AuthClaims)
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ghUser, err := GetUser(r.Context(), string(authClaims.Data))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// add auth claims to context
		ctx := context.WithValue(context.WithValue(r.Context(), ctxKeyAccessToken, string(authClaims.Data)), ctxKeyUser, *ghUser.Login)
		// authenticated, pass it through
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AccessTokenFromContext(ctx context.Context) string {
	return ctx.Value(ctxKeyAccessToken).(string)
}

func UserFromContext(ctx context.Context) string {
	return ctx.Value(ctxKeyUser).(string)
}

func EncryptAccessToken(accessToken string) (string, error) {
	te, err := jwt.EncryptAndSign(env.JweKey, env.JwtKey, []byte(accessToken), accessTokenExpires)
	if err != nil {
		return "", err
	}

	return te, nil
}

func EncodeState(r *http.Request, oauthStateSecret string) (string, error) {
	returnURL := r.FormValue("return_url")

	u, err := url.Parse(returnURL)
	if err != nil {
		return "", err
	}

	if !u.IsAbs() {
		u, err = url.Parse(r.Referer())
		if err != nil {
			return "", err
		}
		u.Path = returnURL
	}

	state := fmt.Sprintf("%s%s%s", oauthStateSecret, oauthStateSep, u)
	return base64.URLEncoding.EncodeToString([]byte(state)), nil
}

func DecodeState(r *http.Request) (string, string) {
	state, _ := base64.URLEncoding.DecodeString(r.FormValue("state"))
	parts := strings.Split(string(state), oauthStateSep)
	return parts[0], parts[1]
}

func ReturnURLWithCode(returnURL, code string, m int) (string, error) {
	u, err := url.Parse(returnURL)
	if err != nil {
		return "", err
	}

	codeToken, err := jwt.EncryptAndSign(env.JweKey, env.JwtKey, []byte(code), codeTokenExpires)
	if err != nil {
		return "", err
	}

	switch {
	case m == RetUrlCodeQuery:
		u.RawQuery = url.Values{
			"code": {codeToken},
		}.Encode()
	case m == RetUrlCodePath:
		u.Path = path.Join(u.Path, codeToken)
	}

	return u.String(), nil
}

func DecryptExchangeCode(code string) (string, error) {
	token, err := jwt.VerifyAndDecrypt(env.JweKey, env.JwtKey, code)
	if err != nil {
		return "", err
	}

	return string(token.Claims.(*jwt.AuthClaims).Data), nil
}
