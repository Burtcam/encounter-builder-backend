package utils

import (
	"fmt"
	"log"

	"os"

	"github.com/Burtcam/encounter-builder-backend/structs"
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
			fmt.Println(gjson.Get(string(data), "system.actionType.value").String())
			os.Exit(1)
			switch gjson.Get(string(data), "system.actionType.value").String() {
			case "passive":
				fmt.Println("Found a passive")
			case "action":
				fmt.Println("found a subaction")
			case "free":
				fmt.Println("Found a free action")
			default:
				fmt.Println("Uncategorized! : ", value.Get("type"))
				f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close() // Ensure the file is closed when we're done.
				logLine := fmt.Sprintf("Uncategorized Item Found, %s in the file %s \n", value.Get("type"), gjson.Get(string(data), "name"))
				// Write a string to the file.
				if _, err := f.WriteString(logLine); err != nil {
					log.Fatal(err)
				}

			}
			//split into passives, free, action ("system.actionType.Value" == passive)
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
		case "weapon":
			fmt.Println("Found a weapon")
		case "armor":
			fmt.Println("Found an Armor")
		case "equipment":
			fmt.Println("Found equipment")
		case "consumable":
			fmt.Println("found a consumable")
		case "effect":
			fmt.Println("found an effect")
		case "treasure":
			fmt.Println("Found a Treasure")
		case "shield":
			fmt.Println("Shield")
		case "backpack":
			fmt.Println("backpack")
		case "condition":
			fmt.Println("condition")
		default:
			fmt.Println("Uncategorized! : ", value.Get("type"))
			f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close() // Ensure the file is closed when we're done.
			logLine := fmt.Sprintf("Uncategorized Item Found, %s in the file %s \n", value.Get("type"), gjson.Get(string(data), "name"))
			// Write a string to the file.
			if _, err := f.WriteString(logLine); err != nil {
				log.Fatal(err)
			}
		}
		return true // Continue iterating
	})
	return nil
}

func ingestJSONList(jsonData []byte, listString string) []string {
	result := gjson.Get(string(jsonData), listString).Array()

	// Convert []gjson.Result to []string
	var values []string
	for _, v := range result {
		values = append(values, v.String())
	}
	return values
}

// written for immunities, but can be used for any list of objects
func extractListOfObjectsValues(jsonData string, path string) []string {
	var types []string

	// Get the JSON array stored in "immunities"
	immunities := gjson.Get(jsonData, path)
	// Iterate over each element in the immunities array
	immunities.ForEach(func(key, value gjson.Result) bool {
		// For each object, get the "type" field
		typ := value.Get("type").String()
		types = append(types, typ)
		return true // Continue iterating
	})

	return types
}

// TODO Left off here
func CreateSenseList(jsonData []byte, path string) []structs.Sense {
	var senses []structs.Sense

	// Get the JSON array stored in "senses"
	sensesData := gjson.Get(string(jsonData), path)
	// Iterate over each element in the senses array
	sensesData.ForEach(func(key, value gjson.Result) bool {
		senses = append(senses, structs.Sense{
			Name:   value.Get("type").String(),
			Range:  value.Get("range").String(),
			Detail: value.Get("details").String(),
			Acuity: value.Get("acuity").String(),
		})
		return true // Continue iterating
	})

	return senses
}

func ExtractSkills(jsonData string) []structs.Skill {
	var skills []structs.Skill

	// Get the "skills" object from the JSON.
	skillsJSON := gjson.Get(jsonData, "system.skills")

	// Iterate over each key-value pair in the "skills" object.
	skillsJSON.ForEach(func(key, value gjson.Result) bool {
		// key.String() is the skill name.
		// value is an object containing "base" and optionally "special".
		baseValue := int(value.Get("base").Int())
		var specials []structs.SkillSpecial

		// Check if a "special" field exists.
		specialArray := value.Get("special")
		if specialArray.Exists() {
			// Iterate over each special item.
			specialArray.ForEach(func(_, specialItem gjson.Result) bool {
				specValue := int(specialItem.Get("base").Int())
				specLabel := specialItem.Get("label").String()

				// Extract "predicate" array.
				var predicates []string
				predicateArray := specialItem.Get("predicate")
				predicateArray.ForEach(func(_, pred gjson.Result) bool {
					predicates = append(predicates, pred.String())
					return true // continue iteration
				})

				// Create a SkillSpecial instance and add to the slice.
				specials = append(specials, structs.SkillSpecial{
					Value:      specValue,
					Label:      specLabel,
					Predicates: predicates,
				})
				return true // continue iteration
			})
		}

		// Append the skill to the final slice.
		skills = append(skills, structs.Skill{
			Name:     key.String(),
			Value:    baseValue,
			Specials: specials,
		})
		return true // continue iterating over skills
	})

	return skills
}

func parseJSON(data []byte) error {
	fmt.Println(gjson.Get(string(data), "name"))
	fmt.Println(gjson.Get(string(data), "system.abilities.cha.mod"))
	// dont delete
	// dbObj := structs.Monster{
	// 	Name: gjson.Get(string(data), "name").String(),
	// 	Traits: structs.Traits{
	// 		Rarity:    gjson.Get(string(data), "system.traits.rarity").String(),
	// 		Size:      gjson.Get(string(data), "system.traits.size.value").String(),
	// 		TraitList: ingestJSONList(data, "system.traits.value"),
	// 	},
	// 	Attributes: structs.Attributes{
	// 		Str: gjson.Get(string(data), "system.abilities.str.mod").String(),
	// 		Dex: gjson.Get(string(data), "system.abilities.dex.mod").String(),
	// 		Con: gjson.Get(string(data), "system.abilities.con.mod").String(),
	// 		Wis: gjson.Get(string(data), "system.abilities.wis.mod").String(),
	// 		Int: gjson.Get(string(data), "system.abilities.int.mod").String(),
	// 		Cha: gjson.Get(string(data), "system.abilities.cha.mod").String(),
	// 	},
	// 	Level: gjson.Get(string(data), "system.details.level.value").String(),
	// 	AClass: structs.AC{
	// 		Value:  gjson.Get(string(data), "system.attributes.ac.value").String(),
	// 		Detail: gjson.Get(string(data), "system.attributes.ac.details.value").String(),
	// 	},
	// 	HP: structs.HP{
	// 		Detail: gjson.Get(string(data), "system.attributes.hp.details.value").String(),
	// 		Value:  int(gjson.Get(string(data), "system.attributes.hp.value").Int()),
	// 	},
	// 	Immunities:  extractListOfObjectsValues(string(data), "system.attributes.immunities"),
	// 	Weaknesses:  extractListOfObjectsValues(string(data), "system.attributes.weaknesses"),
	// 	Resistances: extractListOfObjectsValues(string(data), "system.attributes.resistances"),
	// 	Languages:   ingestJSONList(data, "system.details.languages.value"),
	// 	Senses:      CreateSenseList(data, "system.perception.senses"),
	// 	Perception: structs.Perception{
	// 		Mod:    gjson.Get(string(data), "system.perception.mod").String(),
	// 		Detail: gjson.Get(string(data), "system.perception.details").String(),
	// 	},
	// 	Skills: ExtractSkills(string(data)),
	// }
	ParseItems(data)
	return nil
}
