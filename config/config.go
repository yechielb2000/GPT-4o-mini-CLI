package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

var (
	instance *Config
	once     sync.Once
)

const (
	// ConfigFileName is like this for tests TODO: read from path
	ConfigFileName = "config.yaml"

	RealtimeSessionsPath = "/v1/realtime/sessions"
	RealtimePath         = "/v1/realtime"
)

func GetConfig() *Config {
	once.Do(func() {
		instance = NewApiConfig()
	})
	return instance
}

func NewApiConfig() *Config {
	config := &Config{}
	filePath := getConfigFilePath()
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(data, config); err != nil {
		panic(fmt.Errorf("failed to parse YAML: %w", err))
	}
	return config
}

func (c *Config) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	if err := os.WriteFile(ConfigFileName, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func getConfigFilePath() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(exePath)
	return filepath.Join(dir, ConfigFileName)
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
