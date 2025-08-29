package builders

import (
	"fmt"
	"gpt4omini/events"
	"gpt4omini/types"
)

const ClientName = "user"

func NewClientTextMessage(text string) *types.ClientMessage {
	return &types.ClientMessage{
		Type: events.ResponseCreate,
		Response: types.Response{
			Modalities: []string{types.TextModality},
			Input: []types.ConversationItem{
				{
					Type: types.MessageItem,
					Role: ClientName,
					Content: []types.Content{
						{Text: text, Type: types.InputTextItem},
					},
				},
			},
		},
	}
}

func NewClientFunctionCallResultItem(item types.ConversationItem, result string) events.ConversationItemCreateEvent {
	fmt.Println("new item call id", item.CallID)
	return events.ConversationItemCreateEvent{
		Type: events.ConversationItemCreateEventType,
		Item: types.ConversationItem{
			Type:   types.FunctionCallOutputItem,
			CallID: item.CallID,
			Output: result,
		},
	}
}
