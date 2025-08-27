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
	sessions map[string]Session
}

func GetSessionsManager() *SessionsManager {
	once.Do(func() {
		instance = NewSessionsManager()
	})
	return instance
}

func NewSessionsManager() *SessionsManager {
	return &SessionsManager{make(map[string]Session)}
}

func (sm *SessionsManager) AddSession(session Session) {
	id := session.GetID()
	if _, exists := sm.sessions[id]; exists {
		fmt.Printf("Session %s already exists. Skipping.\n", id)
		return
	}
	sm.sessions[session.GetID()] = session
}

func (sm *SessionsManager) RemoveSession(id string) {
	session, ok := sm.sessions[id]
	if ok {
		session.Close()
		delete(sm.sessions, id)
	}
}

func (sm *SessionsManager) GetSession(id string) (Session, error) {
	session, ok := sm.sessions[id]
	if !ok {
		return nil, errors.New(fmt.Sprintf("session %s not found", id))
	}
	if session.HasClientSecretExpired() {
		sm.RemoveSession(id)
		return nil, errors.New(fmt.Sprintf("session %s has expired", id))
	}
	return session, nil
}

func (sm *SessionsManager) Sessions() map[string]Session {
	return sm.sessions
}
