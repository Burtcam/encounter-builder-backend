package utils

import (
	"encoding/json"
	"fmt"

	"github.com/Burtcam/encounter-builder-backend/structs"
	"github.com/tidwall/gjson"
)

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
func CreateSenseList(jsonData string, path string) []structs.Sense {
	var senses []structs.Sense

	// Get the JSON array stored in "senses"
	sensesData := gjson.Get(jsonData, path)
	// Iterate over each element in the senses array
	sensesData.ForEach(func(key, value gjson.Result) bool {
		// For each object, get the "type" field
		typ := value.Get("type").String()
		detail := value.Get("details").String()
		senses = append(senses, structs.Sense{
			Type:  typ,
			Range: value.Get("range").String(),

			Detail: detail,
			Acuity: value.Get("acuity").String(),
		})
		return true // Continue iterating
	})

	return senses
}

func parseJSON(data []byte) error {
	fmt.Println(gjson.Get(string(data), "name"))
	fmt.Println(gjson.Get(string(data), "system.abilities.cha.mod"))

	dbObj := structs.Monster{
		Name: gjson.Get(string(data), "name").String(),
		Traits: structs.Traits{
			Rarity:    gjson.Get(string(data), "system.traits.rarity").String(),
			Size:      gjson.Get(string(data), "system.traits.size.value").String(),
			TraitList: ingestJSONList(data, "system.traits.value"),
		},
		Attributes: structs.Attributes{
			Str: gjson.Get(string(data), "system.abilities.str.mod").String(),
			Dex: gjson.Get(string(data), "system.abilities.dex.mod").String(),
			Con: gjson.Get(string(data), "system.abilities.con.mod").String(),
			Wis: gjson.Get(string(data), "system.abilities.wis.mod").String(),
			Int: gjson.Get(string(data), "system.abilities.int.mod").String(),
			Cha: gjson.Get(string(data), "system.abilities.cha.mod").String(),
		},
		Level: gjson.Get(string(data), "system.details.level.value").String(),
		AClass: structs.AC{
			Value:  gjson.Get(string(data), "system.attributes.ac.value").String(),
			Detail: gjson.Get(string(data), "system.attributes.ac.details.value").String(),
		},
		HP: structs.HP{
			Detail: gjson.Get(string(data), "system.attributes.hp.details.value").String(),
			Value:  int(gjson.Get(string(data), "system.attributes.hp.value").Int()),
		},
		Immunities: extractListOfObjectsValues(string(data), "system.attributes.immunities"),
		Languages:  ingestJSONList(data, "system.details.languages.value"),
	}

	jsonData, err := json.MarshalIndent(dbObj, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}
	fmt.Println(string(jsonData))
	return nil
}

// dbObj := structs.Monster{
// 	Name:  payload["name"].(string),
// 	Level:
// 	Traits: structs.Traits{
// 		Rarity:    payload["rarity"].(string),
// 		Size:      payload["size"].(string),
// 		TraitList: payload["traits"].([]string),
// 	},
// 	Attributes: structs.Attributes{
// 		Str: gjson.Get(string(jsonData), "System.abilities.str.value").String(),
// 		Dex: gjson.Get(string(jsonData), "System.abilities.dex.value").String(),
// 		Con: gjson.Get(string(jsonData), "System.abilities.con.value").String(),
// 		Wis: gjson.Get(string(jsonData), "System.abilities.wis.value").String(),
// 		Int: gjson.Get(string(jsonData), "System.abilities.int.value").String(),
// 		Cha: gjson.Get(string(jsonData), "System.abilities.cha.value").String(),
// 	},
// 	Size:       gjson.Get(string(jsonData), "System.traits.size.value").String(),
// 	AClass:     gjson.Get(string(jsonData), "System.attributes.ac.value").String(),
// 	HP:         int(gjson.Get(string(jsonData), "System.attributes.hp.value").Int()),
// 	Immunities: ingestJSONList(jsonData, "System.traits.immunities.value"),
// 	Languages:  ingestJSONList(jsonData, "System.traits.languages.value"),
// 	Perception: gjson.Get(string(jsonData), "System.attributes.perception.value").String(),
// 	// Skills: SkillParser(payload),
// }
//fmt.Println(dbObj)
