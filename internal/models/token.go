package models

import (
	"time"
)

const (
	ScopeAuthentication = "authentication"
)

// Token is the type of authentication token
type Token struct {
	PlainText string    `json:"token"`
	UserId    int64     `json:"-"`
	Hash      []byte    `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}
