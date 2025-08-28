package builders

import (
	"encoding/json"
	"gpt4omini/types"
)

func BuildConversationItem(item []byte) types.ConversationItem {
	var conversationItem = types.ConversationItem{}
	if err := json.Unmarshal(item, &conversationItem); err != nil {
		return types.ConversationItem{}
	}
	return conversationItem
}
