package events

import "gpt4omini/types"

type EventError struct {
	Error   types.Error `json:"error,omitempty"`    // Details of the error.
	Type    string      `json:"type,omitempty"`     // Type The event type, must be `error`.
	EventID *string     `json:"event_id,omitempty"` // The unique ID of the server event.
}
