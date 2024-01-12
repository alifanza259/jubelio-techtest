package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secret string
}

type Payload struct {
	Email     string    `json:"email"`
	ID        int       `json:"id"`
	Issuer    string    `json:"iss"`
	ExpiresAt time.Time `json:"exp"`
	IssuedAt  time.Time `json:"iat"`
}

type CustomClaimsGolangJwt struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
	jwt.RegisteredClaims
}

func NewJWTMaker(secret string) Maker {
	return &JWTMaker{
		secret: secret,
	}
}

func (maker *JWTMaker) CreateToken(email string, id int, duration time.Duration) (string, int, error) {
	expiresAt := time.Now().Add(duration)
	claims := &CustomClaimsGolangJwt{
		email,
		id,
		jwt.RegisteredClaims{
			Issuer:    "login",
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := t.SignedString([]byte(maker.secret))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create token: %w", err)
	}

	return tokenString, int(expiresAt.Unix()), nil

}

func (maker *JWTMaker) VerifyToken(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaimsGolangJwt{}, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("token is invalid")
		}

		return []byte(maker.secret), nil
	})

	if err != nil {
		return nil, errors.New("token is invalid")
	}

	claims, ok := token.Claims.(*CustomClaimsGolangJwt)
	if !ok {
		return nil, errors.New("token is invalid")
	}

	payload := &Payload{
		Email:     claims.Email,
		ID:        claims.ID,
		Issuer:    claims.Issuer,
		ExpiresAt: claims.ExpiresAt.Time,
		IssuedAt:  claims.IssuedAt.Time,
	}
	return payload, nil
}
