package utils

import (
	"fmt"
	"log"

	"os"

	"github.com/Burtcam/encounter-builder-backend/structs"
	"github.com/tidwall/gjson"
)

//	func sortSpells(spellCastingList structs.SpellCasting, MasterSpellList []structs.Spell) error {
//		// loop over thhe master spell list, add it to the proper spellcastingList
//		MasterSpellList.ForEach(func(key, value) bool{
//			// match id in each spell to ID in spellcastingList
//			// For each spell, sort it appropriately based on each type.
//			// if found type is Prepared: create a PreparedSlot, and put the spell in it, and then put the Prepared Slot in the spellcasting struct
//			// If found type is Spontaenous put spell in SpellList
//			//If found type is Innate create a spellUse (with spell in it) and put in innate.SpellUses
//			// if found type is Focus put in focus spell list
//		}
//		return nil
//	}
//
// ParseSlots extracts slot data from JSON and returns a slice of Slot structs.
// ParsePreparedSpellCasting extracts spell slots from JSON into PreparedSpellCasting struct
func ParsePreparedSpellCasting(jsonData string) structs.PreparedSpellCasting {
	// Initialize the result struct
	var spellCasting structs.PreparedSpellCasting

	// Extract high-level attributes
	spellCasting.DC = int(gjson.Get(jsonData, "system.spelldc.dc").Int())
	spellCasting.Mod = gjson.Get(jsonData, "system.spelldc.value").String()
	spellCasting.Tradition = gjson.Get(jsonData, "system.tradition.value").String()
	spellCasting.ID = gjson.Get(jsonData, "_id").String()
	spellCasting.Description = gjson.Get(jsonData, "system.description.value").String()

	// Extract spell slots dynamically
	gjson.Get(jsonData, "slots").ForEach(func(slotName, slotData gjson.Result) bool {
		// Iterate over prepared spells in each slot
		slotData.Get("prepared").ForEach(func(_, spell gjson.Result) bool {
			spellCasting.Slots = append(spellCasting.Slots, structs.PreparedSlot{
				Level:   slotName.String(), // Slot name represents the spell level (e.g., "slot0", "slot1")
				SpellID: spell.Get("id").String(),
			})
			return true
		})
		return true
	})

	return spellCasting
}

// ParseSpontaneousSpellCasting extracts relevant data into SpontaneousSpellCasting struct
func ParseSpontaneousSpellCasting(jsonData string) structs.SpontaneousSpellCasting {
	var entry structs.SpontaneousSpellCasting
	// Extract top-level attributes
	entry.ID = gjson.Get(jsonData, "_id").String()
	entry.Tradition = gjson.Get(jsonData, "system.tradition.value").String()
	entry.DC = int(gjson.Get(jsonData, "system.spelldc.dc").Int())
	entry.Mod = gjson.Get(jsonData, "system.spelldc.value").String()

	// Extract slot data dynamically
	gjson.Get(jsonData, "system.slots").ForEach(func(slotName, slotData gjson.Result) bool {
		entry.Slots = append(entry.Slots, structs.Slot{
			Level: slotName.String(),
			Casts: slotData.Get("value").String(), // Converting max casts to string
		})
		return true
	})

	return entry
}

func ParseDamageBlocks(jsonData gjson.Result) []structs.DamageBlock {
	// for each item in DamageRolls, create a DamageBlock and add it to the slice
	var DamageBlocks []structs.DamageBlock
	// Get the JSON array stored in "damageRolls"
	damageRolls := jsonData.Get("system").Get("damageRolls").String()
	// Iterate over each element in the damageRolls array
	damageRolls.ForEach(func(key, value gjson.Result) bool {
		damageBlock := structs.DamageBlock{
			DamageRoll: value.Get(key.String()).Get("damage").String(),
			DamageType: value.Get(key.String()).Get("damageType").String(),
		}
		DamageBlocks = append(DamageBlocks, damageBlock)
		return true
	})
	return DamageBlocks
}

