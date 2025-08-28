package types

type ToolProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ToolParameters struct {
	Type       string                  `json:"type"`
	Properties map[string]ToolProperty `json:"properties"`
}

type Tool struct {
	Description *string        `json:"description,omitempty"` // Function description and usage guidance
	Name        *string        `json:"name,omitempty"`        // Name of the function
	Parameters  ToolParameters `json:"parameters,omitempty"`  // JSON Schema parameters
	Type        *string        `json:"type,omitempty"`        // Must be "function"
}

type ToolResult struct {
	Type   string `json:"type,omitempty"`
	Name   string `json:"name,omitempty"`
	Output any    `json:"output,omitempty"`
}
