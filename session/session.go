package session

import (
	"time"
)

// Session is the common interface for all session types.
type Session interface {
	// Start is the starting point for the new session.
	Start()
	// Close closes all connections (destroying session).
	Close()
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
