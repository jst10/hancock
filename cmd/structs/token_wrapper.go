package structs

import (
	"time"
)

type TokenWrapper struct {
	User       *User
	Token      string
	Expiration *time.Time
}
