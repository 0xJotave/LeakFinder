package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	IgnoreDirs []string          `json:"ignore_dirs"`
	Patterns   map[string]string `json:"patterns"`
}

func LoadConfig(filepath string) (Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
