package utils

import (
	"fmt"
	"os"

	"github.com/Burtcam/encounter-builder-backend/structs"
	"github.com/tidwall/gjson"
)

func ItemSwitch(item string, passiveList *[]structs.Passive, SpellCastingBlocks *structs.SpellCasting, FreeActionList *[]structs.FreeAction, ReactionList *[]structs.Reaction, actionList *[]structs.Action, SpellList *[]structs.Spell, MeleeList *[]structs.Attack, RangedList *[]structs.Attack, Inventory *[]structs.Item) error {
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
		//Safe to skip for now.
	case "weapon":
		fmt.Println("Found a weapon")
		*Inventory = append(*Inventory, ParseItem(item))
	case "armor":
		fmt.Println("Found an Armor")
		*Inventory = append(*Inventory, ParseItem(item))
	case "equipment":
		fmt.Println("Found equipment")
		*Inventory = append(*Inventory, ParseItem(item))
	case "consumable":
		fmt.Println("found a consumable")
		*Inventory = append(*Inventory, ParseItem(item))
	case "effect":
		fmt.Println("found an effect")
		//Can safely ignore because it's in the passives or actives
	case "treasure":
		fmt.Println("Found a Treasure")
		*Inventory = append(*Inventory, ParseItem(item))
	case "shield":
		fmt.Println("Shield")
		*Inventory = append(*Inventory, ParseItem(item))
	case "backpack":
		fmt.Println("backpack")
		*Inventory = append(*Inventory, ParseItem(item))
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
	return nil
}

func ParseItems(data gjson.Result) ([]structs.FreeAction, []structs.Action, []structs.Reaction, []structs.Passive, structs.SpellCasting, []structs.Spell, []structs.Attack, []structs.Attack, []structs.Item) {

	var passiveList []structs.Passive
	var SpellCastingBlocks structs.SpellCasting
	var FreeActionList []structs.FreeAction
	var actionList []structs.Action
	var ReactionList []structs.Reaction
	var SpellMasterList []structs.Spell
	var MeleeList []structs.Attack
	var RangedList []structs.Attack
	var inventory []structs.Item
	arrayJson := data.Array()
	for i := range len(arrayJson) {
		fmt.Println(arrayJson[i])
		ItemSwitch(arrayJson[i].String(),
			&passiveList,
			&SpellCastingBlocks,
			&FreeActionList,
			&ReactionList,
			&actionList,
			&SpellMasterList,
			&MeleeList,
			&RangedList,
			&inventory)

	}
	return FreeActionList,
		actionList,
		ReactionList,
		passiveList,
		SpellCastingBlocks,
		SpellMasterList,
		MeleeList,
		RangedList,
		inventory
}
