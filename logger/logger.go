package logger

import (
    "log/slog"
    "os"
)

// Log is a publicly accessible logger instance.
var Log *slog.Logger

func init() {
    // Initialize the logger with a JSON handler.
    Log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}