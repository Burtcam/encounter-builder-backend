package main

import (
	"fmt"
	"log/slog"

	"github.com/Burtcam/encounter-builder-backend/config"
	localonlyutils "github.com/Burtcam/encounter-builder-backend/local"
	"github.com/Burtcam/encounter-builder-backend/logger"
	"github.com/Burtcam/encounter-builder-backend/utils"
)

// struct encounter {
// 	difficulty string
// 	pSize      int
// 	level      int
// }

func main() {
	cfg := config.Load()
	logger.Log.Info("Backend Initializing",
		slog.String("version", "1.0.0"),
		slog.String("env", "development"),
	)

	xpBudget, err := utils.GetXpBudget("trivial", 4)
	if err != nil {
		logger.Log.Error("Error occurred in someFunction", slog.String("error", err.Error()))
	} else {
		logger.Log.Info((fmt.Sprintf("xpBudget succesfully calculated as: %d", xpBudget)))
	}
	//setup the sync cron for the db.
	go utils.ManageDBSync(*cfg)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	// //TODO Remove this else everytime the ap starts it'll rebuild the db.
	// err = utils.KickOffSync(*cfg)
	// if err != nil {
	// 	logger.Log.Error(err.Error())
	// }

	err = localonlyutils.LocalDataLoad(*cfg)
	if err != nil {
		logger.Log.Error(err.Error())
	}

}
