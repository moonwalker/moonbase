package jwt

import (
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
	te, err := EncryptAndSign(encKey, sigKey, payload, 1)
	if err != nil {
		t.Error(err)
	}

	outPayload := &testData{}
	_, err = VerifyAndDecrypt(encKey, sigKey, te, outPayload)
	if err != nil {
		t.Error(err)
	}

	if outPayload.Email != payload.Email {
		t.Fail()
	}
}
