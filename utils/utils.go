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
	"github.com/google/uuid"

	//"github.com/jackc/pgx/pgtype"

	// "github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/gjson"
)

func NewText(value string) pgtype.Text {
	return pgtype.Text{String: value, Valid: value != ""}
}

func NewInt4(value int) pgtype.Int4 {
	return pgtype.Int4{Int32: int32(value), Valid: true}
}

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
		return 0, errors.New("failed likely due to the difficulty input being incorrect and not in the map")
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

	logger.Log.Info(fmt.Sprintln("Repository archive downloaded successfully!"))
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

func ParseFoundJson(data string) (structs.Monster, error) {
	monster := ParseCoreData(string(data))
	//Parse items and pass it just the items list then attach the return values to monster.
	ItemsList := gjson.Get(string(data), "items")
	var spells []structs.Spell
	var err error
	monster.FreeActions, monster.Actions, monster.Reactions, monster.Passives, monster.SpellCasting, spells, monster.Melees, monster.Ranged, monster.Inventory, err = ParseItems(ItemsList)
	if err != nil {
		return monster, err
	}
	AssignSpell(&spells, &monster.SpellCasting)
	return monster, err
}

func PrepMonsterParams(monster structs.Monster) writeMonsters.InsertMonsterParams {

	monsterParams := writeMonsters.InsertMonsterParams{

		Name:             monster.Name,
		Level:            pgtype.Text{String: monster.Level, Valid: true},
		FocusPoints:      pgtype.Int4{Int32: int32(monster.FocusPoints), Valid: true},
		TraitsRarity:     pgtype.Text{String: monster.Traits.Rarity, Valid: true},
		TraitsSize:       pgtype.Text{String: monster.Traits.Size, Valid: true},
		AttrStr:          pgtype.Text{String: monster.Attributes.Str, Valid: true},
		AttrDex:          pgtype.Text{String: monster.Attributes.Dex, Valid: true},
		AttrCon:          pgtype.Text{String: monster.Attributes.Con, Valid: true},
		AttrWis:          pgtype.Text{String: monster.Attributes.Wis, Valid: true},
		AttrInt:          pgtype.Text{String: monster.Attributes.Int, Valid: true},
		AttrCha:          pgtype.Text{String: monster.Attributes.Cha, Valid: true},
		SavesFort:        pgtype.Text{String: monster.Saves.Fort, Valid: true},
		SavesFortDetail:  pgtype.Text{String: monster.Saves.FortDetail, Valid: true},
		SavesRef:         pgtype.Text{String: monster.Saves.Ref, Valid: true},
		SavesRefDetail:   pgtype.Text{String: monster.Saves.RefDetail, Valid: true},
		SavesWill:        pgtype.Text{String: monster.Saves.Will, Valid: true},
		SavesWillDetail:  pgtype.Text{String: monster.Saves.WillDetail, Valid: true},
		SavesException:   pgtype.Text{String: monster.Saves.Exception, Valid: true},
		AcValue:          pgtype.Text{String: monster.AClass.Value, Valid: true},
		AcDetail:         pgtype.Text{String: monster.AClass.Detail, Valid: true},
		HpValue:          pgtype.Int4{Int32: int32(monster.HP.Value), Valid: true},
		HpDetail:         pgtype.Text{String: monster.HP.Detail, Valid: true},
		PerceptionMod:    pgtype.Text{String: monster.Perception.Mod, Valid: true},
		PerceptionDetail: pgtype.Text{String: monster.Perception.Detail, Valid: true},
	}
	return monsterParams
}

