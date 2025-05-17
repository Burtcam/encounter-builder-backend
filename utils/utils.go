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
	"archive/tar"
	"compress/gzip"
	"path/filepath"
	"strings"
	"encoding/json"
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
    outFile, err := os.Create("repo_archive.tar.gz")
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
func extractTarball(tarFile string, destDir string) error {
    // Open the tar file
    file, err := os.Open(tarFile)
    if err != nil {
        return fmt.Errorf("failed to open tar file: %w", err)
    }
    defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}

	defer gzipReader.Close()



    // Create a tar reader
    tarReader := tar.NewReader(gzipReader)

    // Read each file inside the tar archive
    for {
        header, err := tarReader.Next()
        if err == io.EOF {
            break // End of archive
        }
        if err != nil {
            logger.Log.Error("error reading tar file: %w", err.Error())
			return err
        }

        // Determine the file path
        filePath := destDir + "/" + header.Name

        switch header.Typeflag {
        case tar.TypeDir:
            // Create directories
            if err := os.MkdirAll(filePath, os.FileMode(header.Mode)); err != nil {
				// logger.Log.Error(("error creating directory: %w", err.Error()))
                return err
            }
        case tar.TypeReg:
            // Extract regular files
            outFile, err := os.Create(filePath)
            if err != nil {
				// logger.Log.Error(("error creating file: %w", err.Error()))
                return err
            }
            defer outFile.Close()

            // Copy the file contents from the archive
            if _, err := io.Copy(outFile, tarReader); err != nil {
                return fmt.Errorf("error writing file: %w", err)
            }
        }
    }
    return nil
}

func GetListofJSON(dir string) ([]string, error) {
    var fileList []string

    // Walk through the directory
    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Check if the file has a .json extension
        if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
            fileList = append(fileList, path)
        }

        return nil
    })

    return fileList, err
}

func LoadEachJSON(path string) error {
	data, err := os.ReadFile(path)
	if err != nil{
		logger.Log.Error(err.Error())
		return err
	}
	var payload map[string]interface{}
    err = json.Unmarshal(data, &payload)
    if err != nil {
        logger.Log.Error("Error during Unmarshal(): ", err)
    }
	if payload["type"] == "npc"{ 
		jsonData, err := json.Marshal(payload)
    	if err != nil {
        	logger.Log.Error("Error encoding JSON:", err)
    	}
		logger.Log.Info(string(jsonData))
		os.Exit(1)
	}
	return nil
}


func KickOffSync(cfg config.Config) error {
	logger.Log.Debug("About to go get the archive")
	err := GetRepoArchive(cfg)
	if err != nil {
		logger.Log.Error("Sync failed at archive download")
	}
	err = extractTarball("repo_archive.tar.gz", "files")
	if err != nil {
		logger.Log.Error("Failed to unpack tarball")
		logger.Log.Error(err.Error())
	}
	fileList, err := GetListofJSON("./files")
	if err != nil {
		logger.Log.Error("Failed to get the list of files to process")
		logger.Log.Error(err.Error())
	}
	logger.Log.Info(fmt.Sprintf("%v", fileList))
	// 

	for _, value := range fileList {
		err = LoadEachJSON(value)
		if err != nil {
			logger.Log.Error(err.Error())
		}

	}
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
	{}
	return nil// Keep the program running
}