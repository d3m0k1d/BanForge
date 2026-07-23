package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*slog.Logger
}

func New(debug bool) *Logger {
	logDir := "/var/log/banforge"
	var level slog.Level
	if debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	var output io.Writer = os.Stdout
	if err := os.MkdirAll(logDir, 0750); err == nil {
		fileWriter := &lumberjack.Logger{
			Filename:   filepath.Join(logDir, "banforge.log"),
			MaxSize:    500,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
		}
		output = io.MultiWriter(fileWriter, os.Stdout)
	}

	handler := slog.NewTextHandler(output, &slog.HandlerOptions{
		Level: level,
	})

	return &Logger{
		Logger: slog.New(handler),
	}
}
