package logger

import (
    "log/slog"
    "os"
)

// Log is a publicly accessible logger instance.
var logger *slog.Logger

func init() {
    // Initialize the logger with a JSON handler.
    logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}