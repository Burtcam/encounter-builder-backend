package utils

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Burtcam/encounter-builder-backend/config"
	"github.com/Burtcam/encounter-builder-backend/logger"
	"github.com/Burtcam/encounter-builder-backend/structs"
	"github.com/Burtcam/encounter-builder-backend/writeMonsters"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/gjson"
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
	} else {
		return 0, errors.New("Failed likely due to the difficulty input being incorrect and not in the map")
	}
	return 0, errors.New("unspecfied Error")
}
func GetRepoArchive(cfg config.Config) error {
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
	logger.Log.Info(fmt.Sprintf("Response Return code is: %s", resp.Status))
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

	logger.Log.Info(fmt.Sprintf("Repository archive downloaded successfully!"))
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
			// logger.Log.Error("error reading tar file: %w", err.Error())
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

func ParseFoundJson(data string) structs.Monster {
	monster := ParseCoreData(string(data))
	//Parse items and pass it just the items list then attach the return values to monster.
	ItemsList := gjson.Get(string(data), "items").String()
	var spells []structs.Spell
	monster.FreeActions, monster.Actions, monster.Reactions, monster.Passives, monster.SpellCasting, spells, monster.Melees, monster.Ranged, monster.Inventory = ParseItems(ItemsList)

	AssignSpell(&spells, &monster.SpellCasting)
	return monster
}

func PrepMonsterParams(monster structs.Monster) writeMonsters.InsertMonsterParams {

	monsterParams := writeMonsters.InsertMonsterParams{

		Name:             pgtype.Text{String: monster.Name, Valid: true},
		Level:            pgtype.Text{String: monster.Level, Valid: true},
		FocusPoints:      pgtype.Int4{String: monster.FocusPoints, value: true},
		TraitsRarity:     pgtype.Text{String: monster.Traits.Rarity},
		TraitsSize:       pgtype.Text{String: monster.Traits.Size},
		AttrStr:          pgtype.Text{String: monster.Attributes.Str},
		AttrDex:          pgtype.Text{String: monster.Attributes.Dex},
		AttrCon:          pgtype.Text{String: monster.Attributes.Con},
		AttrWis:          pgtype.Text{String: monster.Attributes.Wis},
		AttrInt:          pgtype.Text{String: monster.Attributes.Int},
		AttrCha:          pgtype.Text{String: monster.Attributes.Cha},
		SavesFort:        pgtype.Text{String: monster.Saves.Fort},
		SavesFortDetail:  pgtype.Text{String: monster.Saves.FortDetail},
		SavesRef:         pgtype.Text{String: monster.Saves.Ref},
		SavesRefDetail:   pgtype.Text{String: monster.Saves.RefDetail},
		SavesWill:        pgtype.Text{String: monster.Saves.Will},
		SavesWillDetail:  pgtype.Text{String: monster.Saves.WillDetail},
		SavesException:   pgtype.Text{String: monster.Saves.Exception},
		AcValue:          pgtype.Text{String: monster.AClass.Value},
		AcDetail:         pgtype.Text{String: monster.AClass.Detail},
		HpValue:          pgtype.int4{Int: monster.HP.Value},
		HpDetail:         pgtype.Text{String: monster.HP.Detail},
		PerceptionMod:    pgtype.Text{String: monster.Perception.Mod},
		PerceptionDetail: pgtype.Text{String: monster.Perception.Detail},
	}
	return monsterParams
}

func writeImmunites(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {

	for i := 0; i < len(monster.Immunities); i++ {
		err := &queries.InsertMonsterImmunity(ctx, writeMonsters.InsertMonsterImmunityParams{
			MonsterID: pgtype.Int4{Int: id, Valid: true},
			Immunity:  pgtype.Text{String: monster.Immunities[i], Valid: true},
		})
		if err != nil {
			logger.Log.Error("Failed to insert immunity %s for monster ID %d: %v", monster.Immunities[i], id, err)
			return fmt.Errorf("failed to insert immunity %s for monster ID %d: %w", monster.Immunities[i], id, err)
		}
		logger.Log.Info(fmt.Sprintf("Succesfully inserted immunity %s for monster ID %d", monster.Immunities[i], id))
	}

	return nil
}

func ProcessWeakAndResist(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {
	var damageBlockList []structs.DamageBlock
	damageBlockList = append(damageBlockList, monster.Weaknesseses)
	damageBlockList = append(damageBlockList, monster.Resistances)

	for i := 0; i < len(damageBlockList); i++ {
		id, err := &queries.InsertMonsterDamageModifier(ctx, writeMonsters.InsertMonsterDamageModifier{
			MonsterID: pgtype.Int4{Int: id, Valid: true},
			Modifier:  pgtype.Text{String: damageBlockList[i].Modifier, Valid: true},
			Detail:    pgtype.Text{String: damageBlockList[i].Detail, Valid: true},
			Type:      pgtype.Text{String: damageBlockList[i].Type, Valid: true},
		})
	}

	return nil
}

func WriteMonsterToDb(monster structs.Monster, cfg config.Config) error {
	ctx := context.Background()
	// ✅ 2. Begin a transaction
	tx, err := cfg.DBPool.Begin(ctx)
	if err != nil {
		logger.Log.Error("failed to start transaction: %v", err)
	}

	//prep main params
	monsterParams := PrepMonsterParams(monster)

	queries := writeMonsters.New(cfg.DBPool)

	id, err := queries.InsertMonster(ctx, monsterParams)
	if err != nil {
		logger.Log.Error("Failed to insert monster %v", err)
	}
	logger.Log.Info(fmt.Sprintf("Succesfully started the transaction for ID %d", id))
	//for each immunities
	err = writeImmunites(ctx, queries, monster, id)
	if err != nil {
		logger.Log.Error("Failed to process immunities: %w", err)
	}
	err = ProcessWeakAndResist(ctx, queries, monster, id)
	if err != nil {
		logger.Log.Error("Failed to process weaknesses and resistances: %w", err)
	}
	return nil

}

func LoadEachJSON(cfg config.Config, path string) error {
	fmt.Println("Path is :", path)
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	if gjson.Get(string(data), "type").String() == "npc" {
		fmt.Println("Found a monster")
		monster := ParseCoreData(string(data))
		//Parse items and pass it just the items list then attach the return values to monster.
		ItemsList := gjson.Get(string(data), "items").String()
		var spells []structs.Spell
		monster.FreeActions, monster.Actions, monster.Reactions, monster.Passives, monster.SpellCasting, spells, monster.Melees, monster.Ranged, monster.Inventory = ParseItems(ItemsList)

		AssignSpell(&spells, &monster.SpellCasting)

		fmt.Println(monster)

		err = WriteMonsterToDb(monster, cfg)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("Unable to  write %s, to db %w", monster.Name, err))
		}
		// err = parseJSON(data)
		// if err != nil {
		// 	logger.Log.Error(fmt.Sprintf("Error Parsing file %s", path))
		// }
		// // WRite it out to a json
		// // err = os.WriteFile("example-monster.json", jsonData, 0644)
		// // if err != nil {
		// // 	logger.Log.Error("Error writting JSON:", err)
		// // }
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
		err = LoadEachJSON(cfg, value)
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
	{
	}
	return nil // Keep the program running
}
