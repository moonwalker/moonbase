package jwt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/moonwalker/moonbase/pkg/jwe"
)

const AUTH_CLAIMS_KEY = "auth-claims"

type authClaims struct {
	jwt.StandardClaims
	Data []byte `json:"data"`
}

func EncryptAndSign(encKey, sigKey []byte, payload interface{}, expiresInMin int) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encData, err := jwe.Encrypt(encKey, data)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &authClaims{
		Data: encData,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Duration(expiresInMin) * time.Minute).Unix(),
		},
	})

	return token.SignedString(sigKey)
}

func VerifyAndDecrypt(encKey, sigKey []byte, tokenString string, payload interface{}) (*jwt.Token, error) {
	claims := &authClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return sigKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims.Data, err = jwe.Decrypt(encKey, claims.Data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(claims.Data, &payload)
	if err != nil {
		return nil, err
	}

	return token, nil
}
