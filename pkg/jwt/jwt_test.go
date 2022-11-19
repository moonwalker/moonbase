package jwt

import (
	"encoding/json"
	"errors"
	"testing"
)

var (
	encKey = []byte("a4f9e6035517aae049edc0de0d815914")
	sigKey = []byte("c9cea3a1132598a1734bcaf03aa2ea98")
)

type testData struct {
	Name  string
	Email string
}

func TestEncryptDecrypt(t *testing.T) {
	payload := &testData{Name: "Foo", Email: "bar@example.com"}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Error(err)
	}

	te, err := EncryptAndSign(encKey, sigKey, data, 1)
	if err != nil {
		t.Error(err)
	}

	token, err := VerifyAndDecrypt(encKey, sigKey, te)
	if err != nil {
		t.Error(err)
	}

	authClaims, ok := token.Claims.(*AuthClaims)
	if !ok {
		t.Error(errors.New("invalid token claims type"))
	}

	outPayload := &testData{}
	err = json.Unmarshal(authClaims.Data, outPayload)
	if err != nil {
		t.Error(err)
	}

	if outPayload.Email != payload.Email {
		t.Fail()
	}
}
