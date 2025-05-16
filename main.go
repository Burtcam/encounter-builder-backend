package main

import (
	"fmt"
	"log/slog"
	"os"
	"github.com/Burtcam/encounter-builder-backend/logger"
	"github.com/Burtcam/encounter-builder-backend/utils"
)

// struct encounter {
// 	difficulty string
// 	pSize      int
// 	level      int
// }


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
