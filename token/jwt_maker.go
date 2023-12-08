package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const minSecreKeySize = 32

type JwtMaker struct {
	secreKey string
}

func (j JwtMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	jwtToken, err := token.SignedString([]byte(j.secreKey))
	return jwtToken, payload, err
}

func (j JwtMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.secreKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {

		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}

func NewJwtMaker(secreKey string) (Maker, error) {
	if len(secreKey) < minSecreKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least: %d characters", minSecreKeySize)
	}

	return &JwtMaker{secreKey: secreKey}, nil
}
