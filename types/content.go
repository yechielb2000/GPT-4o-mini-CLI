package types

/*
Source: https://github.com/openai/openai-python/blob/main/src/openai/types/
I make one content for all content types.
*/

// Content one for all content items
type Content struct {
	Type       ConversationItemType `json:"type,omitempty"`
	Name       string               `json:"name,omitempty"`
	Output     interface{}          `json:"output,omitempty"`
	Text       string               `json:"text,omitempty"`
	ToolCallID string               `json:"tool_call_id,omitempty"`
}
