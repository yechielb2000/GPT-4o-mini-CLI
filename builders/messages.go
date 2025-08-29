package builders

import (
	"gpt4omini/events"
	"gpt4omini/types"
)

const ClientName = "user"

func NewClientTextMessage(text string) *types.ClientMessage {
	return &types.ClientMessage{
		Type: events.ResponseCreate,
		Response: types.Response{
			Modalities: []string{"text"},
			Input: []types.Message{
				{
					Type: "message",
					Role: ClientName,
					Content: []types.Content{
						{Text: text, Type: "input_text"},
					},
				},
			},
		},
	}
}

func NewClientToolResult(name string, result any) types.Content {
	return types.Content{
		Type:   "tool_result",
		Name:   name,
		Output: result,
	}
}

func NewClientItemCreateEvent(item types.ConversationItem) events.ConversationItemCreateEvent {
	return events.ConversationItemCreateEvent{
		Type: events.ConversationItemCreateEventType,
		Item: item,
	}
}
