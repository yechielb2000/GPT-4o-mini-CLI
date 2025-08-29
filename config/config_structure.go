package config

import "gpt4omini/types"

type Api struct {
	Key    string `json:"key"`
	Host   string `json:"host"`
	Schema string `json:"schema"`
}

type Model struct {
	Name        string       `json:"name"`
	Instruction string       `json:"instruction"`
	Tools       []types.Tool `json:"tools,omitempty"`
}

type Config struct {
	Api   Api   `json:"config"`
	Model Model `json:"model"`
}
