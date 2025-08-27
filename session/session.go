package session

import (
	"gpt4omini/config"
	"time"
)

var cfg = config.GetConfig()

// Session is the common interface for all session types.
type Session interface {
	Start()
	Close()
	HasExpired() bool
	GetID() string
}

// BaseSession contains fields shared across all sessions.
type BaseSession struct {
	ID           string
	ClientSecret ClientSecret
	CreatedAt    time.Time
}

func (bs *BaseSession) GetID() string {
	return bs.ID
}

// HasExpired checks if the session secret expired.
func (bs *BaseSession) HasExpired() bool {
	expireTime := time.Unix(bs.ClientSecret.ExpiresAt, 0)
	return time.Now().After(expireTime)
}
