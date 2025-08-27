package events

const ConversationItemDeleteEventType = "conversation.item.delete"

type ConversationItemDeleteEvent struct {
	ItemID  string  `json:"item_id"`            // ID of the item to delete
	Type    string  `json:"type"`               // Must be "conversation.item.delete"
	EventID *string `json:"event_id,omitempty"` // Optional client-generated ID
}
