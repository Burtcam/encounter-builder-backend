package main

import (
	"fmt"
	"os"

	"github.com/Burtcam/encounter-builder-backend/logger"
	"github.com/tidwall/gjson"
)

func ParseItems(data []byte) error {
	rawItems := gjson.Get(string(data), "items")
	// for each item in the block, "items.system"
	// check systems.item.type and the three types (action, melee, ranged, spell, spellcastingEntry)
	rawItems.ForEach(func(key, value gjson.Result) bool {
		switch value.Get("type").String() {
		case "action":
			fmt.Println("Found an Action")
		case "spellcastingEntry":
			fmt.Println("Found a spellcasting entry")
		case "melee":
			fmt.Println("Found a melee")
		case "ranged":
			fmt.Println("Found a Ranged")
		case "spell":
			fmt.Println("found a Spell")
		case "lore":
			fmt.Println("Found a Lore")
		default:
			fmt.Println("Uncategorized! : ", value.Get("type"))

		}
		return true // Continue iterating
	})
	return nil
}

func LoadAJson(path string) error {
	fmt.Println("Path is :", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	// var payload map[string]interface{}
	// err = json.Unmarshal(data, &payload)
	// if err != nil {
	// 	logger.Log.Error("Error during Unmarshal(): %s", path, err)
	// }
	if gjson.Get(string(data), "type").String() == "npc" {
		fmt.Println("Found a monster")
		err = ParseItems(data)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("Error Parsing file %s", path))
		}
		// WRite it out to a json
		// err = os.WriteFile("example-monster.json", jsonData, 0644)
		// if err != nil {
		// 	logger.Log.Error("Error writting JSON:", err)
		// }
		os.Exit(1)
	}

	return nil
}

// func main() {
// 	// Sample JSON containing the skills block
// 	// load a file
// 	err := LoadAJson("/home/cburt/encounter-builder/encounter-builder-backend/files/foundryvtt-pf2e-3bcc5cf/packs/pathfinder-monster-core/fortune-dragon-adult.json")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	// Extract the skills.
// }
