package builders

import (
	"gpt4omini/events"
	"gpt4omini/types"
)

func NewClientTextMessage(text string) *types.ClientMessage {
	return &types.ClientMessage{
		Type: events.ResponseCreate,
		Response: types.Response{
			Modalities: []string{"text"},
			Input: []types.Message{
				{
					Type: "message",
					Role: "user",
					Content: []types.Content{
						types.TextContent{Text: text, Type: "input_text"},
					},
				},
			},
		},
	}
}
