package types

/*
Copied from the source as it is.
Source: https://github.com/openai/openai-python/blob/main/src/openai/types/conversations/message.py
*/

type Message struct {

	// Id The unique ID of the message
	Id string `json:"id,omitempty"`

	// Content The content of the message
	Content []Content `json:"content"`

	/*
		Role The role of the message.
		One of `unknown`, `user`, `assistant`, `system`, `critic`, `discriminator`, `developer`, or `tool`.
	*/
	Role string `json:"role"`

	/*
		Status The status of item.
		One of `in_progress`, `completed`, or `incomplete`. Populated when items are returned via API.
	*/
	Status string `json:"status,omitempty"`

	// Type The type of the message. Always set to `message`.
	Type string `json:"type"`
}

type ClientMessage struct {
	Type     string   `json:"type"`
	Response Response `json:"response"`
}
