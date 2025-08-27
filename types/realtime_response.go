package types

const RealtimeResponseObjectType = "realtime.response"

type RealtimeResponse struct {
	ID                *string                 `json:"id,omitempty"`                  // Unique ID of the response
	ConversationID    *string                 `json:"conversation_id,omitempty"`     // ID of the conversation
	MaxOutputTokens   any                     `json:"max_output_tokens,omitempty"`   // int or "inf"
	Metadata          Metadata                `json:"metadata,omitempty"`            // Optional key-value pairs
	Modalities        []string                `json:"modalities,omitempty"`          // "text", "audio"
	Object            *string                 `json:"object,omitempty"`              // "realtime.response"
	Output            []ConversationItem      `json:"output,omitempty"`              // Generated output items
	OutputAudioFormat *string                 `json:"output_audio_format,omitempty"` // "pcm16", "g711_ulaw", "g711_alaw"
	Status            *string                 `json:"status,omitempty"`              // "completed", "cancelled", etc.
	StatusDetails     *RealtimeResponseStatus `json:"status_details,omitempty"`      // Extra status info
	Temperature       *float64                `json:"temperature,omitempty"`         // Sampling temperature
	Usage             *RealtimeResponseUsage  `json:"usage,omitempty"`               // Usage statistics
}
