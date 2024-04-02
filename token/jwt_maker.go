package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const minSecretKeyLen = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeyLen {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeyLen)
	}
	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}

func (m *JWTMaker) CreateToken(uid uint64, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(uid, duration)
	if err != nil {
		return "", payload, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenStr, err := jwtToken.SignedString([]byte(m.secretKey))
	return tokenStr, payload, err
}

func (m *JWTMaker) VerifyToken(token string) (payload *Payload, err error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(m.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, err
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil

}