func writeImmunites(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {

	for i := 0; i < len(monster.Immunities); i++ {
		err := queries.InsertMonsterImmunities(ctx, writeMonsters.InsertMonsterImmunitiesParams{
			MonsterID: pgtype.Int4{Int32: id, Valid: true},
			Immunity:  pgtype.Text{String: monster.Immunities[i], Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to insert immunity %s for monster ID %d: %w", monster.Immunities[i], id, err)
		}
		logger.Log.Info(fmt.Sprintf("Succesfully inserted immunity %s for monster ID %d", monster.Immunities[i], id))
	}

	return nil
}

func ProcessWeakAndResist(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {

	for i := 0; i < len(monster.Weaknesses); i++ {
		DamageModifierID, err := queries.InsertMonsterDamageModifier(ctx, writeMonsters.InsertMonsterDamageModifierParams{
			MonsterID:        pgtype.Int4{Int32: id, Valid: true},
			ModifierCategory: pgtype.Text{String: "weakness", Valid: true},
			Value:            pgtype.Int4{Int32: int32(monster.Weaknesses[i].Value), Valid: true},
			DamageType:       pgtype.Text{String: monster.Weaknesses[i].Type, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to add damage modifier to DB %w", err)
		}
		// if exceptions len > 0
		if len(monster.Weaknesses[i].Exceptions) > 0 {
			for j := 0; j < len(monster.Weaknesses[i].Exceptions); j++ {
				err = queries.InsertMonsterModifierExceptions(ctx, writeMonsters.InsertMonsterModifierExceptionsParams{
					ModifierID: pgtype.Int4{Int32: DamageModifierID, Valid: true},
					Exception:  pgtype.Text{String: monster.Weaknesses[i].Exceptions[j], Valid: true},
				})
			}
		}
		if err != nil {
			return fmt.Errorf("failed to add damage modifier to DB %w", err)
		}
		// if exceptions len > 0
		if len(monster.Weaknesses[i].Double) > 0 {
			for k := 0; k < len(monster.Weaknesses[i].Double); k++ {
				err = queries.InsertMonsterModifierExceptions(ctx, writeMonsters.InsertMonsterModifierExceptionsParams{
					ModifierID: pgtype.Int4{Int32: DamageModifierID, Valid: true},
					Exception:  pgtype.Text{String: monster.Weaknesses[i].Double[k], Valid: true},
				})
			}
		}
		if err != nil {
			return fmt.Errorf("failed to add damage modifier to DB %w", err)
		}
	}
	for i := 0; i < len(monster.Resistances); i++ {
		DamageModifierID, err := queries.InsertMonsterDamageModifier(ctx, writeMonsters.InsertMonsterDamageModifierParams{
			MonsterID:        pgtype.Int4{Int32: id, Valid: true},
			ModifierCategory: pgtype.Text{String: "resistance", Valid: true},
			Value:            pgtype.Int4{Int32: int32(monster.Resistances[i].Value), Valid: true},
			DamageType:       pgtype.Text{String: monster.Resistances[i].Type, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to add damage modifier to DB Resistances  %w", err)
		}
		// if exceptions len > 0
		if len(monster.Resistances[i].Exceptions) > 0 {
			for j := 0; j < len(monster.Resistances[i].Exceptions); j++ {
				err = queries.InsertMonsterModifierExceptions(ctx, writeMonsters.InsertMonsterModifierExceptionsParams{
					ModifierID: pgtype.Int4{Int32: DamageModifierID, Valid: true},
					Exception:  pgtype.Text{String: monster.Resistances[i].Exceptions[j], Valid: true},
				})
			}
		}
		if err != nil {
			return fmt.Errorf("failed to add damage modifier to DB Exceptions %w", err)
		}
		// if exceptions len > 0
		if len(monster.Resistances[i].Double) > 0 {
			for k := 0; k < len(monster.Resistances[i].Double); k++ {
				err = queries.InsertMonsterModifierExceptions(ctx, writeMonsters.InsertMonsterModifierExceptionsParams{
					ModifierID: pgtype.Int4{Int32: DamageModifierID, Valid: true},
					Exception:  pgtype.Text{String: monster.Resistances[i].Double[k], Valid: true},
				})
			}
		}
		if err != nil {
			return fmt.Errorf("failed to add damage modifier to DB %w", err)
		}
	}
	return nil
}

func ProcessLanguages(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {
	for i := 0; i < len(monster.Languages); i++ {
		err := queries.InsertMonsterLanguages(ctx, writeMonsters.InsertMonsterLanguagesParams{
			MonsterID: pgtype.Int4{Int32: id, Valid: true},
			Language:  pgtype.Text{String: monster.Languages[i], Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to write language %w", err)
		}
	}
	return nil
}

func ProcessSenses(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {
	for i := range len(monster.Senses) {
		err := queries.InsertMonsterSenses(ctx, writeMonsters.InsertMonsterSensesParams{
			MonsterID: pgtype.Int4{Int32: id, Valid: true},
			Name:      pgtype.Text{String: monster.Senses[i].Name, Valid: true},
			Range:     pgtype.Text{String: monster.Senses[i].Range, Valid: true},
			Acuity:    pgtype.Text{String: monster.Senses[i].Acuity, Valid: true},
			Detail:    pgtype.Text{String: monster.Senses[i].Detail, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("unable to write sense to DB %w", err)
		}
	}
	return nil
}

func ProcessMovements(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {
	for i := 0; i < len(monster.Movements); i++ {
		err := queries.InsertMonsterMovements(ctx, writeMonsters.InsertMonsterMovementsParams{
			MonsterID:    pgtype.Int4{Int32: id},
			MovementType: pgtype.Text{String: monster.Movements[i].Type, Valid: true},
			Speed:        pgtype.Text{String: monster.Movements[i].Speed, Valid: true},
			Notes:        pgtype.Text{String: monster.Movements[i].Notes, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to write Movements to DB %w", err)
		}
	}
	return nil
}

func ProcessSkills(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {
	for i := 0; i < len(monster.Skills); i++ {
		skillId, err := queries.InsertMonsterSkills(ctx, writeMonsters.InsertMonsterSkillsParams{
			MonsterID: pgtype.Int4{Int32: id, Valid: true},
			Name:      pgtype.Text{String: monster.Skills[i].Name, Valid: true},
			Value:     pgtype.Int4{Int32: int32(monster.Skills[i].Value), Valid: true},
		})
		if err != nil {
			return fmt.Errorf("unable to write skill %w", err)
		}
		if len(monster.Skills[i].Specials) > 0 {
			for j := 0; j < len(monster.Skills[i].Specials); j++ {
				err := queries.InsertMonsterSkillSpecials(ctx, writeMonsters.InsertMonsterSkillSpecialsParams{
					SkillID:    pgtype.Int4{Int32: skillId, Valid: true},
					Value:      pgtype.Int4{Int32: int32(monster.Skills[i].Specials[j].Value), Valid: true},
					Label:      pgtype.Text{String: monster.Skills[i].Specials[j].Label, Valid: true},
					Predicates: monster.Skills[i].Specials[j].Predicates,
				})
				if err != nil {
					return fmt.Errorf("unable to write skill specials %w", err)
				}
			}
		}
	}
	return nil
}

func ProcessAction(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {
	for i := 0; i < len(monster.Actions); i++ {
		actionId, err := queries.InsertMonsterAction(ctx, writeMonsters.InsertMonsterActionParams{
			MonsterID:  pgtype.Int4{Int32: id, Valid: true},
			ActionType: pgtype.Text{String: "action", Valid: true},
			Name:       pgtype.Text{String: monster.Actions[i].Name, Valid: true},
			Text:       pgtype.Text{String: monster.Actions[i].Text, Valid: true},
			Actions:    pgtype.Text{String: monster.Actions[i].Actions, Valid: true},
			Category:   pgtype.Text{String: monster.Actions[i].Category, Valid: true},
			Rarity:     pgtype.Text{String: monster.Actions[i].Rarity, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("unable to process Monster Action %w", err)
		}
		for j := 0; j < len(monster.Actions[i].Traits); j++ {
			err := queries.InsertMonsterActionTraits(ctx, writeMonsters.InsertMonsterActionTraitsParams{
				MonsterActionID: pgtype.Int4{Int32: actionId, Valid: true},
				Trait:           pgtype.Text{String: monster.Actions[i].Traits[j], Valid: true},
			})
			if err != nil {
				return fmt.Errorf("unable to Process Traits for Actions %w", err)
			}
		}
	}
	return nil
}

func ProcessReaction(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {
	for i := 0; i < len(monster.Reactions); i++ {
		actionId, err := queries.InsertMonsterAction(ctx, writeMonsters.InsertMonsterActionParams{
			MonsterID:  pgtype.Int4{Int32: id, Valid: true},
			ActionType: pgtype.Text{String: "reaction", Valid: true},
			Name:       pgtype.Text{String: monster.Reactions[i].Name, Valid: true},
			Text:       pgtype.Text{String: monster.Reactions[i].Text, Valid: true},
			Category:   pgtype.Text{String: monster.Reactions[i].Category, Valid: true},
			Rarity:     pgtype.Text{String: monster.Reactions[i].Rarity, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("unable to process Monster Reaction %w", err)
		}
		for j := 0; j < len(monster.Reactions[i].Traits); j++ {
			err := queries.InsertMonsterActionTraits(ctx, writeMonsters.InsertMonsterActionTraitsParams{
				MonsterActionID: pgtype.Int4{Int32: actionId, Valid: true},
				Trait:           pgtype.Text{String: monster.Reactions[i].Traits[j], Valid: true},
			})
			if err != nil {
				return fmt.Errorf("unable to Process Traits for Reaction %w", err)
			}
		}
	}
	return nil
}

func ProcessPassive(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {
	for i := 0; i < len(monster.Passives); i++ {
		actionId, err := queries.InsertMonsterAction(ctx, writeMonsters.InsertMonsterActionParams{
			MonsterID:  pgtype.Int4{Int32: id, Valid: true},
			ActionType: pgtype.Text{String: "passive", Valid: true},
			Name:       pgtype.Text{String: monster.Passives[i].Name, Valid: true},
			Text:       pgtype.Text{String: monster.Passives[i].Text, Valid: true},
			Category:   pgtype.Text{String: monster.Passives[i].Category, Valid: true},
			Rarity:     pgtype.Text{String: monster.Passives[i].Rarity, Valid: true},
			Dc:         pgtype.Text{String: monster.Passives[i].DC, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("unable to process Monster Passive %w", err)
		}
		for j := 0; j < len(monster.Passives[i].Traits); j++ {
			err := queries.InsertMonsterActionTraits(ctx, writeMonsters.InsertMonsterActionTraitsParams{
				MonsterActionID: pgtype.Int4{Int32: actionId, Valid: true},
				Trait:           pgtype.Text{String: monster.Passives[i].Traits[j], Valid: true},
			})
			if err != nil {
				return fmt.Errorf("unable to Process Traits for Passive %w", err)
			}
		}
	}
	return nil
}

func ProcessAttacks(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {
	for i := range len(monster.Melees) {
		attackID, err := queries.InsertMonsterAttacks(ctx, writeMonsters.InsertMonsterAttacksParams{
			MonsterID:           pgtype.Int4{Int32: id, Valid: true},
			AttackCategory:      pgtype.Text{String: "melee", Valid: true},
			Name:                pgtype.Text{String: monster.Melees[i].Name, Valid: true},
			AttackType:          pgtype.Text{String: monster.Melees[i].Type, Valid: true},
			ToHitBonus:          pgtype.Text{String: monster.Melees[i].ToHitBonus, Valid: true},
			EffectsCustomString: pgtype.Text{String: monster.Melees[i].Effects.CustomString, Valid: true},
			EffectsValues:       monster.Melees[i].Effects.Value,
		})
		if err != nil {
			return fmt.Errorf("unable to write attack ID %w", err)
		}
		for j := range len(monster.Melees[i].DamageBlocks) {
			err = queries.InsertMonsterAttackDamageBlock(ctx, writeMonsters.InsertMonsterAttackDamageBlockParams{
				AttackID:   pgtype.Int4{Int32: attackID, Valid: true},
				DamageRoll: pgtype.Text{String: monster.Melees[i].DamageBlocks[j].DamageRoll, Valid: true},
				DamageType: pgtype.Text{String: monster.Melees[i].DamageBlocks[j].DamageType, Valid: true},
			})
		}
		if err != nil {
			return fmt.Errorf("unable to write damageblock %w", err)
		}
	}
	for i := range len(monster.Ranged) {
		attackID, err := queries.InsertMonsterAttacks(ctx, writeMonsters.InsertMonsterAttacksParams{
			MonsterID:           pgtype.Int4{Int32: id},
			AttackCategory:      pgtype.Text{String: "ranged"},
			Name:                pgtype.Text{String: monster.Ranged[i].Name},
			AttackType:          pgtype.Text{String: monster.Ranged[i].Type},
			ToHitBonus:          pgtype.Text{String: monster.Ranged[i].ToHitBonus},
			EffectsCustomString: pgtype.Text{String: monster.Ranged[i].Effects.CustomString},
			EffectsValues:       monster.Ranged[i].Effects.Value,
		})
		if err != nil {
			return fmt.Errorf("unable to write attack ID %w", err)
		}
		for j := range len(monster.Ranged[i].DamageBlocks) {
			err = queries.InsertMonsterAttackDamageBlock(ctx, writeMonsters.InsertMonsterAttackDamageBlockParams{
				AttackID:   pgtype.Int4{Int32: attackID},
				DamageRoll: pgtype.Text{String: monster.Ranged[i].DamageBlocks[j].DamageRoll},
				DamageType: pgtype.Text{String: monster.Ranged[i].DamageBlocks[j].DamageType},
			})
		}
		if err != nil {
			return fmt.Errorf("unable to write damageblock %w", err)
		}
	}
	return nil
}

func processSpellGeneric(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, spell structs.Spell) (string, error) {
	spellId, err := queries.InsertSpell(ctx, writeMonsters.InsertSpellParams{
		Name:                        pgtype.Text{String: spell.Name, Valid: true},
		CastLevel:                   pgtype.Text{String: spell.CastLevel, Valid: true},
		SpellBaseLevel:              pgtype.Text{String: spell.SpellBaseLevel, Valid: true},
		Description:                 pgtype.Text{String: spell.Description, Valid: true},
		Range:                       pgtype.Text{String: spell.Range, Valid: true},
		CastTime:                    pgtype.Text{String: spell.CastTime, Valid: true},
		CastRequirements:            pgtype.Text{String: spell.CastRequirements, Valid: true},
		Rarity:                      pgtype.Text{String: spell.Rarity, Valid: true},
		AtWill:                      pgtype.Bool{Bool: spell.AtWill, Valid: true},
		SpellCastingBlockLocationID: pgtype.Text{String: spell.SpellCastingBlockLocationID, Valid: true},
		Uses:                        pgtype.Text{String: spell.Uses, Valid: true},
		Targets:                     pgtype.Text{String: spell.Targets, Valid: true},
		Ritual:                      pgtype.Bool{Bool: spell.Ritual, Valid: true},
	})
	if err != nil {
		return spellId, fmt.Errorf("unable to write spell %w", err)
	}
	err = queries.InsertSpellArea(ctx, writeMonsters.InsertSpellAreaParams{
		SpellID:  pgtype.Text{String: spellId, Valid: true},
		AreaType: pgtype.Text{String: spell.Area.Type, Valid: true},
		Value:    pgtype.Text{String: spell.Area.Value, Valid: true},
		Detail:   pgtype.Text{String: spell.Area.Detail, Valid: true},
	})
	if err != nil {
		return spellId, fmt.Errorf("unable to write spell area %w", err)
	}
	err = queries.InsertSpellDuration(ctx, writeMonsters.InsertSpellDurationParams{
		SpellID:   pgtype.Text{String: spellId, Valid: true},
		Sustained: pgtype.Bool{Bool: spell.Duration.Sustained, Valid: true},
		Duration:  pgtype.Text{String: spell.Duration.Duration, Valid: true},
	})
	if err != nil {
		return spellId, fmt.Errorf("unable to write spell duration %w", err)
	}
	err = queries.InsertSpellDefences(ctx, writeMonsters.InsertSpellDefencesParams{
		SpellID: pgtype.Text{String: spellId, Valid: true},
		Save:    pgtype.Text{String: spell.Defense.Save, Valid: true},
		Basic:   pgtype.Bool{Bool: spell.Defense.Basic, Valid: true},
	})
	if err != nil {
		return spellId, fmt.Errorf("failed to write spell defence block %w", err)
	}

	return spellId, nil
}

func ProcessInnateMagic(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {
	// -- name: InsertInnateSpellCasting :one
	for i := range len(monster.SpellCasting.InnateSpellCasting) {
		castingId, err := queries.InsertInnateSpellCasting(ctx, writeMonsters.InsertInnateSpellCastingParams{
			MonsterID:      pgtype.Int4{Int32: id, Valid: true},
			Dc:             pgtype.Int4{Int32: int32(monster.SpellCasting.InnateSpellCasting[i].DC), Valid: true},
			Tradition:      pgtype.Text{String: monster.SpellCasting.InnateSpellCasting[i].Tradition, Valid: true},
			Mod:            pgtype.Text{String: monster.SpellCasting.InnateSpellCasting[i].Mod, Valid: true},
			SpellcastingID: pgtype.Text{String: monster.SpellCasting.InnateSpellCasting[i].ID, Valid: true},
			Description:    pgtype.Text{String: monster.SpellCasting.InnateSpellCasting[i].Description, Valid: true},
			Name:           pgtype.Text{String: monster.SpellCasting.InnateSpellCasting[i].Name, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("unable to write innatespellcasting %w", err)
		}
		for j := range len(monster.SpellCasting.InnateSpellCasting[i].SpellUses) {
			//For each spell use theres a spell, write it to spell table AND write to spell use table.
			spellId, err := processSpellGeneric(ctx, queries, monster, monster.SpellCasting.InnateSpellCasting[i].SpellUses[j].Spell)
			if err != nil {
				return fmt.Errorf("unable to process spell to db %w", err)
			}
			err = queries.InsertInnateSpellUse(ctx, writeMonsters.InsertInnateSpellUseParams{
				InnateSpellCastingID: pgtype.Int4{Int32: int32(castingId), Valid: true},
				SpellID:              pgtype.Text{String: spellId, Valid: true},
				Level:                pgtype.Int4{Int32: int32(monster.SpellCasting.InnateSpellCasting[i].SpellUses[j].Level), Valid: true},
				Uses:                 pgtype.Text{String: monster.SpellCasting.InnateSpellCasting[i].SpellUses[j].Uses, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("unable to write innatespelluse %w", err)
			}
		}
	}
	return nil
}

func ProcessFocusMagic(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {
	// -- name: InsertFocusSpellCasting :one
	for i := range len(monster.SpellCasting.FocusSpellCasting) {
		castingId, err := queries.InsertFocusSpellCasting(ctx, writeMonsters.InsertFocusSpellCastingParams{
			MonsterID:      pgtype.Int4{Int32: id, Valid: true},
			Dc:             pgtype.Int4{Int32: int32(monster.SpellCasting.FocusSpellCasting[i].DC), Valid: true},
			Tradition:      pgtype.Text{String: monster.SpellCasting.FocusSpellCasting[i].Tradition, Valid: true},
			Mod:            pgtype.Text{String: monster.SpellCasting.FocusSpellCasting[i].Mod, Valid: true},
			SpellcastingID: pgtype.Text{String: monster.SpellCasting.FocusSpellCasting[i].ID, Valid: true},
			Description:    pgtype.Text{String: monster.SpellCasting.FocusSpellCasting[i].Description, Valid: true},
			Name:           pgtype.Text{String: monster.SpellCasting.FocusSpellCasting[i].Name, Valid: true},
			CastLevel:      pgtype.Text{String: monster.SpellCasting.FocusSpellCasting[i].CastLevel, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("unable to write focus spellcasting %w", err)
		}
		for j := range len(monster.SpellCasting.FocusSpellCasting[i].FocusSpellList) {
			spellId, err := processSpellGeneric(ctx, queries, monster, monster.SpellCasting.FocusSpellCasting[i].FocusSpellList[j])
			if err != nil {
				return fmt.Errorf("unable to write focus spell %w", err)
			}
			// Write each spell associatation.
			err = queries.InsertFocusSpellsCasts(ctx, writeMonsters.InsertFocusSpellsCastsParams{
				FocusSpellCastingID: pgtype.Int4{Int32: castingId, Valid: true},
				SpellID:             pgtype.Text{String: spellId, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("unable to write focus spell casts %w", err)
			}
		}
	}
	return nil
}

func ProcessPreparedMagic(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {
	// -- name: InsertPreparedSpellCasting :one
	for i := range len(monster.SpellCasting.PreparedSpellCasting) {
		castingId, err := queries.InsertPreparedSpellCasting(ctx, writeMonsters.InsertPreparedSpellCastingParams{
			MonsterID:      pgtype.Int4{Int32: id},
			Dc:             pgtype.Int4{Int32: int32(monster.SpellCasting.PreparedSpellCasting[i].DC), Valid: true},
			Tradition:      pgtype.Text{String: monster.SpellCasting.PreparedSpellCasting[i].Tradition, Valid: true},
			Mod:            pgtype.Text{String: monster.SpellCasting.PreparedSpellCasting[i].Mod, Valid: true},
			SpellcastingID: pgtype.Text{String: monster.SpellCasting.PreparedSpellCasting[i].ID, Valid: true},
			Description:    pgtype.Text{String: monster.SpellCasting.PreparedSpellCasting[i].Description, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("unable to write prepared Spellcasting %w", err)
		}
		for j := range len(monster.SpellCasting.PreparedSpellCasting[i].Slots) {
			//For each spell use theres a spell, write it to spell table AND write to spell use table.
			spellId, err := processSpellGeneric(ctx, queries, monster, monster.SpellCasting.PreparedSpellCasting[i].Slots[j].Spell)
			if err != nil {
				return fmt.Errorf("unable to process spell to db %w", err)
			}
			err = queries.InsertPreparedSlots(ctx, writeMonsters.InsertPreparedSlotsParams{
				PreparedSpellCastingID: pgtype.Int4{Int32: int32(castingId), Valid: true},
				SpellID:                pgtype.Text{String: spellId, Valid: true},
				Level:                  pgtype.Text{String: monster.SpellCasting.PreparedSpellCasting[i].Slots[j].Level, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("unable to insert prepared spell slots %w", err)
			}
		}
	}
	return nil
}

func ProcessSpontaneousMagic(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {
	// -- name: InsertSpontaneousSpellCasting :one
	for i := range len(monster.SpellCasting.SpontaneousSpellCasting) {
		spellCastingId, err := queries.InsertSpontaneousSpells(ctx, writeMonsters.InsertSpontaneousSpellsParams{
			MonsterID: pgtype.Int4{Int32: id, Valid: true},
			Dc:        pgtype.Int4{Int32: int32(monster.SpellCasting.SpontaneousSpellCasting[i].DC), Valid: true},
			IDString:  pgtype.Text{String: monster.SpellCasting.SpontaneousSpellCasting[i].ID, Valid: true},
			Tradition: pgtype.Text{String: monster.SpellCasting.SpontaneousSpellCasting[i].Tradition, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to insertSpontaneousSpells %w", err)
		}
		for j := range len(monster.SpellCasting.SpontaneousSpellCasting[i].SpellList) {
			spellID, err := processSpellGeneric(ctx, queries, monster, monster.SpellCasting.SpontaneousSpellCasting[i].SpellList[j])
			if err != nil {
				return fmt.Errorf("failed to process generic spell %w", err)
			}
			err = queries.InsertSpontaneousSpellList(ctx, writeMonsters.InsertSpontaneousSpellListParams{
				SpontaneousSpellCastingID: pgtype.Int4{Int32: spellCastingId, Valid: true},
				SpellID:                   pgtype.Text{String: spellID, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("failed to insert spell List stuff %w", err)
			}
		}
		for k := range len(monster.SpellCasting.SpontaneousSpellCasting[i].Slots) {
			err := queries.InsertSpontaneousSpellSlots(ctx, writeMonsters.InsertSpontaneousSpellSlotsParams{
				SpontaneousSpellCastingID: pgtype.Int4{Int32: spellCastingId, Valid: true},
				Level:                     pgtype.Text{String: monster.SpellCasting.SpontaneousSpellCasting[i].Slots[k].Level, Valid: true},
				Casts:                     pgtype.Text{String: monster.SpellCasting.SpontaneousSpellCasting[i].Slots[k].Casts, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("failed to assign spell slots in spontaneous block %w", err)
			}
		}
	}
	return nil
}

func ProcessMagic(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {
	if monster.SpellCasting.InnateSpellCasting != nil {
		err := ProcessInnateMagic(ctx, queries, monster, id)
		if err != nil {
			return fmt.Errorf("failed to parse an innate spellcasting Block %w", err)
		}
	}
	if monster.SpellCasting.FocusSpellCasting != nil {
		err := ProcessFocusMagic(ctx, queries, monster, id)
		if err != nil {
			return fmt.Errorf("failed to parse a focus spellcasting block %w", err)
		}
	}
	if monster.SpellCasting.PreparedSpellCasting != nil {
		err := ProcessPreparedMagic(ctx, queries, monster, id)
		if err != nil {
			return fmt.Errorf("failed to process a prepared spellcasting %w", err)
		}
	}
	if monster.SpellCasting.SpontaneousSpellCasting != nil {
		err := ProcessSpontaneousMagic(ctx, queries, monster, id)
		if err != nil {
			return fmt.Errorf("failed to process a spontaneous spellcasting block %w", err)
		}
	}
	return nil
}

func ProcessItems(ctx context.Context, queries *writeMonsters.Queries, monster structs.Monster, id int32) error {

	for i := range len(monster.Inventory) {
		dbID := uuid.New().String()
		itemId, err := queries.InsertItems(ctx, writeMonsters.InsertItemsParams{
			ID:          dbID,
			MonsterID:   pgtype.Int4{Int32: id, Valid: true},
			Name:        pgtype.Text{String: monster.Inventory[i].Name, Valid: true},
			Category:    pgtype.Text{String: monster.Inventory[i].Category, Valid: true},
			Description: pgtype.Text{String: monster.Inventory[i].Description, Valid: true},
			Level:       pgtype.Text{String: monster.Inventory[i].Level, Valid: true},
			Rarity:      pgtype.Text{String: monster.Inventory[i].Rarity, Valid: true},
			Bulk:        pgtype.Text{String: monster.Inventory[i].Bulk, Valid: true},
			Quantity:    pgtype.Text{String: monster.Inventory[i].Quantity, Valid: true},
			PricePer:    pgtype.Int4{Int32: int32(monster.Inventory[i].Price.Per), Valid: true},
			PriceCp:     pgtype.Int4{Int32: int32(monster.Inventory[i].Price.CP), Valid: true},
			PriceSp:     pgtype.Int4{Int32: int32(monster.Inventory[i].Price.SP), Valid: true},
			PriceGp:     pgtype.Int4{Int32: int32(monster.Inventory[i].Price.GP), Valid: true},
			PricePp:     pgtype.Int4{Int32: int32(monster.Inventory[i].Price.PP), Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to write item %s, %w", monster.Inventory[i].Name, err)
		}
		// write traits.
		for j := range len(monster.Inventory[i].Traits) {
			err := queries.InsertItemTraits(ctx, writeMonsters.InsertItemTraitsParams{
				ItemID: pgtype.Text{String: itemId},
				Trait:  pgtype.Text{String: monster.Inventory[i].Traits[j]},
			})
			if err != nil {
				return fmt.Errorf("failed to write item traits %w", err)
			}
		}
	}
	return nil
}

func WriteMonsterToDb(monster structs.Monster, cfg config.Config) error {
	logger.Log.Info(fmt.Sprintf("%+v", monster))
	ctx := context.Background()
	// ✅ 2. Begin a transaction
	tx, err := cfg.DBPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction %w", err)
	}

	//prep main params
	monsterParams := PrepMonsterParams(monster)

	queries := writeMonsters.New(cfg.DBPool)

	id, err := queries.InsertMonster(ctx, monsterParams)
	if err != nil {
		return fmt.Errorf("failed to insert Monster %w", err)
	}
	logger.Log.Info(fmt.Sprintf("Succesfully started the transaction for ID %d", id))
	//for each immunities
	err = writeImmunites(ctx, queries, monster, id)
	if err != nil {
		return fmt.Errorf("failed to write immunites to db %w", err)
	}
	err = ProcessWeakAndResist(ctx, queries, monster, id)
	if err != nil {
		return fmt.Errorf("failed to process weaknesses and resistances: %w", err)
	}
	err = ProcessLanguages(ctx, queries, monster, id)
	if err != nil {
		return fmt.Errorf("failed to write languages to db %w", err)
	}
	err = ProcessSenses(ctx, queries, monster, id)
	if err != nil {
		return fmt.Errorf("failed to write senses to db %w", err)
	}
	err = ProcessSkills(ctx, queries, monster, id)
	if err != nil {
		return fmt.Errorf("failed to process skills to db %w", err)
	}
	err = ProcessMovements(ctx, queries, monster, id)
	if err != nil {
		return fmt.Errorf("failed to process movements into db %w", err)
	}
	err = ProcessAction(ctx, queries, monster, id)
	if err != nil {
		return fmt.Errorf("failed to process action to db %w", err)
	}
	err = ProcessReaction(ctx, queries, monster, id)
	if err != nil {
		return fmt.Errorf("failed to process reaction to db %w", err)
	}
	err = ProcessPassive(ctx, queries, monster, id)
	if err != nil {
		return fmt.Errorf("failed to process free action to db %w", err)
	}
	err = ProcessAttacks(ctx, queries, monster, id)
	if err != nil {
		return fmt.Errorf("failed to process attack to db %w", err)
	}
	err = ProcessMagic(ctx, queries, monster, id)
	if err != nil {
		return fmt.Errorf("failed to process spellcasting blocks %w", err)
	}

	err = ProcessItems(ctx, queries, monster, id)
	if err != nil {
		return fmt.Errorf("failed to process items into db, %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction close %w", err)
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
		itemsList := gjson.Get(string(data), "items")
		// ItemsList := gjson.Get(string(data), "items").String()
		// fmt.Println(ItemsList)
		var spells []structs.Spell
		monster.FreeActions, monster.Actions, monster.Reactions, monster.Passives, monster.SpellCasting, spells, monster.Melees, monster.Ranged, monster.Inventory, err = ParseItems(itemsList)
		if err != nil {
			return err
		}
		AssignSpell(&spells, &monster.SpellCasting)

		err = WriteMonsterToDb(monster, cfg)
		if err != nil {
			return (err)
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
