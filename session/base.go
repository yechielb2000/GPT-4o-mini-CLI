package session

import (
	"context"
	"encoding/json"
	"gpt4omini/config"
	"gpt4omini/types"
	"sync"
	"time"
)

const ExitSessionFunctionName = "exit_session"

var cfg = config.GetConfig()

// BaseSession contains fields shared across all sessions.
type BaseSession struct {
	ID               string
	Type             string
	outgoingMessages chan types.ConversationItem
	functionCalls    chan types.ConversationItem
	incomingEvents   chan types.Event
	readyForInput    chan struct{}
	clientSecret     types.ClientSecret
	createdAt        time.Time
	ctx              context.Context
	cancel           context.CancelFunc
	wg               sync.WaitGroup
	mu               sync.Mutex
	conversation     []types.ConversationItem
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

// GetCreationTime get the creation time of the session
func (bs *BaseSession) GetCreationTime() types.ClientSecret {
	return bs.clientSecret
}

// AddToConversation any item that you want to model to remember
func (bs *BaseSession) AddToConversation(item types.ConversationItem) {
	item.Object = "" // The model do not expect to see it, so I removed it.
	bs.conversation = append(bs.conversation, item)
}

// GetConversation of all the session messages history
func (bs *BaseSession) GetConversation() []types.ConversationItem {
	return bs.conversation
}

// NewClientMessage that creates the client message based on the conversation history
func (bs *BaseSession) NewClientMessage() types.ClientMessage {
	return types.ClientMessage{
		Type: types.ResponseCreateEvent,
		Response: types.Response{
			Modalities: []types.Modality{types.TextModality},
			Input:      bs.GetConversation(),
		},
	}
}

// String print nice the object
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
