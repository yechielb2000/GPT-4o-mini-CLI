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

type SessionResponse struct {
	Type    string `json:"type"`
	EventId string `json:"event_id"`
	Session struct {
		Object                   string        `json:"object"`
		Id                       string        `json:"id"`
		Model                    string        `json:"model"`
		Modalities               []string      `json:"modalities"`
		Instructions             string        `json:"instructions"`
		Voice                    string        `json:"voice"`
		OutputAudioFormat        string        `json:"output_audio_format"`
		Tools                    []interface{} `json:"tools"`
		ToolChoice               string        `json:"tool_choice"`
		Temperature              float64       `json:"temperature"`
		MaxResponseOutputTokens  string        `json:"max_response_output_tokens"`
		TurnDetection            []interface{} `json:"turn_detection"`
		Speed                    float64       `json:"speed"`
		Tracing                  interface{}   `json:"tracing"`
		Prompt                   interface{}   `json:"prompt"`
		ExpiresAt                int           `json:"expires_at"`
		InputAudioNoiseReduction interface{}   `json:"input_audio_noise_reduction"`
		InputAudioFormat         string        `json:"input_audio_format"`
		InputAudioTranscription  interface{}   `json:"input_audio_transcription"`
		Include                  interface{}   `json:"include"`
	} `json:"session"`
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
