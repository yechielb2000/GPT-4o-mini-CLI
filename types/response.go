package types

// Source: https://github.com/openai/openai-python/blob/main/src/openai/types/beta/realtime/response_create_event.py

type ResponseTool struct {
	Description *string `json:"description,omitempty"`
	Name        *string `json:"name,omitempty"`
	Parameters  any     `json:"parameters,omitempty"`
	Type        *string `json:"type,omitempty"`
}

type Response struct {
	Conversation            *string            `json:"conversation,omitempty"`
	Input                   []ConversationItem `json:"input,omitempty"`
	Instructions            *string            `json:"instructions,omitempty"`
	MaxResponseOutputTokens any                `json:"max_response_output_tokens,omitempty"`
	Metadata                Metadata           `json:"metadata,omitempty"`
	Modalities              []Modality         `json:"modalities,omitempty"`
	OutputAudioFormat       *string            `json:"output_audio_format,omitempty"`
	Temperature             *float64           `json:"temperature,omitempty"`
	ToolChoice              *string            `json:"tool_choice,omitempty"`
	Tools                   []ResponseTool     `json:"tools,omitempty"`
	Voice                   *string            `json:"voice,omitempty"`
}
