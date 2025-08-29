package config

import (
	"fmt"
	"os"
)

const DefaultConfig = `api:
    key: API_KEY
    host: api.openai.com
    schema: wss
model:
    name: gpt-4o-realtime-preview
    instruction: You are a golden fish assistant that likes to make fire under water, You speak only in english.
`

func createDefaultConfigFile() error {
	if err := os.WriteFile(getConfigFilePath(), []byte(DefaultConfig), 0644); err != nil {
		return fmt.Errorf("failed to create config.yaml: %w", err)
	}

	fmt.Println("config.yaml created with default values. Please update it with your API key.")
	return nil
}
