package events

// Source: https://github.com/openai/openai-python/blob/main/src/openai/types/beta/realtime/response_create_event.py

import "gpt4omini/types"

type Metadata = map[string]interface{}

type ResponseTool struct {
	Description *string `json:"description,omitempty"` // Guidance on when/how to call
	Name        *string `json:"name,omitempty"`        // Name of the function
	Parameters  any     `json:"parameters,omitempty"`  // JSON Schema parameters
	Type        *string `json:"type,omitempty"`        // Should be "function"
}

type ResponseCreateEvent struct {
	Type     string    `json:"type"`               // Must be "response.create"
	EventID  *string   `json:"event_id,omitempty"` // Optional client-generated ID
	Response *Response `json:"response,omitempty"` // The response parameters
}

type Response struct {
	Conversation            *string                               `json:"conversation,omitempty"` // "auto", "none", or nil
	Input                   []types.ConversationItemWithReference `json:"input,omitempty"`
	Instructions            *string                               `json:"instructions,omitempty"`
	MaxResponseOutputTokens any                                   `json:"max_response_output_tokens,omitempty"` // int or "inf"
	Metadata                Metadata                              `json:"metadata,omitempty"`
	Modalities              []string                              `json:"modalities,omitempty"`          // "text" or "audio"
	OutputAudioFormat       *string                               `json:"output_audio_format,omitempty"` // "pcm16", "g711_ulaw", "g711_alaw"
	Temperature             *float64                              `json:"temperature,omitempty"`
	ToolChoice              *string                               `json:"tool_choice,omitempty"`
	Tools                   []ResponseTool                        `json:"tools,omitempty"`
	Voice                   *string                               `json:"voice,omitempty"` // "alloy", "ash", etc.
}
