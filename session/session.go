package session

import (
	"context"
	"encoding/json"
	"gpt4omini/config"
	"gpt4omini/types"
	"sync"
	"time"
)

var cfg = config.GetConfig()

// BaseSession contains fields shared across all sessions.
type BaseSession struct {
	ID           string
	Type         string
	clientSecret types.ClientSecret
	createdAt    time.Time
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	mu           sync.Mutex
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
		Type         string    `json:"type"`
		ClientSecret string    `json:"client_secret"`
		CreatedAt    time.Time `json:"created_at"`
		ExpiresAt    time.Time `json:"expires_at"`
	}{
		ID:           bs.GetID(),
		Type:         bs.GetType(),
		ClientSecret: bs.GetClientSecretValue(),
		CreatedAt:    bs.createdAt,
		ExpiresAt:    bs.GetClientSecretExpirationTime(),
	}, "", "  ")
	return string(out)
}

// Session is the common interface for all session types.
type Session interface {
	// Start is the starting point for the new session.
	Start()
	// Close closes all connections (destroying session).
	close()
	// GetID provides the session id.
	GetID() string
	// GetType provides the session type.
	GetType() string
	// HasClientSecretExpired if the secret has expired, we cant keep communicating.
	HasClientSecretExpired() bool
	// GetClientSecretValue provides the secret value
	GetClientSecretValue() string
	// GetClientSecretExpirationTime provides the creation time of the session.
	GetClientSecretExpirationTime() time.Time
}
