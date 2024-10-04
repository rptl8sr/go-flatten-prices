package logger

import (
	"log"
	"log/slog"
	"os"
	"path"
	"time"
)

type logger struct {
	log *slog.Logger
}

var (
	globalLogger     logger
	globalFileLogger logger
)

func Init(lvl slog.Level, lvlFile slog.Level, logsDir string) {
	consoleLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
	globalLogger.log = consoleLogger

	logFilePath := path.Join(logsDir, time.Now().Format("2006-01-02")+"_log.txt")

	dirPath := path.Dir(logFilePath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if e := os.MkdirAll(dirPath, 0755); e != nil {
			log.Fatalf("Не могу создать директорию для лог файла: %v", e)
		}
	}

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Не могу создать или открыть лог файл %v", err)
	}

	globalFileLogger.log = slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: lvlFile}))
}

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

func Debug(msg string, args ...interface{}) {
	globalFileLogger.log.Debug(msg, args...)
	globalLogger.log.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	globalFileLogger.log.Info(msg, args...)
	globalLogger.log.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	globalFileLogger.log.Warn(msg, args...)
	globalLogger.log.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	globalFileLogger.log.Error(msg, args...)
	globalLogger.log.Error(msg, args...)
}
