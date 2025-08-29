package types

/*
ClientMessage holds the Response which defines the content as the message.
Open-ai generate duplicated objects all the time.
The Message is basically the ConversationItem
Source: https://github.com/openai/openai-python/blob/main/src/openai/types/conversations/message.py
In the source there is the fields of Message which is exactly how ConversationItem looks like if we omit the rest.
*/
type ClientMessage struct {
	Type     EventType `json:"type"`
	Response Response  `json:"response"`
}

func NewClientMessage(items []ConversationItem) ClientMessage {
	return ClientMessage{
		Type: ResponseCreateEvent,
		Response: Response{
			Modalities: []Modality{TextModality},
			Input:      items,
		},
	}
}
