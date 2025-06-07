package logger

import (
	"io"
	"log/slog"
	"os"
)

// Log is the global logger instance that writes to both STDOUT and a file.
var Log *slog.Logger

func init() {
	// 1. Open or create the log file in append mode.
	//    os.O_CREATE: creates the file if it doesn't exist.
	//    os.O_WRONLY: open the file write-only.
	//    os.O_APPEND: append new log entries to the end of the file.
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// If the log file can't be opened, panic because logging is critical.
		panic("failed to open log file: " + err.Error())
	}

	// 2. Create an io.MultiWriter that duplicates writes to both STDOUT and the log file.
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// 3. Create a new slog Handler that writes structured log messages to the multiwriter.
	//    You can choose a text or JSON handler. Here we use NewTextHandler for readability.
	//    The HandlerOptions allow you to set parameters such as the minimum log level.
	handler := slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelInfo, // Change this to adjust verbosity.
	})

	// 4. Create the global structured logger using the handler.
	Log = slog.New(handler)
}
