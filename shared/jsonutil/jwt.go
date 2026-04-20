package jsonutil

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	UserID      string   `json:"sub"`
	Login       string   `json:"login"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions,omitempty"`

	jwt.RegisteredClaims
}

func ParseAccessToken(tokenStr string, secret []byte) (*AccessClaims, error) {
	claims := &AccessClaims{}

	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	_, err := parser.ParseWithClaims(
		tokenStr,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	return claims, nil
}
func ParseAccessTokenWithValidation(tokenStr string, secret []byte) (*AccessClaims, error) {
	claims := &AccessClaims{}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token invalid")
	}

	return claims, nil
}
