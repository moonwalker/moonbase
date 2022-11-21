package api

type errorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	errNoAuthToken        = &errorData{1001, "auth token not provided"}
	errUnauthorized       = &errorData{1002, "unauthorized"}
	errInvalidAuthClaims  = &errorData{1003, "invalid auth claims type"}
	errFailEncOAuthState  = &errorData{1004, "failed to encode oauth state"}
	errInvalidOAuthSecret = &errorData{1005, "invalid oauth state secret"}
	errFailedEncRetURL    = &errorData{1006, "failed to encrypt return url with oauth exchange code"}
	errAuthCodeMissing    = &errorData{1007, "auth code missing"}
	errFailDecOAuthCode   = &errorData{1008, "failed to decrypt oauth exchange code"}
	errFailOAuthExchange  = &errorData{1009, "oauth exchange failed"}
	errClientFailGetUser  = &errorData{1010, "github client failed to get user"}
	errFailEncAccessToken = &errorData{1011, "failed to encrypt token"}
)
