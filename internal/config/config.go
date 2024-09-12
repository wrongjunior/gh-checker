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
	FollowerUpdateInterval time.Duration `yaml:"follower_check_interval"` // Изменен ключ
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

	if AppConfig.FollowerUpdateInterval == 0 {
		log.Fatalf("FollowerUpdateInterval cannot be 0 seconds. Please check your config file.")
	}

	log.Printf("Loaded config with FollowerUpdateInterval: %v seconds", AppConfig.FollowerUpdateInterval.Seconds())
}
