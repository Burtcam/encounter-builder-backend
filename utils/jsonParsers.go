package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Burtcam/encounter-builder-backend/logger"
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
func ParsePassive(value string) structs.Passive {
	// Ensure the pointer list exists

	passive := structs.Passive{
		Name:     gjson.Get(value, "name").String(),
		Text:     gjson.Get(value, "system.description.value").String(),
		Traits:   ingestJSONList(value, "system.traits.value"),
		Category: gjson.Get(value, "system.category").String(),
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
	for j := range len(*spellList) {
		for p := range len(castingBlocks.PreparedSpellCasting) {
			//Loop over and check if spell[]
			if (*spellList)[j].SpellCastingBlockLocationID == castingBlocks.PreparedSpellCasting[p].ID {
				for i := range len(castingBlocks.PreparedSpellCasting[p].Slots) {
					if castingBlocks.PreparedSpellCasting[p].Slots[i].SpellID == (*spellList)[j].ID {
						castingBlocks.PreparedSpellCasting[p].Slots[i].Spell = (*spellList)[j]
						break
					}
				}
			}
		}
		for p := range len(castingBlocks.InnateSpellCasting) {
			if (*spellList)[j].SpellCastingBlockLocationID == castingBlocks.InnateSpellCasting[p].ID {
				// create a spell use using spell.
				// Convert CastLevel from string to int
				levelInt := 0
				if l, err := strconv.Atoi((*spellList)[j].CastLevel); err == nil {
					levelInt = l
				}
				castingBlocks.InnateSpellCasting[p].SpellUses =
					append(castingBlocks.InnateSpellCasting[p].SpellUses, structs.SpellUse{
						Spell: (*spellList)[j],
						Level: levelInt,
						Uses:  (*spellList)[j].Uses,
					})
				break
			}
		}
		for p := range len(castingBlocks.SpontaneousSpellCasting) {
			if (*spellList)[j].SpellCastingBlockLocationID == castingBlocks.SpontaneousSpellCasting[p].ID {
				castingBlocks.SpontaneousSpellCasting[p].SpellList =
					append(castingBlocks.SpontaneousSpellCasting[p].SpellList, (*spellList)[j])
				break
			}
		}
		for p := range len(castingBlocks.FocusSpellCasting) {
			if (*spellList)[j].SpellCastingBlockLocationID == castingBlocks.FocusSpellCasting[p].ID {
				castingBlocks.FocusSpellCasting[p].FocusSpellList =
					append(castingBlocks.FocusSpellCasting[p].FocusSpellList, (*spellList)[j])
				break
			}
		}
	}
}

func LoadJSON(path string) (string, error) {
	fmt.Println("Path is :", path)
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Log.Error(err.Error())
		return string(data), err
	}
	return string(data), err
}

func ParseSenses(jsonData string) []structs.Sense {
	var SenseList []structs.Sense
	senseData := gjson.Get(jsonData, "system.perception.senses").Array()
	for i := range len(senseData) {
		sense := structs.Sense{
			Name:   senseData[i].Get("type").String(),
			Acuity: senseData[i].Get("acuity").String(),
			Range:  senseData[i].Get("range").String(),
			Detail: senseData[i].Get("detail").String(),
		}
		SenseList = append(SenseList, sense)
	}
	return SenseList
}

func ParseSaves(jsonData string) structs.Saves {
	save := structs.Saves{
		Fort:       gjson.Get(jsonData, "system.saves.fortitude.value").String(),
		FortDetail: gjson.Get(jsonData, "system.saves.fortitude.saveDetail").String(),
		Will:       gjson.Get(jsonData, "system.saves.will.value").String(),
		WillDetail: gjson.Get(jsonData, "system.saves.will.saveDetail").String(),
		Ref:        gjson.Get(jsonData, "system.saves.reflex.value").String(),
		RefDetail:  gjson.Get(jsonData, "system.saves.reflex.saveDetail").String(),
	}

	return save
}

func ParsePrice(jsonData string) structs.PriceBlock {
	price := structs.PriceBlock{
		GP:  int(gjson.Get(jsonData, "system.price.value.gp").Int()),
		SP:  int(gjson.Get(jsonData, "system.price.value.sp").Int()),
		CP:  int(gjson.Get(jsonData, "system.price.value.cp").Int()),
		PP:  int(gjson.Get(jsonData, "system.price.value.pp").Int()),
		Per: int(gjson.Get(jsonData, "system.price.per").Int()),
	}
	return price
}

func ParseItem(jsonData string) structs.Item {

	item := structs.Item{
		Name:        gjson.Get(jsonData, "name").String(),
		ID:          gjson.Get(jsonData, "_id").String(),
		Category:    gjson.Get(jsonData, "system.category").String(),
		Description: stripHTMLUsingBluemonday(gjson.Get(jsonData, "system.description.value").String()),
		Level:       gjson.Get(jsonData, "system.level.value").String(),
		Price:       ParsePrice(jsonData),
		Type:        gjson.Get(jsonData, "type").String(),
		Traits:      ingestJSONList(jsonData, "system.traits.value"),
		Rarity:      gjson.Get(jsonData, "system.traits.rarity").String(),
		Range:       gjson.Get(jsonData, "system.range").String(),
		Size:        gjson.Get(jsonData, "system.size").String(),
		Reload:      gjson.Get(jsonData, "system.reload.value").String(),
		Bulk:        gjson.Get(jsonData, "system.bulk.value").String(),
		Quantity:    gjson.Get(jsonData, "system.quantity").String(),
	}
	return item
}
