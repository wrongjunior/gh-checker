package logger

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
)

var (
	logger      *slog.Logger
	once        sync.Once
	logFile     *os.File
	fileHandler slog.Handler
	fileLevel   = &slog.LevelVar{} // Уровень логирования для файла
)

// LogConfig - конфигурация логгера
type LogConfig struct {
	FileLevel slog.Level // Уровень логирования для файла
	FilePath  string     // Путь до файла логов
}

// InitializeLogger инициализирует логгер с переданными параметрами
func InitializeLogger(config LogConfig) error {
	var err error
	once.Do(func() {
		// Устанавливаем уровень логирования
		fileLevel.Set(config.FileLevel)

		// Инициализация файла для логирования
		if config.FilePath != "" {
			logFile, err = os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				return
			}
			fileHandler = slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: fileLevel})
			logger = slog.New(fileHandler)
		}
	})

	return err
}

// SetFileLevel изменяет уровень логирования для файла
func SetFileLevel(level slog.Level) {
	fileLevel.Set(level)
}

// CloseLogger закрывает файл логирования (если был открыт)
func CloseLogger() {
	if logFile != nil {
		_ = logFile.Close()
	}
}

// Info выводит информационные логи
func Info(msg string, attrs ...any) {
	logger.Info(msg, attrs...)
}

// Error выводит ошибки
func Error(msg string, err error, attrs ...any) {
	logger.Error(msg, append(attrs, slog.String("error", err.Error()))...)
}

// Debug выводит debug-логи
func Debug(msg string, attrs ...any) {
	logger.Debug(msg, attrs...)
}

// Warn выводит предупреждения
func Warn(msg string, attrs ...any) {
	logger.Warn(msg, attrs...)
}

// ParseLogLevel преобразует строковое представление уровня логирования в slog.Level
func ParseLogLevel(level string) (slog.Level, error) {
	switch level {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("неверный уровень логирования: %s", level)
	}
}
