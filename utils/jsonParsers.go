package utils

import (
	"fmt"
	"strings"

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
func ParseFocusSpellCasting(jsonData string) structs.FocusSpellCasting {
	fmt.Println(gjson.Get(jsonData, "system.spelldc.value").String())
	entry := structs.FocusSpellCasting{
		DC:             int(gjson.Get(jsonData, "system.spelldc.dc").Int()),
		Mod:            gjson.Get(jsonData, "system.spelldc.value").String(),
		Tradition:      gjson.Get(jsonData, "system.tradition.value").String(),
		ID:             gjson.Get(jsonData, "_id").String(),
		FocusSpellList: []structs.Spell{},
		Name:           gjson.Get(jsonData, "name").String(),
		Description:    stripHTMLUsingBluemonday(gjson.Get(jsonData, "system.description.value").String()),
		CastLevel:      "",
	}
	return entry
}

func ParseInnateSpellCasting(jsonData string) structs.InnateSpellCasting {
	entry := structs.InnateSpellCasting{
		DC:          int(gjson.Get(jsonData, "system.spelldc.dc").Int()),
		Mod:         gjson.Get(jsonData, "system.spelldc.value").String(),
		Tradition:   gjson.Get(jsonData, "system.tradition.value").String(),
		ID:          gjson.Get(jsonData, "_id").String(),
		SpellUses:   []structs.SpellUse{},
		Name:        gjson.Get(jsonData, "name").String(),
		Description: stripHTMLUsingBluemonday(gjson.Get(jsonData, "system.description.value").String()),
	}
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
func ParseReaction(jsonData string) structs.Reaction {
	reaction := structs.Reaction{
		Name:     gjson.Get(jsonData, "name").String(),
		Text:     stripHTMLUsingBluemonday(gjson.Get(jsonData, "system.description.value").String()),
		Traits:   ingestJSONList(jsonData, "system.traits.value"),
		Category: gjson.Get(jsonData, "system.category").String(),
		Rarity:   gjson.Get(jsonData, "system.traits.rarity").String(),
	}
	return reaction
}
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
func ParseAction(jsonData string) structs.Action {
	action := structs.Action{
		Name:     gjson.Get(jsonData, ("name")).String(),
		Text:     stripHTMLUsingBluemonday(gjson.Get(jsonData, "system.description.value").String()),
		Traits:   ingestJSONList(jsonData, "system.traits.value"),
		Category: gjson.Get(jsonData, "system.category").String(),
		Actions:  gjson.Get(jsonData, "system.actions.value").String(),
		Rarity:   gjson.Get(jsonData, "system.traits.rarity").String(),
	}
	return action
}
func SpellLevelParser(jsonData string) string {
	if gjson.Get(jsonData, "system.location.heightenedLevel").String() != "" {
		return gjson.Get(jsonData, "system.location.heightenedLevel").String()
	} else {
		return gjson.Get(jsonData, "system.level.value").String()
	}
}
func ParseSpellArea(jsonData string) structs.SpellArea {
	spellArea := structs.SpellArea{
		Type:   gjson.Get(jsonData, "system.area.type").String(),
		Value:  gjson.Get(jsonData, "system.area.value").String(),
		Detail: gjson.Get(jsonData, "system.area.detail").String(),
	}
	return spellArea
}
func ParseDurationBlock(jsonData string) structs.DurationBlock {
	duration := structs.DurationBlock{
		Sustained: gjson.Get(jsonData, "system.duration.sustained").Bool(),
		Duration:  gjson.Get(jsonData, "system.duration.value").String(),
	}
	return duration
}
func ParseDefenseBlock(jsonData string) structs.DefenseBlock {
	defense := structs.DefenseBlock{
		Save:  gjson.Get(jsonData, "system.defense.save.statistic").String(),
		Basic: gjson.Get(jsonData, "system.defense.save.basic").Bool(),
	}
	return defense
}
func DetectRitual(jsonData string) (bool, structs.RitualData) {
	if gjson.Get(jsonData, "system.ritual").String() != "" {
		ritualBool := true
		ritualData := structs.RitualData{
			PrimaryCheck:     gjson.Get(jsonData, "system.ritual.primary.check").String(),
			SecondaryCasters: gjson.Get(jsonData, "system.ritual.secondary.casters").String(),
			SecondaryCheck:   gjson.Get(jsonData, "system.ritual.secondary.checks").String(),
		}
		return ritualBool, ritualData
	} else {
		return false, structs.RitualData{}
	}
}

// spell plan
// 1. Go get all spells and spellcasting blocks
// 2. For each spell, tie to spellcasting block
func ParseSpell(jsonData string) structs.Spell {
	ritualBool, ritualData := DetectRitual(jsonData)
	var uses string
	var AtWill bool
	Name := gjson.Get(jsonData, "name").String()
	if strings.Contains(Name, "(At Will)") {
		AtWill = true
		uses = "unlimited"
	} else {
		AtWill = false
	}
	if gjson.Get(jsonData, "system.location.uses.value").Exists() && AtWill == false {
		uses = gjson.Get(jsonData, "system.location.uses.value").String()
	} else if !(gjson.Get(jsonData, "system.location.uses.value").Exists()) && AtWill == false {
		uses = "1"
	}
	spell := structs.Spell{
		ID:                          gjson.Get(jsonData, "_id").String(),
		Name:                        gjson.Get(jsonData, "name").String(),
		CastLevel:                   SpellLevelParser(jsonData),
		SpellBaseLevel:              gjson.Get(jsonData, "system.level.value").String(),
		Description:                 stripHTMLUsingBluemonday(gjson.Get(jsonData, "system.description.value").String()),
		Range:                       gjson.Get(jsonData, "system.range.value").String(),
		Area:                        ParseSpellArea(jsonData),
		Duration:                    ParseDurationBlock(jsonData),
		Targets:                     gjson.Get(jsonData, "system.target.value").String(),
		Traits:                      ingestJSONList(jsonData, "system.traits.value"),
		Defense:                     ParseDefenseBlock(jsonData),
		CastTime:                    gjson.Get(jsonData, "system.time.value").String(),
		CastRequirements:            gjson.Get(jsonData, "system.requirements").String(),
		Rarity:                      gjson.Get(jsonData, "system.traits.rarity").String(),
		SpellCastingBlockLocationID: gjson.Get(jsonData, "system.location.value").String(),
		Uses:                        uses,
		Ritual:                      ritualBool,
		RitualData:                  ritualData,
		AtWill:                      AtWill,
	}
	return spell
}

func AssignSpell(spellList *[]structs.Spell, castingBlocks *structs.SpellCasting) {
	//For each spell in spell list. Look in each list of castingblocks for the match
	// castingblocks.PreparedSpellList[0].ID == spellList[i].SpellCastingBlockLocationID
	for j := 0; j < len(*spellList); j++ {
		for p := 0; p < len(castingBlocks.PreparedSpellCasting); p++ {
			//Loop over and check if spell[]
			if (*spellList)[j].SpellCastingBlockLocationID == castingBlocks.PreparedSpellCasting[p].ID {
				for i := 0; i < len(castingBlocks.PreparedSpellCasting[p].Slots); i++ {
					if castingBlocks.PreparedSpellCasting[p].Slots[i].SpellID == (*spellList)[j].ID {
						castingBlocks.PreparedSpellCasting[p].Slots[i].Spell = (*spellList)[j]
					}
				}
			}
		}
		for p := 0; p < len(castingBlocks.InnateSpellCasting); p++ {
			if (*spellList)[j].SpellCastingBlockLocationID == castingBlocks.InnateSpellCasting[p].ID {
				// create a spell use using spell. 
				castingBlocks.InnateSpellCasting[p].SpellUses = append(castingBlocks.InnateSpellCasting[p].SpellUses, structs.SpellUse{
					Spell: *(spellList[j]), 
					Level: int(*spellList)[j].CastLevel), 
					Uses: (*spellList)[j].Uses,
				}) 
			}
			

		}

	}
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
