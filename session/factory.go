package session

import (
	"errors"
	"fmt"
)

// Using this factory to create all session types available.

var Factory = map[string]func() (Session, error){
	"realtime": func() (Session, error) { return NewRealtimeSession() },
}

func NewSessionByType(type_ string) (Session, error) {
	if factory, ok := Factory[type_]; ok {
		return factory()
	}
	return nil, errors.New(fmt.Sprintf("Session type \"%s\" not support", type_))
}
