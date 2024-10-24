package config

import (
	"LeakGFinder/scanner"
	"encoding/json"
	"os"
)

type Config struct {
	IgnoreDirs []string          `json:"ignore_dirs"`
	Patterns   map[string]string `json:"patterns"`
}

func LoadConfig(filepath string) Config {
	file, err := os.Open(filepath)
	if err != nil {
		scanner.HandleError("Não foi possível abrir o arquivo: %s\n", filepath)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		scanner.HandleError("Não foi possível decodificar a configuração\n", "")
	}

	if err := validateConfig(config); err != nil {
		scanner.HandleError("Configuração Inválida\n", "")
	}

	return config
}

func validateConfig(config Config) error {
	if len(config.IgnoreDirs) == 0 {
		scanner.HandleError("A lista de diretórios a ignorar não pode estar vazia\n", "")
	}

	return nil
}
