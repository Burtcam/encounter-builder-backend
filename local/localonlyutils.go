package localonlyutils

import (
	"github.com/Burtcam/encounter-builder-backend/config"
	"github.com/Burtcam/encounter-builder-backend/utils"
)

func LocalDataLoad(cfg config.Config) error {
	fileList := []string{"files/foundryvtt-pf2e-4cbdaa3/packs/wardens-of-wildwood-bestiary/book-3-shepherd-of-decay/primal-warden-of-zibik.json"}
	for i := range len(fileList) {
		utils.LoadEachJSON(cfg, fileList[i])
	}
	return nil
}
