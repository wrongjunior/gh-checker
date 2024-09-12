package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	GitHub struct {
		APIKey string `yaml:"api_key"`
	} `yaml:"github"`
}

var AppConfig Config

// LoadConfig загружает конфигурацию из YAML файла
func LoadConfig(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
}
