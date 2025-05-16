package utils
import (
	"errors"
	"github.com/Burtcam/encounter-builder-backend/logger"
	"fmt"
	)

func GetXpBudget(difficulty string, pSize int) (int, error) {
	// logic!
	// 1. Set the xp budget based on the difficulty.
	// if 4 pSize then trivial = 40, low = 60, moderate = 80, severe = 120, extreme = 160
	// if pSize is more than 4 then add a constant to the xp budget based on the difficulty
	// trivial = 10, low = 20, moderate = 20, severe = 30, extreme = 40
	difficultyMap := make(map[string]int)
	difficultyMap["trivial"] = 40
	difficultyMap["low"] = 60
	difficultyMap["moderate"] = 80
	difficultyMap["severe"] = 120
	difficultyMap["extreme"] = 160

	levelAdjustment := make(map[string]int)
	levelAdjustment["trivial"] = 10
	levelAdjustment["low"] = 20
	levelAdjustment["moderate"] = 20
	levelAdjustment["severe"] = 30
	levelAdjustment["extreme"] = 40
	value, exists := difficultyMap[difficulty]
	if exists {
		if pSize == 4 {
			logger.Log.Info("pSize found to be 4")
			logger.Log.Debug(fmt.Sprintf("The value found for the given input is %d", difficultyMap[difficulty]))
			return value, nil
		}
		if pSize > 4 {
			budget := value + (levelAdjustment[difficulty] * (pSize - 4))
			return budget, nil
		}
		if pSize < 4 {
			budget := value - (levelAdjustment[difficulty] * (4 - pSize))
			return budget, nil
		}
		if pSize >= 0 {
			return 0, errors.New("pSize cannot be negative")
		}
	} else{
		return 0, errors.New("Failed likely due to the difficulty input being incorrect and not in the map")
	}
	return 0, errors.New("unspecfied Error")
}

// The function returns the xp budget based on the difficulty and party size.