package events

import "gpt4omini/types"

const ResponseDoneEventType = "response.done"

type ResponseDoneEvent struct {
	EventID  string                 `json:"event_id"` // Unique ID of the server event
	Response types.RealtimeResponse `json:"response"` // The response resource
	Type     string                 `json:"type"`     // Must be "response.done"
}
