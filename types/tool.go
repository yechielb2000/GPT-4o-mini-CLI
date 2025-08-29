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
	Description *string        `json:"description,omitempty"`
	Name        *string        `json:"name,omitempty"`
	Parameters  ToolParameters `json:"parameters,omitempty"`
	Type        *string        `json:"type,omitempty"`
}
