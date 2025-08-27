package session

import (
	"encoding/json"
	"gpt4omini/config"
	"gpt4omini/types"
	"time"
)

var cfg = config.GetConfig()

// BaseSession contains fields shared across all sessions.
type BaseSession struct {
	ID           string
	Type         string
	clientSecret types.ClientSecret
	createdAt    time.Time
}

func (bs *BaseSession) GetID() string {
	return bs.ID
}

func (bs *BaseSession) GetType() string {
	return bs.Type
}

func (bs *BaseSession) GetClientSecretExpirationTime() time.Time {
	return time.Unix(bs.clientSecret.ExpiresAt, 0)
}

func (bs *BaseSession) HasClientSecretExpired() bool {
	return time.Now().After(bs.GetClientSecretExpirationTime())
}

func (bs *BaseSession) GetClientSecretValue() string {
	return bs.clientSecret.Value
}

func (bs *BaseSession) GetCreationTime() types.ClientSecret {
	return bs.clientSecret
}

func (bs *BaseSession) String() string {
	out, _ := json.MarshalIndent(struct {
		ID           string    `json:"id"`
		ClientSecret string    `json:"client_secret"`
		CreatedAt    time.Time `json:"created_at"`
		Type         string    `json:"type"`
		ExpiresAt    time.Time `json:"expires_at"`
	}{
		ID:           bs.GetID(),
		ClientSecret: bs.GetClientSecretValue(),
		CreatedAt:    bs.createdAt,
		Type:         bs.GetType(),
		ExpiresAt:    bs.GetClientSecretExpirationTime(),
	}, "", "  ")
	return string(out)
}

// Session is the common interface for all session types.
type Session interface {
	Start()
	Close()
	GetID() string
	GetType() string
	HasExpired() bool
	GetClientSecretValue() string
	GetClientSecretExpirationTime() time.Time
}
