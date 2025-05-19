package logger

import (
	"log/slog"
	"os"
)

// Log is a publicly accessible logger instance.
var Log *slog.Logger

func init() {

	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	defer file.Close()
	// Initialize the logger with a JSON handler.
	Log = slog.New(slog.NewJSONHandler(file, nil))
}
