package utils

import (
	"fmt"

	"github.com/Burtcam/encounter-builder-backend/structs"
	"github.com/microcosm-cc/bluemonday"
	"github.com/tidwall/gjson"
)

func ParsePreparedSpellCasting(jsonData string) structs.PreparedSpellCasting {
	// Initialize the result struct
	var spellCasting structs.PreparedSpellCasting
	var preparedSlots []structs.PreparedSlot

	// Extract high-level attributes
	spellCasting.DC = int(gjson.Get(jsonData, "system.spelldc.dc").Int())
	spellCasting.Mod = gjson.Get(jsonData, "system.spelldc.value").String()
	spellCasting.Tradition = gjson.Get(jsonData, "system.tradition.value").String()
	spellCasting.ID = gjson.Get(jsonData, "_id").String()
	spellCasting.Description = gjson.Get(jsonData, "system.description.value").String()

	// Extract spell slots dynamically
	gjson.Get(jsonData, "system.slots").ForEach(func(slotName, slotData gjson.Result) bool {
		// Iterate over prepared spells in each slot
		slotData.Get("prepared").ForEach(func(_, spell gjson.Result) bool {
			fmt.Println("Slot Name:", slotName.String())
			fmt.Println("Spell ID:", spell.Get("id").String())
			slot := structs.PreparedSlot{
				Level:   slotName.String(), // Slot name represents the spell level (e.g., "slot0", "slot1")
				SpellID: spell.Get("id").String(),
			}
			preparedSlots = append(preparedSlots, slot)
			return true
		})
		return true
	})
	spellCasting.Slots = preparedSlots
	return spellCasting
}

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

func ParseDamageBlocks(jsonData string) []structs.DamageBlock {
	//loop over each item
	var DamageBlocks []structs.DamageBlock
	gjson.Get(jsonData, "system.damageRolls").ForEach(func(key, value gjson.Result) bool {
		fmt.Println(value.Get("damage").String())
		// key."damage"
		damageBlock := structs.DamageBlock{
			DamageRoll: value.Get("damage").String(),
			DamageType: value.Get("damageType").String(),
		}
		DamageBlocks = append(DamageBlocks, damageBlock)
		return true
	})
	return DamageBlocks
}

func ParseDamageEffects(jsonData string) structs.DamageEffect {
	fmt.Println(gjson.Get(jsonData, "system.attackEffects.value").String())
	effectBlock := structs.DamageEffect{
		CustomString: gjson.Get(jsonData, "system.attackEffects.custom").String(),
		Value:        ingestJSONList(jsonData, "system.attackEffects.value"),
	}
	return effectBlock
}
func stripHTMLUsingBluemonday(input string) string {
	p := bluemonday.StripTagsPolicy()
	return p.Sanitize(input)
}

func ParseWeapon(jsonData string) structs.Attack {
	var TypeDefinition string
	if !(gjson.Get(jsonData, "system.weaponType.value").Exists()) {
		TypeDefinition = "melee"
	} else {
		TypeDefinition = gjson.Get(jsonData, "system.weaponType.value").String()
	}
	damageBlocks := ParseDamageBlocks(jsonData)

	attackInScope := structs.Attack{
		Type:         TypeDefinition,
		ToHitBonus:   gjson.Get(jsonData, "system.bonus.value").String(),
		DamageBlocks: damageBlocks,
		Traits:       ingestJSONList(jsonData, "system.traits.value"),
		Effects:      ParseDamageEffects(jsonData),
		Name:         gjson.Get(jsonData, "name").String(),
	}
	return attackInScope
}

func ParseFreeAction(jsonData string) structs.FreeAction {
	fmt.Println("Found a free action")
	freeAction := structs.FreeAction{
		Name:     gjson.Get(jsonData, "name").String(),
		Text:     stripHTMLUsingBluemonday(gjson.Get(jsonData, "system.description.value").String()),
		Traits:   ingestJSONList(jsonData, "system.traits.value"),
		Category: gjson.Get(jsonData, "system.category").String(),
		Rarity:   gjson.Get(jsonData, "system.traits.rarity").String(),
	}
	return freeAction
}

// func CompareSpellCastingIDs(spellCasting structs.SpellCasting, spell structs.Spell, value gjson.Result) {
// 	// Check InnateSpellCasting

// 	locationID := value.Get("system").Get("location").Get("value").String()

// 	for _, item := range spellCasting.InnateSpellCasting {
// 		if item.ID == locationID {
// 			fmt.Printf("Match found in InnateSpellCasting: %s\n", item.ID)
// 		}
// 	}

