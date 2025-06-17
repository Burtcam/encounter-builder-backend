package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"strconv"

	"github.com/Burtcam/encounter-builder-backend/config"
	"github.com/Burtcam/encounter-builder-backend/logger"
	"github.com/Burtcam/encounter-builder-backend/utils"
	"github.com/Burtcam/encounter-builder-backend/writeMonsters"
)

func CalculatexpBudget(w http.ResponseWriter, r *http.Request) {
	psizeStr := r.URL.Query().Get("psize")
	Psize, err := strconv.Atoi(psizeStr)
	if err != nil {
		http.Error(w, "Invalid psize parameter", http.StatusBadRequest)
		return
	}
	difficulty := r.URL.Query().Get("difficulty")

	xpBudget, err := utils.GetXpBudget(difficulty, Psize)

	if err != nil {
		logger.Log.Error("unable to calculate xp budget %w", err)
	}
	fmt.Fprintf(w, "XP budget is %d", xpBudget)
}

func getMonstersbyLevel(cfg config.Config, ctx context.Context, w http.ResponseWriter, r *http.Request) {
	queries := writeMonsters.New(cfg.DBPool)
	minStr := r.URL.Query().Get("min")
	maxStr := r.URL.Query().Get("max")
	//
	monsters, err := queries.GetMonstersByLevelRange(ctx, writeMonsters.GetMonstersByLevelRangeParams{
		Level:   utils.NewText(minStr),
		Level_2: utils.NewText(maxStr),
	})
	if err != nil {

	}
}

func api(ctx context.Context, cfg config.Config) error {
	http.HandleFunc("/calculatebudget", CalculatexpBudget)
	http.HandleFunc("/MonstersInLevelRange", getMonstersbyLevel(cfg, ctx))
	logger.Log.Info("listening on :5000")
	http.ListenAndServe(":5000", nil)
	return nil
}

func main() {
	cfg := config.Load()
	logger.Log.Info("Backend Initializing",
		slog.String("version", "1.0.0"),
		slog.String("env", "development"),
	)

	//setup the sync cron for the db.
	go utils.ManageDBSync(*cfg)
	// //TODO Remove this else everytime the ap starts it'll rebuild the db.
	// err := utils.KickOffSync(*cfg)
	// if err != nil {
	// 	logger.Log.Error(err.Error())
	//}

	err := api()
	if err != nil {
		logger.Log.Error("Unable to initialize APIS %w", err)
	}
}
