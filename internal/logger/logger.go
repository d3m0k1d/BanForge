package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func New(debug bool) *Logger {
	var level slog.Level
	if debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return &Logger{
		Logger: slog.New(handler),
	}
}
