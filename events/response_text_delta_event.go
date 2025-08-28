package events

import "gpt4omini/types"

type Event struct {
	ContentIndex int                    `json:"content_index"` // Index in the item's content array
	Delta        string                 `json:"delta"`         // Text delta
	EventID      string                 `json:"event_id"`      // Unique server event ID
	ItemID       string                 `json:"item_id"`       // ID of the item
	OutputIndex  int                    `json:"output_index"`  // Index of the output item in the response
	ResponseID   string                 `json:"response_id"`   // ID of the response
	Type         string                 `json:"type"`          // Must be "response.text.delta"
	Item         types.ConversationItem `json:"item,omitempty"`
	Obfuscation  string                 `json:"obfuscation,omitempty"`
}