func ParseWeapon(value gjson.Result) (structs.Attack, error) {

	if !(value.Get("system").Get("weaponType").Get("value").Exists()) {
		TypeDefinition := "melee"
	} else {
		TypeDefinition := value.Get("system").Get("weaponType").Get("value").String()
	}
	damageBlocks, err := ParseDamageBlocks(jsonData)

	attackInScope := structs.Attack{
		Type:         TypeDefinition,
		ToHitBonus:   gjson.Get("system").Get("bonus").Get("value"),
		DamageBlocks: damageBlocks,
		Traits:       ingestJSONList(jsonData, "system.traits.value"),
		Effects:      ingestJSONList(jsonData, "system.effects.value"),
	}
	return attackInScope, nil
}

// CompareSpellCastingIDs loops over each slice in SpellCasting and checks if
// any item's ID matches the provided spellCastingLocation.
func CompareSpellCastingIDs(spellCasting structs.SpellCasting, spell structs.Spell, value gjson.Result) {
	// Check InnateSpellCasting

	locationID := value.Get("system").Get("location").Get("value").String()

	for _, item := range spellCasting.InnateSpellCasting {
		if item.ID == locationID {
			fmt.Printf("Match found in InnateSpellCasting: %s\n", item.ID)
		}
	}

	// Check PreparedSpellCasting
	for _, item := range spellCasting.PreparedSpellCasting {
		if item.ID == locationID {
			fmt.Printf("Match found in PreparedSpellCasting: %s\n", item.ID)
		}
	}

	// Check SpontaneousSpellCasting
	for _, item := range spellCasting.SpontaneousSpellCasting {
		if item.ID == locationID {
			fmt.Printf("Match found in SpontaneousSpellCasting: %s\n", item.ID)
		}
	}

	// Check FocusSpellCasting
	for _, item := range spellCasting.FocusSpellCasting {
		if item.ID == locationID {
			fmt.Printf("Match found in FocusSpellCasting: %s\n", item.ID)
		}
	}
}

func HandleSpell(spellCastingBlocks *structs.SpellCasting, value gjson.Result) (structs.SpellCasting, error) {
	// Check if the referenced spellcasting exists in the spellcasting blocks. If it does, add to the spell to the right spot,
	// ELSE create the spellcasting block, add it to the spellcasting THEN add the found spell.
	var Sustained bool
	var defenseSaveType bool
	if value.Get("system").Get("duration").Get("sustained").String() == "false" {
		Sustained = false
	} else {
		Sustained = true
	}
	if value.Get("system").Get("defense").Get("save").Get("basic") == "true" {
		defenseSaveType = true
	} else if value.Get("system").Get("defense").Get("save").Get("basic") == "false" {
		defenseSaveType = false
	} else {
		defenseSaveType = nil
	}
	Traits := extractListOfObjectsValues(value.String(), "system.traits.value")

	spell := structs.Spell{
		ID:          value.Get("_id").String(),
		Name:        value.Get("name").String(),
		Level:       value.Get("system").Get("level").Get("value").String(),
		Description: value.Get("system").Get("description").Get("value").String(),
		Range:       value.Get("system").Get("range").Get("value").String(),
		Area: structs.SpellArea{
			Type:   value.Get("system").Get("area").Get("type").String(),
			Value:  value.Get("system").Get("area").Get("value").String(),
			Detail: value.Get("system").Get("area").Get("details").String(),
		},
		Duration: structs.DurationBlock{
			Sustained: Sustained,
			Duration:  value.Get("system").Get("duration").Get("value").String(),
		},
		Targets: value.Get("system").Get("target").Get("value").String(),
		Traits:  Traits,
		Defense: structs.DefenseBlock{
			Save:  value.Get("system").Get("defense").Get("save").Get("statistic").String(),
			Basic: defenseSaveType,
		},
		CastTime:       value.Get("system").Get("time").Get("value").String(),
		CastComponents: value.Get("system").Get("cost").Get("value").String(),
		Rarity:         value.Get("system").Get("traits").Get("rarity").String(),
	}
	spellCastingLocation := value.Get("system").Get("location").Get("value")
	// find which spellcasting block it belongs in.
}

