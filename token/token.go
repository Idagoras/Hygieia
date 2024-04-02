package token

import "time"

type Maker interface {
	CreateToken(uid uint64, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
