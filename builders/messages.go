package builders

import (
	"gpt4omini/events"
	"gpt4omini/types"
)

const (
	ClientName = "user"
	ModelName  = "assistant"
)

func NewClientTextConversationItem(text string) types.ConversationItem {
	return types.ConversationItem{
		Type: types.MessageItem,
		Role: ClientName,
		Content: []types.Content{
			{Text: text, Type: types.InputTextItem},
		},
	}
}

func NewClientFunctionCallConversationItem(item types.ConversationItem, result string) types.ConversationItem {
	return types.ConversationItem{
		Type:   types.FunctionCallOutputItem,
		CallID: item.CallID,
		Output: result,
	}
}

func NewClientConversationEvent(items []types.ConversationItem) types.ClientMessage {
	return types.ClientMessage{
		Type: events.ResponseCreate,
		Response: types.Response{
			Modalities: []string{types.TextModality},
			Input:      items,
		},
	}
}
