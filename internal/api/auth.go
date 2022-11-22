package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/go-github/v48/github"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"

	"github.com/moonwalker/moonbase/internal/env"
	"github.com/moonwalker/moonbase/internal/jwt"
	"github.com/moonwalker/moonbase/internal/log"
)

// testing the flow:
// http://localhost:8080/login/github?returnURL=/login/github/authenticate

const (
	userCtxKey              = "user-token"
	oauthStateSep           = "|"
	retUrlCodeQuery     int = 0
	retUrlCodePath          = 1
	codeTokenExpiresMin     = 1
	userTokenExpiresMin     = 30
)

var (
	ghScopes         = []string{"user:email", "read:org", "repo"}
	oauthStateSecret = xid.New().String()
)

type User struct {
	Login *string `json:"login,omitempty"`
	Image *string `json:"image,omitempty"`
	Token string  `json:"token"`
}

func githubConfig() *oauth2.Config {
	return &oauth2.Config{
		Scopes:       ghScopes,
		Endpoint:     githuboauth.Endpoint,
		ClientID:     env.GithubClientID,
		ClientSecret: env.GithubClientSecret,
	}
}

func githubAuth(w http.ResponseWriter, r *http.Request) {
	state, err := encodeState(r)
	if err != nil {
		log.Error(err)
		jsonResponse(w, http.StatusInternalServerError, errFailEncOAuthState)
		return
	}

	url := githubConfig().AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func githubCallback(w http.ResponseWriter, r *http.Request) {
	secret, returnURL := decodeState(r)
	if secret != oauthStateSecret {
		err := fmt.Errorf("expected: %s actual: %s", oauthStateSecret, secret)
		log.Error(err).Msg(errInvalidOAuthSecret.Message)
		jsonResponse(w, http.StatusInternalServerError, errInvalidOAuthSecret)
		return
	}

	code := r.FormValue("code")
	url, err := returnURLWithCode(returnURL, code, retUrlCodePath)
	if err != nil {
		log.Error(err).Msg(errFailedEncRetURL.Message)
		jsonResponse(w, http.StatusInternalServerError, errFailedEncRetURL)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	if code == "" {
		code = chi.URLParam(r, "code")
		if code == "" {
			jsonResponse(w, http.StatusInternalServerError, errAuthCodeMissing)
			return
		}
	}

	decoded, err := decryptExchangeCode(code)
	if err != nil {
		log.Error(err).Msg(errFailDecOAuthCode.Message)
		jsonResponse(w, http.StatusInternalServerError, errFailDecOAuthCode)
		return
	}

	githubConfig := githubConfig()
	token, err := githubConfig.Exchange(oauth2.NoContext, decoded)
	if err != nil {
		log.Error(err).Msg(errFailOAuthExchange.Message)
		jsonResponse(w, http.StatusInternalServerError, errFailOAuthExchange)
		return
	}

	oauthClient := githubConfig.Client(oauth2.NoContext, token)
	ghClient := github.NewClient(oauthClient)
	ghUser, _, err := ghClient.Users.Get(context.Background(), "")
	if err != nil {
		log.Error(err).Msg(errClientFailGetUser.Message)
		jsonResponse(w, http.StatusInternalServerError, errClientFailGetUser)
		return
	}

	et, err := encryptAccessToken(token.AccessToken)
	if err != nil {
		log.Error(err).Msg(errFailEncAccessToken.Message)
		jsonResponse(w, http.StatusInternalServerError, errFailEncAccessToken)
		return
	}

	usr := &User{
		Login: ghUser.Login,
		Image: ghUser.AvatarURL,
		Token: et,
	}

	jsonResponse(w, http.StatusOK, usr)
}

func encryptAccessToken(accessToken string) (string, error) {
	te, err := jwt.EncryptAndSign(env.JweKey, env.JwtKey, []byte(accessToken), userTokenExpiresMin)
	if err != nil {
		return "", err
	}

	return te, nil
}

func encodeState(r *http.Request) (string, error) {
	returnURL := r.FormValue("returnURL")

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

	codeToken, err := jwt.EncryptAndSign(env.JweKey, env.JwtKey, []byte(code), codeTokenExpiresMin)
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
			jsonResponse(w, http.StatusUnauthorized, errNoAuthToken)
			return
		}

		token, err := jwt.VerifyAndDecrypt(env.JweKey, env.JwtKey, tokenString)
		if err != nil {
			jsonResponse(w, http.StatusUnauthorized, errUnauthorized)
			return
		}

		authClaims, ok := token.Claims.(*jwt.AuthClaims)
		if !ok {
			jsonResponse(w, http.StatusInternalServerError, errInvalidAuthClaims)
			return
		}

		// add auth claims to context
		ctx := context.WithValue(r.Context(), userCtxKey, string(authClaims.Data))

		// authenticated, pass it through
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