// []structs.Action, []structs.FreeAction, []structs.Attack, []structs.Attack, []structs.Reaction, []structs.Passive, []structs.SpellCasting,
func ParseItems(data []byte) ([]structs.FreeAction,
	[]structs.Action,
	[]structs.Reaction,
	[]structs.Passive,
	structs.SpellCasting,
	error) {
	rawItems := gjson.Get(string(data), "items")
	// for each item in the block, "items.system"
	// check systems.item.type and the three types (action, melee, ranged, spell, spellcastingEntry)
	// if it's a spell, ignore for now but create a map to look it up later?
	// if it's a action, encode it into the proper struct (passive, reaction, free, action)
	// for each, make a slice of type each (so []Action, []Reaction, []Passive, []FreeAction)
	// return those slices.
	var passiveList []structs.Passive
	var SpellCastingBlocks structs.SpellCasting
	var FreeActionList []structs.FreeAction
	var actionList []structs.Action
	var ReactionList []structs.Reaction // If it's a spell, create the spell and put it in master spell list (with the location.id) as well as the uses
	// If it's a spellcasting, create the spellcasting entry and attach it to the SpellCasting struct
	rawItems.ForEach(func(key, value gjson.Result) bool {
		switch value.Get("type").String() {
		case "action":
			fmt.Println("Found an Action")
			switch value.Get("system").Get("actionType").Get("value").String() {
			case "passive":
				fmt.Println("Found a passive")
				passive := structs.Passive{
					Name:     value.Get("name").String(),
					Text:     value.Get("system").Get("description").Get("value").String(),
					Traits:   extractListOfObjectsValues(string(data), "system.traits.value"),
					Category: value.Get("system").Get("category").String(),
				}
				passiveList = append(passiveList, passive)
			case "action":
				fmt.Println("found a subaction")
				action := structs.Action{
					Name:     value.Get("name").String(),
					Text:     value.Get("system").Get("description").Get("value").String(),
					Traits:   ingestJSONList(data, "system.traits.value"),
					Category: value.Get("system").Get("category").String(),
					Actions:  value.Get("system").Get("actions").Get("value").String(),
					Rarity:   value.Get("system").Get("traits").Get("rarity").String(),
				}
				actionList = append(actionList, action)
			case "free":
				fmt.Println("Found a free action")
				freeAction := structs.FreeAction{
					Name:     value.Get("name").String(),
					Text:     value.Get("system").Get("description").Get("value").String(),
					Traits:   extractListOfObjectsValues(string(data), "system.traits.value"),
					Category: value.Get("system").Get("category").String(),
					Rarity:   value.Get("system").Get("traits").Get("rarity").String(),
				}
				FreeActionList = append(FreeActionList, freeAction)

			case "reaction":
				fmt.Println("Found a reaction")
				reaction := structs.Reaction{
					Name:     value.Get("name").String(),
					Text:     value.Get("system").Get("description").Get("value").String(),
					Traits:   ingestJSONList(data, "system.traits.value"),
					Category: value.Get("system").Get("category").String(),
					Rarity:   value.Get("system").Get("traits").Get("rarity").String(),
				}
				ReactionList = append(ReactionList, reaction)
			default:
				fmt.Println("Uncategorized! : ", value.Get("system").Get("actionType").Get("value"))
				f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close() // Ensure the file is closed when we're done.
				logLine := fmt.Sprintf("Uncategorized Item Found, %s in the file %s \n", value.Get("system").Get("actionType").Get("value"), gjson.Get(string(data), "name"))
				// Write a string to the file.
				if _, err := f.WriteString(logLine); err != nil {
					log.Fatal(err)
				}

			}
			//split into passives, free, action ("system.actionType.Value" == passive)
		case "spellcastingEntry":
			fmt.Println("Found a spellcasting entry")
			if value.Get("system").Get("prepared").Get("value").String() == "innate" {
				fmt.Println("Innate Handler Called")
				innate := structs.InnateSpellCasting{
					ID:          value.Get("_id").String(),
					Tradition:   value.Get("system").Get("tradition").Get("value").String(),
					DC:          int(value.Get("system").Get("spelldc").Get("dc").Int()),
					Mod:         value.Get("system").Get("spelldc").Get("value").String(),
					Description: value.Get("system").Get("description").Get("value").String(),
				}
				SpellCastingBlocks.InnateSpellCasting = append(SpellCastingBlocks.InnateSpellCasting, innate)
			} else if value.Get("system").Get("prepared").Get("value").String() == "prepared" {
				fmt.Println("Innate Handler Called")
				SpellCastingBlocks.PreparedSpellCasting = append(SpellCastingBlocks.PreparedSpellCasting, ParsePreparedSpellCasting(value.String()))
				// Create a spellcastingEntry of the type,
			} else if value.Get("system").Get("prepared").Get("value").String() == "spontaneous" {
				fmt.Println("spontaneous Handler Called")
				SpellCastingBlocks.SpontaneousSpellCasting = append(SpellCastingBlocks.SpontaneousSpellCasting, ParseSpontaneousSpellCasting(value.String()))
				// Create a spellcastingEntry of the type,
			} else if value.Get("system").Get("prepared").Get("value").String() == "focus" {
				fmt.Println("focus Handler Called")
				// Create a spellcastingEntry of the type,
				focus := structs.FocusSpellCasting{
					ID:          value.Get("_id").String(),
					Tradition:   value.Get("system").Get("tradition").Get("value").String(),
					DC:          int(value.Get("system").Get("spelldc").Get("dc").Int()),
					Mod:         value.Get("system").Get("spelldc").Get("value").String(),
					Description: value.Get("system").Get("description").Get("value").String(),
				}
				SpellCastingBlocks.FocusSpellCasting = append(SpellCastingBlocks.FocusSpellCasting, focus)

			} else if value.Get("system").Get("prepared").Get("value").String() == "items" {
				fmt.Println("Spell Item Handler Called")
			} else {
				fmt.Println("Uncategorized! : ", value.Get("type"))
				f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close() // Ensure the file is closed when we're done.
				logLine := fmt.Sprintf("Uncategorized Spellcasting entry Found, %s in the file %s \n", value.Get("system").Get("prepared").Get("value").String(), gjson.Get(string(data), "name"))
				// Write a string to the file.
				if _, err := f.WriteString(logLine); err != nil {
					log.Fatal(err)
				}

			}

		// return a []Attack each one has a system.weaponType which will be ranged or melee
		case "melee":
			fmt.Println("Found a melee")
			weapon, err := ParseWeapon(value.String())
			if err != nil {
				fmt.Println("Some other shit!")
				fmt.Println("Uncategorized! : ", value.Get("type"))
				f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close() // Ensure the file is closed when we're done.
				logLine := fmt.Sprintf("Uncategorized attack Found, %s in the file %s \n", value.Get("system").Get("weaponType").Get("value").String(), gjson.Get(string(data), "name"))
				// Write a string to the file.
				if _, err := f.WriteString(logLine); err != nil {
					log.Fatal(err)
				}
			}

		case "spell":
			fmt.Println("found a Spell")
			// Check if the referenced spellcasting struct exists yet, if it does, add it to that, (location.value) and (location.uses)
			SpellCastingBlocks, err := HandleSpell(&SpellCastingBlocks, value)
			if err != nil {
				fmt.Println("Failed to handle spell")
			}
			// If it hasn't THEN loop over the item list till we find it and create it.
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
	return FreeActionList, actionList, ReactionList, passiveList, SpellCastingBlocks, nil
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
