package main

import (
	"gh-checker/internal/config"
	"gh-checker/internal/database"
	"gh-checker/internal/handlers"
	"gh-checker/internal/lib/logger"
	"gh-checker/internal/services"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Загружаем конфигурацию
	if err := config.LoadConfig("config.yaml"); err != nil {
		// Поскольку логгер еще не инициализирован, используем slog для логирования ошибки
		slog.Error("Error loading config", "error", err)
		os.Exit(1) // Завершение программы при ошибке загрузки конфигурации
	}

	// Инициализируем логгер
	fileLevel, err := logger.ParseLogLevel(config.AppConfig.Logging.FileLevel)
	if err != nil {
		slog.Error("Invalid file log level", "error", err)
		os.Exit(1) // Завершение программы при ошибке парсинга уровня логирования
	}
	loggingConfig := logger.LogConfig{
		FileLevel: fileLevel,
		FilePath:  config.AppConfig.Logging.FilePath,
	}

	if err = logger.InitializeLogger(loggingConfig); err != nil {
		slog.Error("Failed to initialize logger", "error", err)
		os.Exit(1) // Завершение программы при ошибке инициализации логгера
	}
	defer logger.CloseLogger()

	// Теперь можно использовать кастомный логгер
	logger.Info("Configuration and logger initialized")

	// Проверка наличия GitHub API Key
	githubAPIKey := config.AppConfig.GitHub.APIKey
	if githubAPIKey == "" {
		logger.Error("GitHub API key not found in config file", nil)
		os.Exit(1) // Завершение программы при отсутствии API ключа
	}
	services.SetGitHubAPIKey(githubAPIKey)

	// Инициализация базы данных
	if err := database.InitDB(config.AppConfig.Database.Path); err != nil {
		logger.Error("Failed to initialize database", err)
		os.Exit(1) // Завершение программы при ошибке инициализации базы данных
	}
	logger.Info("Database initialized")

	// Настройка роутера
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/api/subscribe", handlers.SubscribeHandler) // TODO: сделать на /check-followers
	r.Post("/check-star", handlers.StarCheckHandler)

	logger.Info("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Error("Server failed to start", err)
		os.Exit(1) // Завершение программы при ошибке старта сервера
	}
}
