package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	GitHub struct {
		APIKey string `yaml:"api_key"`
	} `yaml:"github"`
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
	FollowerUpdateInterval time.Duration `yaml:"follower_check_interval"`
	Logging                struct {
		FileLevel    string `yaml:"file_level"`
		ConsoleLevel string `yaml:"console_level"`
		FilePath     string `yaml:"file_path"`
	} `yaml:"logging"`
}

var AppConfig Config

// LoadConfig загружает конфигурацию из файла
func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Error reading config file", "error", err)
		return err
	}

	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		slog.Error("Error parsing config file", "error", err)
		return err
	}

	if AppConfig.FollowerUpdateInterval == 0 {
		err = fmt.Errorf("FollowerUpdateInterval cannot be 0 seconds")
		slog.Error("Invalid FollowerUpdateInterval in config file", "error", err)
		return err
	}

	slog.Info("Loaded config successfully")
	return nil
}
