package types

type InputTokenDetails struct {
	AudioTokens  *int `json:"audio_tokens,omitempty"`  // Number of audio tokens used
	CachedTokens *int `json:"cached_tokens,omitempty"` // Number of cached tokens used
	TextTokens   *int `json:"text_tokens,omitempty"`   // Number of text tokens used
}

type OutputTokenDetails struct {
	AudioTokens *int `json:"audio_tokens,omitempty"` // Number of audio tokens used
	TextTokens  *int `json:"text_tokens,omitempty"`  // Number of text tokens used
}

type RealtimeResponseUsage struct {
	InputTokenDetails  *InputTokenDetails  `json:"input_token_details,omitempty"`  // Details about input tokens
	InputTokens        *int                `json:"input_tokens,omitempty"`         // Total input tokens
	OutputTokenDetails *OutputTokenDetails `json:"output_token_details,omitempty"` // Details about output tokens
	OutputTokens       *int                `json:"output_tokens,omitempty"`        // Total output tokens
	TotalTokens        *int                `json:"total_tokens,omitempty"`         // Total input + output tokens
}
