package session

type ClientSecret struct {
	Value     string `json:"value"`
	ExpiresAt int64  `json:"expires_at"`
}

type CreateSessionHTTPRequest struct {
	Modalities   []string `json:"modalities"`
	Model        string   `json:"model"`
	Instructions string   `json:"instructions"`
}

type CreateSessionHTTPResponse struct {
	Id           string       `json:"id"`
	Object       string       `json:"object"`
	Model        string       `json:"model"`
	Modalities   []string     `json:"modalities"`
	Instructions string       `json:"instructions"`
	ExpiresAt    int          `json:"expires_at"`
	ClientSecret ClientSecret `json:"client_secret"`
}

type MessageResponseConfig struct {
	Modalities   []string `json:"modalities"`
	Instructions string   `json:"instructions"`
	Conversation string   `json:"conversation,omitempty"`
}

type MessageRequest struct {
	Type     string                `json:"type"`
	Response MessageResponseConfig `json:"response"`
}
