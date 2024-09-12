package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GitHub struct {
		APIKey string `yaml:"api_key"`
	} `yaml:"github"`
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
	FollowerUpdateInterval time.Duration `yaml:"follower_update_interval"`
}

var AppConfig Config

func LoadConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// Предполагается, что интервал указывается в секундах
	AppConfig.FollowerUpdateInterval = AppConfig.FollowerUpdateInterval * time.Second
}
