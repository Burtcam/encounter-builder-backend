package utils
import (
	"errors"
	"github.com/Burtcam/encounter-builder-backend/logger"
	"fmt"
	"net/http"
	"os"
	"time"
	"github.com/robfig/cron/v3"
	"github.com/Burtcam/encounter-builder-backend/config"
	"io"
	)

func GetXpBudget(difficulty string, pSize int) (int, error) {
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
	logger.Log.Info(fmt.Sprintf("%d", pSize))
	value, exists := difficultyMap[difficulty]
	if exists {
		if pSize <= 0 {
			return 0, errors.New("pSize cannot be negative")
		}
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
	} else{
		return 0, errors.New("Failed likely due to the difficulty input being incorrect and not in the map")
	}
	return 0, errors.New("unspecfied Error")
}
func GetRepoArchive(cfg config.Config) (error){
	client := &http.Client{}
	// call to the repoUrl and get the archive downloaded. 
	req, err := http.NewRequest("GET", cfg.REPO_URL, nil)
    if err != nil {
        return fmt.Errorf("error creating request: %w", err)
    }
 	req.Header.Set("User-Agent", "MyGoClient/1.0")
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.GH_TOKEN))
    req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	    // Execute the request
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    // Handle the response (assume we want to save it)
    outFile, err := os.Create("repo_archive.tar")
    if err != nil {
        return fmt.Errorf("error creating file: %w", err)
    }
    defer outFile.Close()

    // Copy response body to the file
    _, err = io.Copy(outFile, resp.Body)
    if err != nil {
        return fmt.Errorf("error saving archive: %w", err)
    }

    fmt.Println("Repository archive downloaded successfully!")
    return nil
}
func KickOffSync(cfg config.Config) error {
	err := GetRepoArchive(cfg)
	if err != nil {
		logger.Log.Error("Sync failed at archive download")
	}
	// else {
		
	// }
	return nil
}
func ManageDBSync(cfg config.Config) error {
	// Go routine to wait until a certain unix time (3AM PST by default but managed by config) then go get the new tarball every week and sync it to the db 
	loc, err := time.LoadLocation("America/Los_Angeles")
    if err != nil {
        return fmt.Errorf("invalid timezone: %w", err)
    }

    // Define cron syntax: Run at 3 AM PST every Tuesday
    c := cron.New(cron.WithLocation(loc))
    _, err = c.AddFunc("0 3 * * 2", func() { KickOffSync(cfg) }) // "2" represents Tuesday in cron syntax
    if err != nil {
        return err
    }

    c.Start()
    fmt.Println("Scheduled job to run every Tuesday at 3 AM PST")
    select {} // Keep the program running
	return nil
}