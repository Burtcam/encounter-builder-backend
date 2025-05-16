package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Burtcam/encounter-builder-backend/utils"
)

// struct encounter {
// 	difficulty string
// 	pSize      int
// 	level      int
// }

var logger *slog.Logger

func init() {
    // Initialize the global logger using a JSON handler.
    logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func main() {
	logger.Info("Backend Initializing",
		slog.String("version", "1.0.0"),
		slog.String("env", "development"),
	)

	xpBudget, err := utils.GetXpBudget("Trivial", 4)
	if err != nil {
		logger.Error("Error occurred in someFunction", slog.String("error", err.Error()))
	} else {
		logger.Info(fmt.Sprintf(fmt.Sprintf("xpBudget succesfully calculated as: %d", xpBudget)))
	}
}
