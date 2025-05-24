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
			// weapon, err := ParseWeapon(value.String())
			// if err != nil {
			// 	fmt.Println("Some other shit!")
			// 	fmt.Println("Uncategorized! : ", value.Get("type"))
			// 	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	defer f.Close() // Ensure the file is closed when we're done.
			// 	logLine := fmt.Sprintf("Uncategorized attack Found, %s in the file %s \n", value.Get("system").Get("weaponType").Get("value").String(), gjson.Get(string(data), "name"))
			// 	// Write a string to the file.
			// 	if _, err := f.WriteString(logLine); err != nil {
			// 		log.Fatal(err)
			// 	}
			// }

		case "spell":
			fmt.Println("found a Spell")
			// // Check if the referenced spellcasting struct exists yet, if it does, add it to that, (location.value) and (location.uses)
			// SpellCastingBlocks, err := HandleSpell(&SpellCastingBlocks, value)
			// if err != nil {
			// 	fmt.Println("Failed to handle spell")
			// }
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
