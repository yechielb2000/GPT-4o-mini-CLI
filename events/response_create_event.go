package events

import "gpt4omini/types"

type ResponseCreateEvent struct {
	Type     string          `json:"type"`               // Must be "response.create"
	EventID  *string         `json:"event_id,omitempty"` // Optional client-generated ID
	Response *types.Response `json:"response,omitempty"` // The response parameters
}
