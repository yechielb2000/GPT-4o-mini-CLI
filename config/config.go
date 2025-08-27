package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"sync"
)

var (
	instance *Config
	once     sync.Once
)

const (
	// ConfigFilePath is like this for tests TODO: read from path
	ConfigFilePath = "./tests/config.yaml"

	RealtimeSessionsPath = "/v1/realtime/sessions"
	RealtimePath         = "/v1/realtime"
)

func NewApiConfig() *Config {
	config := &Config{}
	data, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(data, config); err != nil {
		panic(fmt.Errorf("failed to parse YAML: %w", err))
	}
	return config
}

func GetConfig() *Config {
	once.Do(func() {
		instance = NewApiConfig()
	})
	return instance
}

// GetURL provides full url.URL object. path is provided manually.
func GetURL(path string) url.URL {
	config := GetConfig()
	return url.URL{
		Scheme:   config.Api.Schema,
		Host:     config.Api.Host,
		Path:     path,
		RawQuery: fmt.Sprintf("model=%s", config.Model.Name),
	}
}
