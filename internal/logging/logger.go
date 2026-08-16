package logging

import (
	"log/slog"
	"os"
)

func New(level string) *slog.Logger {

	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{},
	)

	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger
}