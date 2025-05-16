package utils
import ("errors"
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

	if pSize == 4 {
		logger.Info("pSize found to be 4")
		return difficultyMap[difficulty], nil
	}
	if pSize > 4 {
		budget := difficultyMap[difficulty] + (levelAdjustment[difficulty] * (pSize - 4))
		return budget, nil
	}
	if pSize < 4 {
		budget := difficultyMap[difficulty] - (levelAdjustment[difficulty] * (4 - pSize))
		return budget, nil
	}
	return 0, errors.New("Failed")
}
