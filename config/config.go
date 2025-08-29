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
	FileName = "config.yaml"
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
	if _, err := os.Stat(filePath); err != nil {
		fmt.Println("Config file doesn't exist, creating default config file")
		if err = createDefaultConfigFile(); err != nil {
			fmt.Println(err)
		}
		return nil
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("failed to read config file: %v\n", err)
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
	if err := os.WriteFile(FileName, data, 0644); err != nil {
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
	return filepath.Join(dir, FileName)
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
