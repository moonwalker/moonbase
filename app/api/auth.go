package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v48/github"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"

	"github.com/moonwalker/moonbase/pkg/env"
	"github.com/moonwalker/moonbase/pkg/jwt"
)

const (
	oauthStateSeparator = "|"
)

var (
	ghScopes         = []string{"user:email", "read:org"}
	oauthStateSecret = xid.New().String()
)

type TokenData struct {
	Login       string `json:"login"`
	AccessToken string `json:"accessToken"`
}

type User struct {
	Login *string `json:"login,omitempty"`
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Image *string `json:"image,omitempty"`
	Token string  `json:"token,omitempty"`
}

func githubConfig() *oauth2.Config {
	return &oauth2.Config{
		Scopes:       ghScopes,
		Endpoint:     githuboauth.Endpoint,
		ClientID:     env.GithubClientID,
		ClientSecret: env.GithubClientSecret,
	}
}

// Github login
//
//	@Summary		login via github
//	@Description	github login
//	@ID				gh-login
//	@Accept			json
//	@Produce		json
//	@Success		200		{string}	string			"ok"
//	@Router			/login/github [get]
func githubAuth(c *gin.Context) {
	state := encodeState(c.Request)
	url := githubConfig().AuthCodeURL(state, oauth2.AccessTypeOnline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func githubCallback(c *gin.Context) {
	githubConfig := githubConfig()

	secret, returnURL := decodeState(c.Request)
	if secret != oauthStateSecret {
		httpError(c, -1, "invalid oauth state secret", fmt.Errorf("expected: %s, actual: %s", oauthStateSecret, secret))
		return
	}

	code := c.Request.FormValue("code")

	// IDEA: here we can return with the 'code' (encrypted) to the client
	// and the client can initiate the oauth exchange step
	// so we won't need popup or iframe solutions
	println(returnURLWithCode(returnURL, code))

	token, err := githubConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		httpError(c, -1, "oauth exchange failed", err)
		return
	}

	oauthClient := githubConfig.Client(oauth2.NoContext, token)
	ghClient := github.NewClient(oauthClient)
	ghUser, _, err := ghClient.Users.Get(context.Background(), "")
	if err != nil {
		httpError(c, -1, "github client failed to get user", err)
		return
	}

	et, err := encryptAccessToken(ghUser, token.AccessToken)
	if err != nil {
		httpError(c, -1, "failed to encrypt token", err)
		return
	}

	u := &User{
		Login: ghUser.Login,
		Name:  ghUser.Name,
		Email: ghUser.Email,
		Image: ghUser.AvatarURL,
		Token: et,
	}

	c.JSON(http.StatusOK, u)
}

func encryptAccessToken(user *github.User, accessToken string) (string, error) {
	payload := &TokenData{Login: *user.Login, AccessToken: accessToken}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	te, err := jwt.EncryptAndSign(env.JweKey, env.JwtKey, data, 30)
	if err != nil {
		return "", err
	}

	return te, nil
}

func encodeState(r *http.Request) string {
	returnURL := r.FormValue("returnURL")
	if returnURL == "" {
		returnURL = r.Referer()
	}
	state := fmt.Sprintf("%s%s%s", oauthStateSecret, oauthStateSeparator, returnURL)
	return base64.URLEncoding.EncodeToString([]byte(state))
}

func decodeState(r *http.Request) (string, string) {
	state, _ := base64.URLEncoding.DecodeString(r.FormValue("state"))
	parts := strings.Split(string(state), oauthStateSeparator)
	return parts[0], parts[1]
}

func returnURLWithCode(returnURL, code string) (string, error) {
	codeJWT, err := jwt.EncryptAndSign(env.JweKey, env.JwtKey, []byte(code), 30)
	if err != nil {
		return "", err
	}

	u := fmt.Sprintf("%s?code=%s", returnURL, codeJWT)
	return u, nil
}

const USER_CTX_KEY = "gh-user"

func withUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// get token from authorization header
		bearer := c.Request.Header.Get("Authorization")
		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			tokenString = bearer[7:]
		}
		if len(tokenString) == 0 {
			httpError(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), fmt.Errorf("no auth token"))
			return
		}

		token, err := jwt.VerifyAndDecrypt(env.JweKey, env.JwtKey, tokenString)
		if err != nil {
			httpError(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), err)
			return
		}

		authClaims, ok := token.Claims.(*jwt.AuthClaims)
		if !ok {
			httpError(c, http.StatusInternalServerError, "invalid auth claims type", nil)
			return
		}

		ghUser := &github.User{}
		err = json.Unmarshal(authClaims.Data, ghUser)
		if err != nil {
			httpError(c, http.StatusInternalServerError, "failed to unmarshal auth claims data", err)
			return
		}

		// add auth claims to context
		ctx := context.WithValue(c.Request.Context(), USER_CTX_KEY, ghUser)
		c.Request = c.Request.WithContext(ctx)

		// authenticated, pass it through
		c.Next()
	}
}
