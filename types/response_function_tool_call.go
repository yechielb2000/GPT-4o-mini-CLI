package types

const ResponseFunctionToolCallType = "function_call"

type ResponseFunctionToolCall struct {
	Arguments string  `json:"arguments"`        // JSON string of arguments
	CallID    string  `json:"call_id"`          // Unique ID of the function call
	Name      string  `json:"name"`             // Name of the function to run
	Type      string  `json:"type"`             // Always "function_call"
	ID        *string `json:"id,omitempty"`     // Optional unique ID
	Status    *string `json:"status,omitempty"` // "in_progress", "completed", "incomplete"
}
