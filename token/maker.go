package token

import (
	"time"
)

type Maker interface {
	CreateToken(email string, id int, duration time.Duration) (string, int, error)
	VerifyToken(tokenString string) (*Payload, error)
}
