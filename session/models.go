package session

import "gpt4omini/types"

type ConfigureModelRequest struct {
	Modalities   []string `json:"modalities"`
	Model        string   `json:"model"`
	Instructions string   `json:"instructions"`
}

type ConfigureModelResponse struct {
	Id           string       `json:"id"`
	Object       string       `json:"object"`
	ClientSecret ClientSecret `json:"client_secret"`
}

type ClientSecret struct {
	Value     string `json:"value"`
	ExpiresAt int64  `json:"expires_at"`
}

type Message struct {
	Type     string          `json:"type"`
	Response MessageResponse `json:"response"`
}

type MessageResponse struct {
	Modalities []string        `json:"modalities"`
	Input      []types.Message `json:"input"`
}
