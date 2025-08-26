package session

import (
	"errors"
	"fmt"
	"sync"
)

var (
	instance *SessionsManager
	once     sync.Once
)

type SessionsManager struct {
	sessions map[string]*Session
}

func GetSessionsManager() *SessionsManager {
	once.Do(func() {
		instance = NewSessionsManager()
	})
	return instance
}

func NewSessionsManager() *SessionsManager {
	return &SessionsManager{make(map[string]*Session)}
}

func (sm *SessionsManager) AddSession(session *Session) {
	// this can probably override existing id
	sm.sessions[session.SessionID] = session
}

func (sm *SessionsManager) RemoveSession(id string) {
	session, ok := sm.sessions[id]
	if ok {
		session.Close()
		delete(sm.sessions, id)
	}
}

func (sm *SessionsManager) GetSession(id string) (*Session, error) {
	session, ok := sm.sessions[id]
	if !ok {
		return nil, errors.New(fmt.Sprintf("session %s not found", id))
	}
	if session.HasExpired() {
		return nil, errors.New(fmt.Sprintf("session %s expired", id))
	}
	return session, nil
}

func (sm *SessionsManager) Sessions() map[string]*Session {
	return sm.sessions
}
