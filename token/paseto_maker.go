package token

import (
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
	"time"
)

type PasetoMaker struct {
	paseto        *paseto.V2
	symmetricKey  []byte
	asymmetricKey []byte
	footer        string
}

func NewPasetoMaker(symmetricKey string, asymmetricKey string, footer string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:        paseto.NewV2(),
		symmetricKey:  []byte(symmetricKey),
		asymmetricKey: []byte(asymmetricKey),
		footer:        footer,
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(uid uint64, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(uid, duration)
	if err != nil {
		return "", payload, err
	}
	tokenStr, err := maker.paseto.Encrypt(maker.symmetricKey, payload, maker.footer)
	return tokenStr, payload, err
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	var jsonToken Payload
	var newFooter string
	err := maker.paseto.Decrypt(token, maker.symmetricKey, &jsonToken, &newFooter)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = jsonToken.Valid()
	if err != nil {
		return nil, ErrExpiredToken
	}
	return &jsonToken, nil
}
