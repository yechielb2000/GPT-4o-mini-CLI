package types

/*
Source: https://github.com/openai/openai-python/blob/main/src/openai/types/
There is no exact definition for content. it's just a TypeAlias for all the defined content types.
We declare an interface for the content and then make content types out of it.
*/

// Content The interface for all content types
type Content interface {
	ContentType() string
}

// BaseContent shares params across all content types
type BaseContent struct {
	// Type The type of the input item.
	Type string `json:"type"`
}

func (b BaseContent) ContentType() string {
	return b.Type
}

type TextContent struct {
	BaseContent
	Type string `json:"type"`
	Text string `json:"text"`
}

type FunctionCallContent struct {
	BaseContent
	Type   string `json:"type,omitempty"`
	Name   string `json:"name,omitempty"`
	Output any    `json:"output,omitempty"` // or result?
}
