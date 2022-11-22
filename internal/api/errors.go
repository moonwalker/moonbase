package api

import (
	"net/http"
	"strconv"

	"github.com/rs/xid"

	"github.com/moonwalker/moonbase/internal/log"
)

var (
	errMsg = map[int]string{
		1000: "unauthorized",
		1001: "auth token not provided",
		1002: "invalid auth claims type",
		1003: "failed to encode oauth state",
		1004: "invalid oauth state secret",
		1005: "failed to encrypt return url with oauth exchange code",
		1006: "auth code missing",
		1007: "failed to decrypt oauth exchange code",
		1008: "oauth exchange failed",
		1009: "github client failed to get user",
		1010: "failed to encrypt token",
	}
	errUnauthorized       = func() *errorData { return makeError(401, 1000) }
	errNoAuthToken        = func() *errorData { return makeError(401, 1001) }
	errInvalidAuthClaims  = func() *errorData { return makeError(500, 1002) }
	errFailEncOAuthState  = func() *errorData { return makeError(500, 1003) }
	errInvalidOAuthSecret = func() *errorData { return makeError(500, 1004) }
	errFailedEncRetURL    = func() *errorData { return makeError(500, 1005) }
	errAuthCodeMissing    = func() *errorData { return makeError(500, 1006) }
	errFailDecOAuthCode   = func() *errorData { return makeError(500, 1007) }
	errFailOAuthExchange  = func() *errorData { return makeError(500, 1008) }
	errClientFailGetUser  = func() *errorData { return makeError(500, 1009) }
	errFailEncAccessToken = func() *errorData { return makeError(500, 1010) }
)

func makeError(statusCode, errorCode int) *errorData {
	return &errorData{xid.New().String(), strconv.Itoa(statusCode), strconv.Itoa(errorCode), errMsg[errorCode]}
}

type errorData struct {
	ID      string `json:"id,omitempty"`
	Status  string `json:"status"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

func (e *errorData) Json(w http.ResponseWriter) *errorData {
	status, _ := strconv.Atoi(e.Status)
	jsonResponse(w, status, e)
	return e
}

func (e *errorData) Log(err error) *errorData {
	log.Error(err).Str("id", e.ID).Msg(e.Message)
	return e
}
