package types

import (
	"encoding/json"
)

// Source: https://github.com/openai/openai-python/blob/main/src/openai/types/beta/realtime/conversation_item_with_reference.py

type ConversationItemType string

const (
	MessageItem            ConversationItemType = "message"
	FunctionCallItem       ConversationItemType = "function_call"
	FunctionCallOutputItem ConversationItemType = "function_call_output"
	ReferenceItem          ConversationItemType = "item_reference"
	InputTextItem          ConversationItemType = "input_text"
)

type ConversationItem struct {
	ID        string               `json:"id,omitempty"`        // Unique ID or reference to previous item
	Arguments string               `json:"arguments,omitempty"` // Arguments for function call
	CallID    string               `json:"call_id,omitempty"`   // Function call ID
	Content   []Content            `json:"content,omitempty"`   // Message content
	Name      string               `json:"name,omitempty"`      // Function name
	Object    string               `json:"object,omitempty"`    // Should be "realtime.item"
	Output    string               `json:"output,omitempty"`    // Output for function_call_output
	Role      string               `json:"role,omitempty"`      // "user", "assistant", "system"
	Status    string               `json:"status,omitempty"`    // "completed", "incomplete", "in_progress"
	Type      ConversationItemType `json:"type,omitempty"`      // "message", "function_call", "function_call_output", "item_reference"
}

func (c *ConversationItem) GetArguments() (map[string]interface{}, error) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(c.Arguments), &m); err != nil {
		return nil, err
	}
	return m, nil
}
