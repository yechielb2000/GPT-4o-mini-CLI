package types

type RealtimeResponseStatus struct {
	Error  *Error  `json:"error,omitempty"`  // Error details if the response failed
	Reason *string `json:"reason,omitempty"` // "turn_detected", "client_cancelled", "max_output_tokens", "content_filter"
	Type   *string `json:"type,omitempty"`   // "completed", "cancelled", "incomplete", "failed"
}
