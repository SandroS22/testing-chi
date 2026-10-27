package config

import (
	"log"
	"os"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Port string `yaml:"port"`
}

func LoadConfig(s string) Config {
	log.Println("Loading config from", s)
	config := Config{}

	err := readYAMLConfig(s, &config)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	return config
}

func readYAMLConfig(filename string, config *Config) error {
	ymlFile, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(ymlFile, config)
	if err != nil {
		return err
	}

	return nil
}
