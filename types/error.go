package types

type Error struct {
	Message string  `json:"message,omitempty"`  // Human-readable error message
	Type    string  `json:"type,omitempty"`     // Type of error, e.g., "invalid_request_error"
	Code    *string `json:"code,omitempty"`     // Optional error code
	EventID *string `json:"event_id,omitempty"` // Optional client event ID
	Param   *string `json:"param,omitempty"`    // Optional related parameter
}