// 	// Check PreparedSpellCasting
// 	for _, item := range spellCasting.PreparedSpellCasting {
// 		if item.ID == locationID {
// 			fmt.Printf("Match found in PreparedSpellCasting: %s\n", item.ID)
// 		}
// 	}

// 	// Check SpontaneousSpellCasting
// 	for _, item := range spellCasting.SpontaneousSpellCasting {
// 		if item.ID == locationID {
// 			fmt.Printf("Match found in SpontaneousSpellCasting: %s\n", item.ID)
// 		}
// 	}

// 	// Check FocusSpellCasting
// 	for _, item := range spellCasting.FocusSpellCasting {
// 		if item.ID == locationID {
// 			fmt.Printf("Match found in FocusSpellCasting: %s\n", item.ID)
// 		}
// 	}
// }

// func HandleSpell(spellCastingBlocks *structs.SpellCasting, value gjson.Result) (structs.SpellCasting, error) {
// 	// Check if the referenced spellcasting exists in the spellcasting blocks. If it does, add to the spell to the right spot,
// 	// ELSE create the spellcasting block, add it to the spellcasting THEN add the found spell.
// 	var Sustained bool
// 	var defenseSaveType bool
// 	if value.Get("system").Get("duration").Get("sustained").String() == "false" {
// 		Sustained = false
// 	} else {
// 		Sustained = true
// 	}
// 	if value.Get("system").Get("defense").Get("save").Get("basic") == "true" {
// 		defenseSaveType = true
// 	} else if value.Get("system").Get("defense").Get("save").Get("basic") == "false" {
// 		defenseSaveType = false
// 	} else {
// 		defenseSaveType = nil
// 	}
// 	Traits := extractListOfObjectsValues(value.String(), "system.traits.value")

// 	spell := structs.Spell{
// 		ID:          value.Get("_id").String(),
// 		Name:        value.Get("name").String(),
// 		Level:       value.Get("system").Get("level").Get("value").String(),
// 		Description: value.Get("system").Get("description").Get("value").String(),
// 		Range:       value.Get("system").Get("range").Get("value").String(),
// 		Area: structs.SpellArea{
// 			Type:   value.Get("system").Get("area").Get("type").String(),
// 			Value:  value.Get("system").Get("area").Get("value").String(),
// 			Detail: value.Get("system").Get("area").Get("details").String(),
// 		},
// 		Duration: structs.DurationBlock{
// 			Sustained: Sustained,
// 			Duration:  value.Get("system").Get("duration").Get("value").String(),
// 		},
// 		Targets: value.Get("system").Get("target").Get("value").String(),
// 		Traits:  Traits,
// 		Defense: structs.DefenseBlock{
// 			Save:  value.Get("system").Get("defense").Get("save").Get("statistic").String(),
// 			Basic: defenseSaveType,
// 		},
// 		CastTime:       value.Get("system").Get("time").Get("value").String(),
// 		CastComponents: value.Get("system").Get("cost").Get("value").String(),
// 		Rarity:         value.Get("system").Get("traits").Get("rarity").String(),
// 	}
// 	spellCastingLocation := value.Get("system").Get("location").Get("value")
// 	// find which spellcasting block it belongs in.
// }

// TODO Left off here
// func CreateSenseList(jsonData []byte, path string) []structs.Sense {
// 	var senses []structs.Sense

// 	// Get the JSON array stored in "senses"
// 	sensesData := gjson.Get(string(jsonData), path)
// 	// Iterate over each element in the senses array
// 	sensesData.ForEach(func(key, value gjson.Result) bool {
// 		senses = append(senses, structs.Sense{
// 			Name:   value.Get("type").String(),
// 			Range:  value.Get("range").String(),
// 			Detail: value.Get("details").String(),
// 			Acuity: value.Get("acuity").String(),
// 		})
// 		return true // Continue iterating
// 	})

// 	return senses
// }

func ExtractSkills(value gjson.Result) []structs.Skill {
	var skills []structs.Skill
	// Iterate over each key-value pair in the "skills" object.
	value.ForEach(func(key, value gjson.Result) bool {
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

func ingestJSONList(jsonData string, listString string) []string {
	result := gjson.Get(jsonData, listString).Array()

	// Convert []gjson.Result to []string
	var values []string
	for _, v := range result {
		values = append(values, v.String())
	}
	return values
}

func ParsePassives(value gjson.Result) structs.Passive {
	// Ensure the pointer list exists

	passive := structs.Passive{
		Name:     value.Get("name").String(),
		Text:     value.Get("system").Get("description").Get("value").String(),
		Traits:   ingestJSONList((value.String()), "system.traits.value"),
		Category: value.Get("system").Get("category").String(),
	}
	return passive
}
