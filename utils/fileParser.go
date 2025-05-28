package utils

import (
	"fmt"
	"os"

	"github.com/Burtcam/encounter-builder-backend/structs"
	"github.com/tidwall/gjson"
)

func StandinFunc() bool {
	return true
}

func ItemSwitch(item string,
	passiveList *[]structs.Passive,
	SpellCastingBlocks *structs.SpellCasting,
	FreeActionList *[]structs.FreeAction,
	ReactionList *[]structs.Reaction,
	actionList *[]structs.Action,
	SpellList *[]structs.Spell,
	MeleeList *[]structs.Attack,
	RangedList *[]structs.Attack) error {
	switch gjson.Get(item, "type").String() {
	case "action":
		switch gjson.Get(item, "system.actionType.value").String() {
		case "action":
			*actionList = append(*actionList, ParseAction(item))
		case "passive":
			*passiveList = append(*passiveList, ParsePassive(item))
		case "free":
			*FreeActionList = append(*FreeActionList, ParseFreeAction(item))
		case "reaction":
			*ReactionList = append(*ReactionList, ParseReaction(item))
		default:
			fmt.Println("Uncategorized!: ", gjson.Get(item, "system.actionType.value"))
			f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer f.Close() // Ensure the file is closed when we're done.
		}
	case "spellcastingentry":
		switch gjson.Get(item, "system.prepared.value").String() {
		case "prepared":
			SpellCastingBlocks.PreparedSpellCasting = append(SpellCastingBlocks.PreparedSpellCasting, ParsePreparedSpellCasting(item))
		case "spontaneous":
			SpellCastingBlocks.SpontaneousSpellCasting = append(SpellCastingBlocks.SpontaneousSpellCasting, ParseSpontaneousSpellCasting(item))
		case "focus":
			SpellCastingBlocks.FocusSpellCasting = append(SpellCastingBlocks.FocusSpellCasting, ParseFocusSpellCasting(item))
		case "innate":
			SpellCastingBlocks.InnateSpellCasting = append(SpellCastingBlocks.InnateSpellCasting, ParseInnateSpellCasting(item))
		}
	case "spell":
		*SpellList = append(*SpellList, ParseSpell(item))
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
		//Can safely ignore because it's in the passives
	case "melee":
		switch gjson.Get(item, "system.weaponType.value").String() {
		case "melee":
			*MeleeList = append(*MeleeList, ParseWeapon(item))
		case "ranged":
			*RangedList = append(*RangedList, ParseWeapon(item))
		default:
			*MeleeList = append(*MeleeList, ParseWeapon(item))
		} //switch on the different types and call the ingesters

	}

}

func ParseItems(data string) ([]structs.FreeAction,
	[]structs.Action,
	[]structs.Reaction,
	[]structs.Passive,
	structs.SpellCasting,
	[]structs.Spell,
	[]structs.Attack,
	[]structs.Attack) {

	var passiveList []structs.Passive
	var SpellCastingBlocks structs.SpellCasting
	var FreeActionList []structs.FreeAction
	var actionList []structs.Action
	var ReactionList []structs.Reaction
	var SpellMasterList []structs.Spell
	var MeleeList []structs.Attack
	var RangedList []structs.Attack

	itemsList := gjson.Get(data, "items").Array()

	for i := 0; i < len(itemsList); i++ {
		ItemSwitch(itemsList[i].String(),
			&passiveList,
			&SpellCastingBlocks,
			&FreeActionList,
			&ReactionList,
			&actionList,
			&SpellMasterList,
			&MeleeList,
			&RangedList)
	}
	return FreeActionList,
		actionList,
		ReactionList,
		passiveList,
		SpellCastingBlocks,
		SpellMasterList,
		MeleeList,
		RangedList
}

// func parseJSON(data []byte) error {
// 	fmt.Println(gjson.Get(string(data), "name"))
// 	fmt.Println(gjson.Get(string(data), "system.abilities.cha.mod"))
// 	// dont delete
// 	// dbObj := structs.Monster{
// 	// 	Name: gjson.Get(string(data), "name").String(),
// 	// 	Traits: structs.Traits{
// 	// 		Rarity:    gjson.Get(string(data), "system.traits.rarity").String(),
// 	// 		Size:      gjson.Get(string(data), "system.traits.size.value").String(),
// 	// 		TraitList: ingestJSONList(data, "system.traits.value"),
// 	// 	},
// 	// 	Attributes: structs.Attributes{
// 	// 		Str: gjson.Get(string(data), "system.abilities.str.mod").String(),
// 	// 		Dex: gjson.Get(string(data), "system.abilities.dex.mod").String(),
// 	// 		Con: gjson.Get(string(data), "system.abilities.con.mod").String(),
// 	// 		Wis: gjson.Get(string(data), "system.abilities.wis.mod").String(),
// 	// 		Int: gjson.Get(string(data), "system.abilities.int.mod").String(),
// 	// 		Cha: gjson.Get(string(data), "system.abilities.cha.mod").String(),
// 	// 	},
// 	// 	Level: gjson.Get(string(data), "system.details.level.value").String(),
// 	// 	AClass: structs.AC{
// 	// 		Value:  gjson.Get(string(data), "system.attributes.ac.value").String(),
// 	// 		Detail: gjson.Get(string(data), "system.attributes.ac.details.value").String(),
// 	// 	},
// 	// 	HP: structs.HP{
// 	// 		Detail: gjson.Get(string(data), "system.attributes.hp.details.value").String(),
// 	// 		Value:  int(gjson.Get(string(data), "system.attributes.hp.value").Int()),
// 	// 	},
// 	// 	Immunities:  extractListOfObjectsValues(string(data), "system.attributes.immunities"),
// 	// 	Weaknesses:  extractListOfObjectsValues(string(data), "system.attributes.weaknesses"),
// 	// 	Resistances: extractListOfObjectsValues(string(data), "system.attributes.resistances"),
// 	// 	Languages:   ingestJSONList(data, "system.details.languages.value"),
// 	// 	Senses:      CreateSenseList(data, "system.perception.senses"),
// 	// 	Perception: structs.Perception{
// 	// 		Mod:    gjson.Get(string(data), "system.perception.mod").String(),
// 	// 		Detail: gjson.Get(string(data), "system.perception.details").String(),
// 	// 	},
// 	// 	Skills: ExtractSkills(string(data)),
// 	// }
// 	ParseItems(data)
// 	return nil
// }
