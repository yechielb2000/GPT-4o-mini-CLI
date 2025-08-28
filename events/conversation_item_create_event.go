package events

import (
	"gpt4omini/types"
)

type ConversationItemCreateEvent struct {
	Item           types.ConversationItem `json:"item"`                       // The item to add
	Type           string                 `json:"type"`                       // Must be "conversation.item.create"
	EventID        *string                `json:"event_id,omitempty"`         // Optional client-generated ID
	PreviousItemID *string                `json:"previous_item_id,omitempty"` // Where to insert the new item
}
