package localonlyutils

import (
	"github.com/Burtcam/encounter-builder-backend/config"
	"github.com/Burtcam/encounter-builder-backend/logger"
	"github.com/Burtcam/encounter-builder-backend/utils"
)

func LocalDataLoad(cfg config.Config) error {
	fileList := []string{"files/foundryvtt-pf2e-4cbdaa3/packs/wardens-of-wildwood-bestiary/book-3-shepherd-of-decay/primal-warden-of-zibik.json",
		"files/foundryvtt-pf2e-4cbdaa3/packs/age-of-ashes-bestiary/book-1-hellknight-hill/charau-ka.json",
		"files/foundryvtt-pf2e-4cbdaa3/packs/age-of-ashes-bestiary/book-1-hellknight-hill/town-hall-fire.json",
		"files/foundryvtt-pf2e-4cbdaa3/packs/age-of-ashes-bestiary/book-4-fires-of-the-haunted-city/saggorak-poltergeist.json",
		"files/foundryvtt-pf2e-4cbdaa3/packs/age-of-ashes-bestiary/book-5-against-the-scarlet-triad/scarlet-triad-enforcer.json",
		"files/foundryvtt-pf2e-4cbdaa3/packs/age-of-ashes-bestiary/book-4-fires-of-the-haunted-city/king-harral.json"}
	for i := range len(fileList) {
		err := utils.LoadEachJSON(cfg, fileList[i])
		if err != nil {
			logger.Log.Error("Failed to write file %s. Err: %v", fileList[i], err)
		}
	}
	return nil
}
