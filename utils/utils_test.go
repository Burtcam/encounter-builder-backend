package utils

import (
	"fmt"
	"testing"

	"github.com/Burtcam/encounter-builder-backend/structs"
	"github.com/tidwall/gjson"
)

func TestGetXpBudget(t *testing.T) {
	tests := []struct {
		difficulty string
		pSize      int
		expected   int
	}{
		{"trivial", 4, 40},
		{"low", 4, 60},
		{"moderate", 4, 80},
		{"severe", 4, 120},
		{"extreme", 4, 160},
		{"trivial", 5, 50},
		{"low", 5, 80},
		{"moderate", 5, 100},
		{"severe", 5, 150},
		{"extreme", 5, 200},
	}

	for _, test := range tests {
		result, err := GetXpBudget(test.difficulty, test.pSize)
		if err != nil {
			t.Errorf("Error: %v", err)
			continue
		}
		if result != test.expected {
			t.Errorf("GetXpBudget(%q, %d) = %d; want %d", test.difficulty, test.pSize, result, test.expected)
		}
	}
}
func TestGetXpBudgetInvalidDifficulty(t *testing.T) {
	tests := []struct {
		difficulty string
		pSize      int
	}{
		{"invalid", 4},
		{"unknown", 5},
	}

	for _, test := range tests {
		result, err := GetXpBudget(test.difficulty, test.pSize)
		if err == nil {
			t.Errorf("Expected error for difficulty %q and pSize %d, got result %d", test.difficulty, test.pSize, result)
		}
	}
}
func TestGetXpBudgetNegativePartySize(t *testing.T) {
	tests := []struct {
		difficulty string
		pSize      int
	}{
		{"trivial", -1},
		{"low", -2},
	}

	for _, test := range tests {
		result, err := GetXpBudget(test.difficulty, test.pSize)
		if err == nil {
			t.Errorf("Expected error for difficulty %q and pSize %d, got result %d", test.difficulty, test.pSize, result)
		}
	}
}
func TestGetXpBudgetZeroPartySize(t *testing.T) {
	tests := []struct {
		difficulty string
		pSize      int
	}{
		{"trivial", 0},
		{"low", 0},
	}

	for _, test := range tests {
		result, err := GetXpBudget(test.difficulty, test.pSize)
		if err == nil {
			t.Errorf("Expected error for difficulty %q and pSize %d, got result %d", test.difficulty, test.pSize, result)
		}
	}
}
func TestParsePassives(t *testing.T) {
	// Sample JSON input
	jsonData := `{
            "_id": "UZivbag22tDFAKdd",
            "_stats": {
                "compendiumSource": "Compendium.pf2e.bestiary-ability-glossary-srd.Item.kdhbPaBMK1d1fpbA"
            },
            "img": "systems/pf2e/icons/actions/Passive.webp",
            "name": "Telepathy 100 feet",
            "sort": 700000,
            "system": {
                "actionType": {
                    "value": "passive"
                },
                "actions": {
                    "value": null
                },
                "category": "interaction",
                "description": {
                    "value": "<p>@Localize[PF2E.NPC.Abilities.Glossary.Telepathy]</p>"
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": "Pathfinder Bestiary"
                },
                "rules": [],
                "slug": "telepathy",
                "traits": {
                    "rarity": "common",
                    "value": [
                        "aura",
                        "magical"
                    ]
                }
            },
            "type": "action"
        }`

	// Parse the JSON using gjson
	result := gjson.Parse(jsonData).String()

	// Call parsePassives function
	passive := ParsePassive(result)

	// Validate parsed values
	if passive.Name != "Telepathy 100 feet" {
		t.Errorf("Expected Name 'Telepathy 100 feet', got '%s'", passive.Name)
	}

	if passive.Text != "<p>@Localize[PF2E.NPC.Abilities.Glossary.Telepathy]</p>" {
		t.Errorf("Expected Text '<p>@Localize[PF2E.NPC.Abilities.Glossary.Telepathy]</p>', got '%s'", passive.Text)
	}

	if passive.Category != "interaction" {
		t.Errorf("Expected Category 'interaction', got '%s'", passive.Category)
	}

	expectedTraits := []string{"aura", "magical"}
	for i, trait := range expectedTraits {
		if i >= len(passive.Traits) || passive.Traits[i] != trait {
			t.Errorf("Expected Trait %s at index %d, got '%s'", trait, i, passive.Traits[i])
		}
	}
}
func TestExtractListOfObjectsValues(t *testing.T) {
	// Sample JSON input
	jsonData := `{
		"immunities": [
			{
				"type": "fire",
				"value": "immune"
			},
			{
				"type": "cold",
				"value": "resist"
			}
		]
	}`

	expectedTypes := []string{"fire", "cold"}

	result := extractListOfObjectsValues(jsonData, "immunities")

	if len(result) != len(expectedTypes) {
		t.Fatalf("Expected %d types, got %d", len(expectedTypes), len(result))
	}

	for i, typ := range expectedTypes {
		if result[i] != typ {
			t.Errorf("Expected type '%s' at index %d, got '%s'", typ, i, result[i])
		}
	}
}
func TestIngestJSONList(t *testing.T) {
	// Sample JSON input
	jsonData := `{
		"list": [
			"value1",
			"value2",
			"value3"
		]
	}`

	expectedValues := []string{"value1", "value2", "value3"}

	result := ingestJSONList(jsonData, "list")

	if len(result) != len(expectedValues) {
		t.Fatalf("Expected %d values, got %d", len(expectedValues), len(result))
	}

	for i, val := range expectedValues {
		if result[i] != val {
			t.Errorf("Expected value '%s' at index %d, got '%s'", val, i, result[i])
		}
	}
}

func TestExtractSkills(t *testing.T) {
	// Sample JSON input
	jsonData := `{
            "acrobatics": {
                "base": 21
            },
            "athletics": {
                "base": 27
            },
            "deception": {
                "base": 24
            },
            "intimidation": {
                "base": 26
            },
            "nature": {
                "base": 25
            },
            "stealth": {
                "base": 21,
                "special": [
                    {
                        "base": 25,
                        "label": "in forests",
                        "predicate": [
                            "terrain:forest"
							"terrain:moon"
							"terrain:yourmommasass"
                        ]
                    }
					{
					"base": 54,
					"label": "UNDERWATER",
					"predicate": [
						"terrain:water"
						"terrain:lakes"
						"terrain:ocean"
					]
                    }
                ]
					}`
	expectedValues := []structs.Skill{
		{
			Name:     "acrobatics",
			Value:    21,
			Specials: nil,
		},
		{
			Name:     "athletics",
			Value:    27,
			Specials: nil,
		},
		{
			Name:     "deception",
			Value:    24,
			Specials: nil,
		},
		{
			Name:     "intimidation",
			Value:    26,
			Specials: nil,
		},
		{
			Name:     "nature",
			Value:    25,
			Specials: nil,
		},
		{
			Name:  "stealth",
			Value: 21,
			Specials: []structs.SkillSpecial{
				{
					Value:      25,
					Label:      "in forests",
					Predicates: []string{"terrain:forest", "terrain:moon", "terrain:yourmommasass"},
				},
				{
					Value:      54,
					Label:      "UNDERWATER",
					Predicates: []string{"terrain:water", "terrain:lakes", "terrain:ocean"},
				},
			},
		},
	}
	result := ExtractSkills(gjson.Parse(jsonData))
	if len(result) != len(expectedValues) {
		fmt.Printf("Result: %v\n", result)
		t.Fatalf("Expected %d skills, got %d", len(expectedValues), len(result))
	}
	for i, val := range expectedValues {
		if result[i].Name != val.Name {
			t.Errorf("Expected skill name '%s' at index %d, got '%s'", val.Name, i, result[i].Name)
		}
		if result[i].Value != val.Value {
			t.Errorf("Expected skill value '%d' at index %d, got '%d'", val.Value, i, result[i].Value)
		}
		if len(result[i].Specials) != len(val.Specials) {
			t.Fatalf("Expected %d specials for skill '%s', got %d", len(val.Specials), val.Name, len(result[i].Specials))
		}
		for j, spec := range val.Specials {
			if result[i].Specials[j].Value != spec.Value {
				t.Errorf("Expected special value '%d' at index %d for skill '%s', got '%d'", spec.Value, j, val.Name, result[i].Specials[j].Value)
			}
			if result[i].Specials[j].Label != spec.Label {
				t.Errorf("Expected special label '%s' at index %d for skill '%s', got '%s'", spec.Label, j, val.Name, result[i].Specials[j].Label)
			}
			if len(result[i].Specials[j].Predicates) != len(spec.Predicates) {
				t.Fatalf("Expected %d predicates for special '%s' of skill '%s', got %d", len(spec.Predicates), spec.Label, val.Name, len(result[i].Specials[j].Predicates))
			}
			for k, pred := range spec.Predicates {
				if result[i].Specials[j].Predicates[k] != pred {
					t.Errorf("Expected predicate '%s' at index %d for special '%s' of skill '%s', got '%s'", pred, k, spec.Label, val.Name, result[i].Specials[j].Predicates[k])
				}
			}
		}
	}
}

func TestParseSpontaneousSpellCasting(t *testing.T) {
	// Sample JSON input
	jsonData := `{
            "_id": "6PZisICkQg9iEoQs",
            "img": "systems/pf2e/icons/default-icons/spellcastingEntry.svg",
            "name": "Occult Spontaneous Spells",
            "sort": 100000,
            "system": {
                "autoHeightenLevel": {
                    "value": null
                },
                "description": {
                    "value": ""
                },
                "prepared": {
                    "value": "spontaneous"
                },
                "proficiency": {
                    "value": 1
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": ""
                },
                "rules": [],
                "slots": {
                    "slot1": {
                        "max": 4,
                        "value": 4
                    },
                    "slot2": {
                        "max": 4,
                        "value": 4
                    },
                    "slot3": {
                        "max": 4,
                        "value": 4
                    },
                    "slot4": {
                        "max": 4,
                        "value": 4
                    },
                    "slot5": {
                        "max": 4,
                        "value": 4
                    }
                },
                "slug": null,
                "spelldc": {
                    "dc": 29,
                    "value": 21
                },
                "tradition": {
                    "value": "occult"
                }
            },
            "type": "spellcastingEntry"
        },`
	expected := structs.SpontaneousSpellCasting{
		ID:        "6PZisICkQg9iEoQs",
		Tradition: "occult",
		DC:        29,
		Mod:       "21",
		Slots: []structs.Slot{
			{Level: "slot1", Casts: "4"},
			{Level: "slot2", Casts: "4"},
			{Level: "slot3", Casts: "4"},
			{Level: "slot4", Casts: "4"},
			{Level: "slot5", Casts: "4"},
		},
	}
	result := ParseSpontaneousSpellCasting(jsonData)
	if result.ID != expected.ID {
		t.Errorf("Expected ID '%s', got '%s'", expected.ID, result.ID)
	}
	if result.Tradition != expected.Tradition {
		t.Errorf("Expected Tradition '%s', got '%s'", expected.Tradition, result.Tradition)
	}
	if result.DC != expected.DC {
		t.Errorf("Expected DC '%d', got '%d'", expected.DC, result.DC)
	}
	if result.Mod != expected.Mod {
		t.Errorf("Expected Mod '%s', got '%s'", expected.Mod, result.Mod)
	}
	if len(result.Slots) != len(expected.Slots) {
		t.Fatalf("Expected %d slots, got %d", len(expected.Slots), len(result.Slots))
	}
	for i, slot := range expected.Slots {
		if result.Slots[i].Level != slot.Level {
			t.Errorf("Expected Slot Level '%s' at index %d, got '%s'", slot.Level, i, result.Slots[i].Level)
		}
		if result.Slots[i].Casts != slot.Casts {
			t.Errorf("Expected Slot Casts '%s' at index %d, got '%s'", slot.Casts, i, result.Slots[i].Casts)
		}
	}
	// Check for any additional fields in the result
	if len(result.Slots) > len(expected.Slots) {
		t.Errorf("Expected no additional slots, got %d", len(result.Slots)-len(expected.Slots))
	}
	if len(result.Slots) < len(expected.Slots) {
		t.Errorf("Expected no missing slots, got %d", len(expected.Slots)-len(result.Slots))
	}
	if result.Slots[0].Level != expected.Slots[0].Level {
		t.Errorf("Expected Slot Level '%s', got '%s'", expected.Slots[0].Level, result.Slots[0].Level)
	}
	if result.Slots[0].Casts != expected.Slots[0].Casts {
		t.Errorf("Expected Slot Casts '%s', got '%s'", expected.Slots[0].Casts, result.Slots[0].Casts)
	}
	if result.Slots[1].Level != expected.Slots[1].Level {
		t.Errorf("Expected Slot Level '%s', got '%s'", expected.Slots[1].Level, result.Slots[1].Level)
	}
	if result.Slots[1].Casts != expected.Slots[1].Casts {
		t.Errorf("Expected Slot Casts '%s', got '%s'", expected.Slots[1].Casts, result.Slots[1].Casts)
	}
	if result.Slots[2].Level != expected.Slots[2].Level {
		t.Errorf("Expected Slot Level '%s', got '%s'", expected.Slots[2].Level, result.Slots[2].Level)
	}
	if result.Slots[2].Casts != expected.Slots[2].Casts {
		t.Errorf("Expected Slot Casts '%s', got '%s'", expected.Slots[2].Casts, result.Slots[2].Casts)
	}
	if result.Slots[3].Level != expected.Slots[3].Level {
		t.Errorf("Expected Slot Level '%s', got '%s'", expected.Slots[3].Level, result.Slots[3].Level)
	}
	if result.Slots[3].Casts != expected.Slots[3].Casts {
		t.Errorf("Expected Slot Casts '%s', got '%s'", expected.Slots[3].Casts, result.Slots[3].Casts)
	}
	if result.Slots[4].Level != expected.Slots[4].Level {
		t.Errorf("Expected Slot Level '%s', got '%s'", expected.Slots[4].Level, result.Slots[4].Level)
	}
	if result.Slots[4].Casts != expected.Slots[4].Casts {
		t.Errorf("Expected Slot Casts '%s', got '%s'", expected.Slots[4].Casts, result.Slots[4].Casts)
	}
}

func TestParsePreparedSpellCasting(t *testing.T) {
	data := `{
            "_id": "9h6KJeGxzm8rEPaD",
            "img": "systems/pf2e/icons/default-icons/spellcastingEntry.svg",
            "name": "Primal Prepared Spells",
            "sort": 100000,
            "system": {
                "autoHeightenLevel": {
                    "value": 6
                },
                "description": {
                    "value": ""
                },
                "prepared": {
                    "flexible": false,
                    "value": "prepared"
                },
                "proficiency": {
                    "value": 1
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": ""
                },
                "rules": [],
                "showSlotlessLevels": {
                    "value": false
                },
                "slots": {
                    "slot0": {
                        "max": 5,
                        "prepared": [
                            {
                                "id": "cgw07bSj0UprtiUE"
                            },
                            {
                                "id": "GeRqpkpFNtXrmbgm"
                            },
                            {
                                "id": "tLuFR0oqghOXKzbd"
                            },
                            {
                                "id": "wmqu97fbZeHaDCYh"
                            },
                            {
                                "id": "ELMWrZpjcRl1T4RG"
                            }
                        ]
                    },
                    "slot1": {
                        "max": 3,
                        "prepared": [
                            {
                                "id": "K2hzbKGlsnbs4Oim"
                            },
                            {
                                "id": "YfWayh8Vf56Z3brL"
                            },
                            {
                                "id": "ZiYYZgtUKyVmJTXf"
                            }
                        ]
                    },
                    "slot2": {
                        "max": 3,
                        "prepared": [
                            {
                                "id": "lyHhpzUmgixU51K3"
                            },
                            {
                                "id": "JxsY3WYSjn7MwRgz"
                            },
                            {
                                "id": "9YyN3ZnrZrlMGETw"
                            }
                        ]
                    },
                    "slot3": {
                        "max": 3,
                        "prepared": [
                            {
                                "id": "uu8jCMiKsmK3daVq"
                            },
                            {
                                "id": "kKqJb4vg5dRnYkWw"
                            },
                            {
                                "id": "gSRFsZkX8Qu19CEz"
                            }
                        ]
                    },
                    "slot4": {
                        "max": 3,
                        "prepared": [
                            {
                                "id": "Pc8OabeDh0D0QoNn"
                            },
                            {
                                "id": "T6VXVjgqGBXusSVY"
                            },
                            {
                                "id": "VVTdSugZYXwWMIqG"
                            }
                        ]
                    },
                    "slot5": {
                        "max": 3,
                        "prepared": [
                            {
                                "id": "Pr9Ih78tzMSQfxvf"
                            },
                            {
                                "id": "7YdPP01kBJ4BN5CS"
                            },
                            {
                                "id": "D5sHvAzd2vbdfA3E"
                            }
                        ]
                    },
                    "slot6": {
                        "max": 3,
                        "prepared": [
                            {
                                "id": "iUZaBJdkAt5wfkw9"
                            },
                            {
                                "id": "Qc0rR7NFVpIq7lgF"
                            },
                            {
                                "id": "ECGCJIVLGkNeDpoK"
                            }
                        ]
                    }
                },
                "slug": null,
                "spelldc": {
                    "dc": 34,
                    "mod": 0,
                    "value": 28
                },
                "tradition": {
                    "value": "primal"
                }
            },
            "type": "spellcastingEntry"
        }`
	expected := structs.PreparedSpellCasting{
		DC:          34,
		Mod:         "28",
		Tradition:   "primal",
		ID:          "9h6KJeGxzm8rEPaD",
		Description: "",
		Slots: []structs.PreparedSlot{
			{Level: "slot0", SpellID: "cgw07bSj0UprtiUE"},
			{Level: "slot0", SpellID: "GeRqpkpFNtXrmbgm"},
			{Level: "slot0", SpellID: "tLuFR0oqghOXKzbd"},
			{Level: "slot0", SpellID: "wmqu97fbZeHaDCYh"},
			{Level: "slot0", SpellID: "ELMWrZpjcRl1T4RG"},
			{Level: "slot1", SpellID: "K2hzbKGlsnbs4Oim"},
			{Level: "slot1", SpellID: "YfWayh8Vf56Z3brL"},
			{Level: "slot1", SpellID: "ZiYYZgtUKyVmJTXf"},
			{Level: "slot2", SpellID: "lyHhpzUmgixU51K3"},
			{Level: "slot2", SpellID: "JxsY3WYSjn7MwRgz"},
			{Level: "slot2", SpellID: "9YyN3ZnrZrlMGETw"},
			{Level: "slot3", SpellID: "uu8jCMiKsmK3daVq"},
			{Level: "slot3", SpellID: "kKqJb4vg5dRnYkWw"},
			{Level: "slot3", SpellID: "gSRFsZkX8Qu19CEz"},
			{Level: "slot4", SpellID: "Pc8OabeDh0D0QoNn"},
			{Level: "slot4", SpellID: "T6VXVjgqGBXusSVY"},
			{Level: "slot4", SpellID: "VVTdSugZYXwWMIqG"},
			{Level: "slot5", SpellID: "Pr9Ih78tzMSQfxvf"},
			{Level: "slot5", SpellID: "7YdPP01kBJ4BN5CS"},
			{Level: "slot5", SpellID: "D5sHvAzd2vbdfA3E"},
			{Level: "slot6", SpellID: "iUZaBJdkAt5wfkw9"},
			{Level: "slot6", SpellID: "Qc0rR7NFVpIq7lgF"},
			{Level: "slot6", SpellID: "ECGCJIVLGkNeDpoK"},
		},
	}
	result := ParsePreparedSpellCasting(data)
	fmt.Println(result)
	if result.DC != expected.DC {
		t.Errorf("Expected DC '%d', got '%d'", expected.DC, result.DC)
	}
	if result.Mod != expected.Mod {
		t.Errorf("Expected Mod '%s', got '%s'", expected.Mod, result.Mod)
	}
	if result.Tradition != expected.Tradition {
		t.Errorf("Expected Tradition '%s', got '%s'", expected.Tradition, result.Tradition)
	}
	if result.ID != expected.ID {
		t.Errorf("Expected ID '%s', got '%s'", expected.ID, result.ID)
	}
	if result.Description != expected.Description {
		t.Errorf("Expected Description '%s', got '%s'", expected.Description, result.Description)
	}
	if len(result.Slots) != len(expected.Slots) {
		t.Fatalf("Expected %d slots, got %d", len(expected.Slots), len(result.Slots))
	}
	for i, slot := range expected.Slots {
		if result.Slots[i].Level != slot.Level {
			t.Errorf("Expected Slot Level '%s' at index %d, got '%s'", slot.Level, i, result.Slots[i].Level)
		}
		if result.Slots[i].SpellID != slot.SpellID {
			t.Errorf("Expected Slot SpellID '%s' at index %d, got '%s'", slot.SpellID, i, result.Slots[i].SpellID)
		}
	}
	if len(result.Slots) > len(expected.Slots) {
		t.Errorf("Expected no additional slots, got %d", len(result.Slots)-len(expected.Slots))
	}
	if len(result.Slots) < len(expected.Slots) {
		t.Errorf("Expected no missing slots, got %d", len(expected.Slots)-len(result.Slots))
	}
	if result.Slots[0].Level != expected.Slots[0].Level {
		t.Errorf("Expected Slot Level '%s', got '%s'", expected.Slots[0].Level, result.Slots[0].Level)
	}
	if result.Slots[0].SpellID != expected.Slots[0].SpellID {
		t.Errorf("Expected Slot SpellID '%s', got '%s'", expected.Slots[0].SpellID, result.Slots[0].SpellID)
	}
	if result.Slots[1].Level != expected.Slots[1].Level {
		t.Errorf("Expected Slot Level '%s', got '%s'", expected.Slots[1].Level, result.Slots[1].Level)
	}
	if result.Slots[1].SpellID != expected.Slots[1].SpellID {
		t.Errorf("Expected Slot SpellID '%s', got '%s'", expected.Slots[1].SpellID, result.Slots[1].SpellID)
	}
	if result.Slots[2].Level != expected.Slots[2].Level {
		t.Errorf("Expected Slot Level '%s', got '%s'", expected.Slots[2].Level, result.Slots[2].Level)
	}
	if result.Slots[2].SpellID != expected.Slots[2].SpellID {
		t.Errorf("Expected Slot SpellID '%s', got '%s'", expected.Slots[2].SpellID, result.Slots[2].SpellID)
	}
	if result.Slots[3].Level != expected.Slots[3].Level {
		t.Errorf("Expected Slot Level '%s', got '%s'", expected.Slots[3].Level, result.Slots[3].Level)
	}
	if result.Slots[3].SpellID != expected.Slots[3].SpellID {
		t.Errorf("Expected Slot SpellID '%s', got '%s'", expected.Slots[3].SpellID, result.Slots[3].SpellID)
	}
	if result.Slots[4].Level != expected.Slots[4].Level {
		t.Errorf("Expected Slot Level '%s', got '%s'", expected.Slots[4].Level, result.Slots[4].Level)
	}
	if result.Slots[4].SpellID != expected.Slots[4].SpellID {
		t.Errorf("Expected Slot SpellID '%s', got '%s'", expected.Slots[4].SpellID, result.Slots[4].SpellID)
	}
	if result.Slots[5].Level != expected.Slots[5].Level {
		t.Errorf("Expected Slot Level '%s', got '%s'", expected.Slots[5].Level, result.Slots[5].Level)
	}
	if result.Slots[5].SpellID != expected.Slots[5].SpellID {
		t.Errorf("Expected Slot SpellID '%s', got '%s'", expected.Slots[5].SpellID, result.Slots[5].SpellID)
	}
	if result.Slots[6].Level != expected.Slots[6].Level {
		t.Errorf("Expected Slot Level '%s', got '%s'", expected.Slots[6].Level, result.Slots[6].Level)
	}
	if result.Slots[6].SpellID != expected.Slots[6].SpellID {
		t.Errorf("Expected Slot SpellID '%s', got '%s'", expected.Slots[6].SpellID, result.Slots[6].SpellID)
	}
	if result.Slots[7].Level != expected.Slots[7].Level {
		t.Errorf("Expected Slot Level '%s', got '%s'", expected.Slots[7].Level, result.Slots[7].Level)
	}
	if result.Slots[7].SpellID != expected.Slots[7].SpellID {
		t.Errorf("Expected Slot SpellID '%s', got '%s'", expected.Slots[7].SpellID, result.Slots[7].SpellID)
	}
}
func TestParseFocusSpellCasting(t *testing.T) {
	jsonData := ` {
            "_id": "EDPFYDhj0ZOTpRmX",
            "img": "systems/pf2e/icons/default-icons/spellcastingEntry.svg",
            "name": "Animal Order Spells",
            "sort": 200000,
            "system": {
                "autoHeightenLevel": {
                    "value": null
                },
                "description": {
                    "value": ""
                },
                "prepared": {
                    "flexible": false,
                    "value": "focus"
                },
                "proficiency": {
                    "value": 1
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": ""
                },
                "rules": [],
                "slots": {},
                "slug": null,
                "spelldc": {
                    "dc": 24,
                    "mod": 0,
                    "value": 16
                },
                "tradition": {
                    "value": "primal"
                }
            },
            "type": "spellcastingEntry"
        },`
	expected := structs.FocusSpellCasting{
		DC:             24,
		Mod:            "16",
		Tradition:      "primal",
		ID:             "EDPFYDhj0ZOTpRmX",
		FocusSpellList: []structs.Spell{},
		Description:    stripHTMLUsingBluemonday(""),
		Name:           "Animal Order Spells",
		CastLevel:      "",
	}
	result := ParseFocusSpellCasting(jsonData)
	if result.DC != expected.DC {
		t.Errorf("Expected DC '%d', got '%d'", expected.DC, result.DC)
	}
	if result.Mod != expected.Mod {
		t.Errorf("Expected Mod '%s', got '%s'", expected.Mod, result.Mod)
	}
	if result.Tradition != expected.Tradition {
		t.Errorf("Expected Tradition %s, got %s", expected.Tradition, result.Tradition)
	}
	if result.ID != expected.ID {
		t.Errorf("Expected ID %s, got %s", expected.ID, result.ID)
	}
	if result.Description != expected.Description {
		t.Errorf("Expected Description %s, got %s", expected.Description, result.Description)
	}
	if result.Name != expected.Name {
		t.Errorf("Expected Name '%s', got '%s'", expected.Name, result.Name)
	}
	if len(expected.FocusSpellList) != len(result.FocusSpellList) {
		t.Errorf("Expected %d yet, got %d", len(expected.FocusSpellList), len(result.FocusSpellList))
	}
	if expected.CastLevel != result.CastLevel {
		t.Errorf("Expected castLevel %s, got %s", expected.CastLevel, result.CastLevel)
	}
}
func TestParseInnateSpellCasting(t *testing.T) {
	jsonData := ` {
            "_id": "yI8fil9Hp8Ob0BcY",
            "img": "systems/pf2e/icons/default-icons/spellcastingEntry.svg",
            "name": "Occult Innate Spells",
            "sort": 100000,
            "system": {
                "autoHeightenLevel": {
                    "value": null
                },
                "description": {
                    "value": ""
                },
                "prepared": {
                    "flexible": false,
                    "value": "innate"
                },
                "proficiency": {
                    "value": 0
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": ""
                },
                "rules": [],
                "slots": {},
                "slug": null,
                "spelldc": {
                    "dc": 31,
                    "mod": 0,
                    "value": 21
                },
                "tradition": {
                    "value": "occult"
                }
            },
            "type": "spellcastingEntry"
        },`

	expected := structs.InnateSpellCasting{
		DC:          31,
		Mod:         "21",
		SpellUses:   []structs.SpellUse{},
		Tradition:   "occult",
		ID:          "yI8fil9Hp8Ob0BcY",
		Description: stripHTMLUsingBluemonday(""),
		Name:        "Occult Innate Spells",
	}

	result := ParseInnateSpellCasting(jsonData)
	if result.DC != expected.DC {
		t.Errorf("Expected DC '%d', got '%d'", expected.DC, result.DC)
	}
	if result.Mod != expected.Mod {
		t.Errorf("Expected Mod '%s', got '%s'", expected.Mod, result.Mod)
	}
	if result.Tradition != expected.Tradition {
		t.Errorf("Expected Tradition %s, got %s", expected.Tradition, result.Tradition)
	}
	if result.ID != expected.ID {
		t.Errorf("Expected ID %s, got %s", expected.ID, result.ID)
	}
	if result.Description != expected.Description {
		t.Errorf("Expected Description %s, got %s", expected.Description, result.Description)
	}
	if result.Name != expected.Name {
		t.Errorf("Expected Name '%s', got '%s'", expected.Name, result.Name)
	}
	if len(expected.SpellUses) != len(result.SpellUses) {
		t.Errorf("Expected %d yet, got %d", len(expected.SpellUses), len(result.SpellUses))
	}

}

func TestParseDamageBlocks(t *testing.T) {
	jsonData := `{
            "_id": "7SJO477OusJy7wpB",
            "img": "systems/pf2e/icons/default-icons/melee.svg",
            "name": "Jaws",
            "sort": 3200000,
            "system": {
                "attack": {
                    "value": ""
                },
                "attackEffects": {
                    "custom": "",
                    "value": []
                },
                "bonus": {
                    "value": 29
                },
                "damageRolls": {
                    "0": {
                        "damage": "3d10+13",
                        "damageType": "piercing"
                    },
                    "e30481rbp6g1b2cgivij": {
                        "damage": "2d6",
                        "damageType": "poison"
                    }
                },
                "description": {
                    "value": ""
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": ""
                },
                "rules": [],
                "slug": null,
                "traits": {
                    "rarity": "common",
                    "value": [
                        "magical",
                        "reach-15",
                        "unarmed"
                    ]
                },
                "weaponType": {
                    "value": "melee"
                }
            },
            "type": "melee"
        },`
	expected := []structs.DamageBlock{
		{
			DamageRoll: "3d10+13",
			DamageType: "piercing",
		},
		{
			DamageRoll: "2d6",
			DamageType: "poison",
		},
	}
	result := ParseDamageBlocks(jsonData)
	if len(result) != len(expected) {
		t.Fatalf("Expected %d damage blocks, got %d", len(expected), len(result))
	}
	for i, block := range expected {
		if result[i].DamageRoll != block.DamageRoll {
			t.Errorf("Expected DamageRoll '%s' at index %d, got '%s'", block.DamageRoll, i, result[i].DamageRoll)
		}
		if result[i].DamageType != block.DamageType {
			t.Errorf("Expected DamageType '%s' at index %d, got '%s'", block.DamageType, i, result[i].DamageType)
		}
	}
}

func TestParseDamageEffects(t *testing.T) {
	jsonData := `{
            "_id": "7SJO477OusJy7wpB",
            "img": "systems/pf2e/icons/default-icons/melee.svg",
            "name": "Jaws",
            "sort": 3200000,
            "system": {
                "attack": {
                    "value": ""
                },
                "attackEffects": {
                    "custom": "hello",
                    "value": ["stunned", "dazed", "slowed", "Confused"]
                },
                "damageRolls": {
                    "0": {
                        "damage": "3d10+13",
                        "damageType": "piercing"
                    },
                    "e30481rbp6g1b2cgivij": {
                        "damage": "2d6",
                        "damageType": "poison"
                    }
                },
                "description": {
                    "value": ""
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": ""
                },
                "rules": [],
                "slug": null,
                "traits": {
                    "rarity": "common",
                    "value": [
                        "magical",
                        "reach-15",
                        "unarmed"
                    ]
                },
                "weaponType": {
                    "value": "melee"
                }
            },
            "type": "melee"
        },`
	expected := structs.DamageEffect{
		CustomString: "hello",
		Value:        []string{"stunned", "dazed", "slowed", "Confused"},
	}
	result := ParseDamageEffects(jsonData)
	if result.CustomString != expected.CustomString {
		t.Errorf("Expected CustomString '%s', got '%s'", expected.CustomString, result.CustomString)
	}
	if len(result.Value) != len(expected.Value) {
		t.Fatalf("Expected %d effects, got %d", len(expected.Value), len(result.Value))
	}
	for i, effect := range expected.Value {
		if result.Value[i] != effect {
			t.Errorf("Expected effect '%s' at index %d, got '%s'", effect, i, result.Value[i])
		}
	}
}

func TestParseWeapon(t *testing.T) {
	jsonData := `{
            "_id": "7SJO477OusJy7wpB",
            "img": "systems/pf2e/icons/default-icons/melee.svg",
            "name": "Jaws",
            "sort": 3200000,
            "system": {
                "attack": {
                    "value": ""
                },
                "attackEffects": {
                    "custom": "",
                    "value": []
                },
                "bonus": {
                    "value": 29
                },
                "damageRolls": {
                    "0": {
                        "damage": "3d10+13",
                        "damageType": "piercing"
                    },
                    "e30481rbp6g1b2cgivij": {
                        "damage": "2d6",
                        "damageType": "poison"
                    }
                },
                "description": {
                    "value": ""
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": ""
                },
                "rules": [],
                "slug": null,
                "traits": {
                    "rarity": "common",
                    "value": [
                        "magical",
                        "reach-15",
                        "unarmed"
                    ]
                },
                "weaponType": {
                    "value": "melee"
                }
            },
            "type": "melee"
        },`
	expected := structs.Attack{
		Type:       "melee",
		Name:       "Jaws",
		ToHitBonus: "29",
		DamageBlocks: []structs.DamageBlock{
			{
				DamageRoll: "3d10+13",
				DamageType: "piercing",
			},
			{
				DamageRoll: "2d6",
				DamageType: "poison",
			},
		},
		Effects: structs.DamageEffect{
			CustomString: "",
			Value:        []string{},
		},
		Traits: []string{"magical", "reach-15", "unarmed"},
	}
	result := ParseWeapon(jsonData)
	fmt.Println(result)
	if result.Type != expected.Type {
		t.Errorf("Expected Type '%s', got '%s'", expected.Type, result.Type)
	}
	if result.Name != expected.Name {
		t.Errorf("Expected Name '%s', got '%s'", expected.Name, result.Name)
	}
	if result.ToHitBonus != expected.ToHitBonus {
		t.Errorf("Expected ToHitBonus '%s', got '%s'", expected.ToHitBonus, result.ToHitBonus)
	}
	if len(result.DamageBlocks) != len(expected.DamageBlocks) {
		t.Fatalf("Expected %d damage blocks, got %d", len(expected.DamageBlocks), len(result.DamageBlocks))
	}
	for i, block := range expected.DamageBlocks {
		if result.DamageBlocks[i].DamageRoll != block.DamageRoll {
			t.Errorf("Expected DamageRoll '%s' at index %d, got '%s'", block.DamageRoll, i, result.DamageBlocks[i].DamageRoll)
		}
		if result.DamageBlocks[i].DamageType != block.DamageType {
			t.Errorf("Expected DamageType '%s' at index %d, got '%s'", block.DamageType, i, result.DamageBlocks[i].DamageType)
		}
	}
	if result.Effects.CustomString != expected.Effects.CustomString {
		t.Errorf("Expected Effects CustomString '%s', got '%s'", expected.Effects.CustomString, result.Effects.CustomString)
	}
	if len(result.Effects.Value) != len(expected.Effects.Value) {
		t.Fatalf("Expected %d effects, got %d", len(expected.Effects.Value), len(result.Effects.Value))
	}
	if len(result.Traits) != len(expected.Traits) {
		t.Fatalf("Expected %d traits, got %d", len(expected.Traits), len(result.Traits))
	}
	for i, trait := range expected.Traits {
		if result.Traits[i] != trait {
			t.Errorf("Expected Trait '%s' at index %d, got '%s'", trait, i, result.Traits[i])
		}
	}
	for i, effect := range expected.Effects.Value {
		if result.Effects.Value[i] != effect {
			t.Errorf("Expected effect '%s' at index %d, got '%s'", effect, i, result.Effects.Value[i])
		}
	}
	// Check for any additional fields in the result
	if len(result.DamageBlocks) > len(expected.DamageBlocks) {
		t.Errorf("Expected no additional damage blocks, got %d", len(result.DamageBlocks)-len(expected.DamageBlocks))
	}
	if len(result.DamageBlocks) < len(expected.DamageBlocks) {
		t.Errorf("Expected no missing damage blocks, got %d", len(expected.DamageBlocks)-len(result.DamageBlocks))
	}
	if len(result.Effects.Value) > len(expected.Effects.Value) {
		t.Errorf("Expected no additional effects, got %d", len(result.Effects.Value)-len(expected.Effects.Value))
	}
	if len(result.Effects.Value) < len(expected.Effects.Value) {
		t.Errorf("Expected no missing effects, got %d", len(expected.Effects.Value)-len(result.Effects.Value))
	}
	if len(result.Traits) > len(expected.Traits) {
		t.Errorf("Expected no additional traits, got %d", len(result.Traits)-len(expected.Traits))
	}
	if len(result.Traits) < len(expected.Traits) {
		t.Errorf("Expected no missing traits, got %d", len(expected.Traits)-len(result.Traits))
	}
}

func TestParseFreeAction(t *testing.T) {
	jsonData := `{
            "_id": "JPj4ayUtkVtkvYCy",
            "img": "systems/pf2e/icons/actions/FreeAction.webp",
            "name": "Consume Light",
            "sort": 600000,
            "system": {
                "actionType": {
                    "value": "free"
                },
                "actions": {
                    "value": null
                },
                "category": "offensive",
                "description": {
                    "value": "<p><strong>Trigger</strong> The voidglutton casts @UUID[Compendium.pf2e.spells-srd.Item.Darkness]</p>\n<hr />\n<p><strong>Effect</strong> The voidglutton extinguishes its Glow as part of Casting the Spell. It becomes @UUID[Compendium.pf2e.conditionitems.Item.Invisible] as long as it remains in the area of darkness. If the voidglutton uses a hostile action, its invisibility ends as soon as the hostile action is completed.</p>"
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": ""
                },
                "rules": [],
                "slug": null,
                "traits": {
                    "rarity": "common",
                    "value": [
                        "darkness",
                        "occult"
                    ]
                }
            },
            "type": "action"
        },`
	expected := structs.FreeAction{
		Name: "Consume Light",
		Text: stripHTMLUsingBluemonday("<p><strong>Trigger</strong> The voidglutton casts @UUID[Compendium.pf2e.spells-srd.Item.Darkness]</p>\n<hr />\n<p><strong>Effect</strong> The voidglutton extinguishes its Glow as part of Casting the Spell. It becomes @UUID[Compendium.pf2e.conditionitems.Item.Invisible] as long as it remains in the area of darkness. If the voidglutton uses a hostile action, its invisibility ends as soon as the hostile action is completed.</p>"),

		Traits:   []string{"darkness", "occult"},
		Category: "offensive",
		Rarity:   "common",
	}
	result := ParseFreeAction(jsonData)
	if result.Name != expected.Name {
		t.Errorf("Expected Name '%s', got '%s'", expected.Name, result.Name)
	}
	if result.Text != expected.Text {
		t.Errorf("Expected Text '%s', got '%s'", expected.Text, result.Text)
	}
	if result.Category != expected.Category {
		t.Errorf("Expected Category '%s', got '%s'", expected.Category, result.Category)
	}
	if result.Rarity != expected.Rarity {
		t.Errorf("Expected Rarity '%s', got '%s'", expected.Rarity, result.Rarity)
	}
	if len(result.Traits) != len(expected.Traits) {
		t.Fatalf("Expected %d traits, got %d", len(expected.Traits), len(result.Traits))
	}

	for i, trait := range expected.Traits {
		if result.Traits[i] != trait {
			t.Errorf("Expected Trait '%s' at index %d, got '%s'", trait, i, result.Traits[i])
		}
	}

}

func TestParseReaction(t *testing.T) {
	jsonData := ` {
            "_id": "tHyoisfllt6q3L0n",
            "img": "systems/pf2e/icons/actions/Reaction.webp",
            "name": "Fed by Water",
            "sort": 3800000,
            "system": {
                "actionType": {
                    "value": "reaction"
                },
                "actions": {
                    "value": null
                },
                "category": "defensive",
                "description": {
                    "value": "<p><strong>Frequency</strong> once per hour</p>\n<p><strong>Trigger</strong> The forest dragon is targeted with a water spell or effect</p>\n<hr />\n<p><strong>Effect</strong> The forest dragon gains [[/r 35 #Temporary Hit Points]]{35 temporary Hit Points}.</p>"
                },
                "frequency": {
                    "max": 1,
                    "per": "PT1H"
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": ""
                },
                "rules": [],
                "slug": null,
                "traits": {
                    "rarity": "common",
                    "value": [
                        "healing",
                        "primal"
                    ]
                }
            },
            "type": "action"
        },`
	expected := structs.Reaction{
		Name: "Fed by Water",
		Text: stripHTMLUsingBluemonday("<p><strong>Frequency</strong> once per hour</p>\n<p><strong>Trigger</strong> The forest dragon is targeted with a water spell or effect</p>\n<hr />\n<p><strong>Effect</strong> The forest dragon gains [[/r 35 #Temporary Hit Points]]{35 temporary Hit Points}.</p>"),

		Traits:   []string{"healing", "primal"},
		Category: "defensive",
		Rarity:   "common",
	}
	result := ParseReaction(jsonData)
	if result.Name != expected.Name {
		t.Errorf("Expected Name '%s', got '%s'", expected.Name, result.Name)
	}
	if result.Text != expected.Text {
		t.Errorf("Expected Text '%s', got '%s'", expected.Text, result.Text)
	}
	if result.Category != expected.Category {
		t.Errorf("Expected Category '%s', got '%s'", expected.Category, result.Category)
	}
	if result.Rarity != expected.Rarity {
		t.Errorf("Expected Rarity '%s', got '%s'", expected.Rarity, result.Rarity)
	}
	if len(result.Traits) != len(expected.Traits) {
		t.Fatalf("Expected %d traits, got %d", len(expected.Traits), len(result.Traits))
	}

	for i, trait := range expected.Traits {
		if result.Traits[i] != trait {
			t.Errorf("Expected Trait '%s' at index %d, got '%s'", trait, i, result.Traits[i])
		}
	}

}

func TestParseAction(t *testing.T) {
	jsonData := ` {
            "_id": "k2Num39uDHGiZwTm",
            "img": "systems/pf2e/icons/actions/TwoActions.webp",
            "name": "Breath Weapon",
            "sort": 4000000,
            "system": {
                "actionType": {
                    "value": "action"
                },
                "actions": {
                    "value": 2
                },
                "category": "offensive",
                "description": {
                    "value": "<p>The dragon unleashes a swarm of insects that deals @Damage[14d6[piercing]|options:area-damage] damage in a @Template[cone|distance:40] (@Check[reflex|dc:34|basic|options:area-effect] save) before dispersing.</p>\n<p>A creature that critically fails is @UUID[Compendium.pf2e.conditionitems.Item.Stunned]{Stunned 2} from the insects' venom; this is a poison effect.</p>\n<p>The dragon can't use Breath Weapon again for [[/gmr 1d4 #Recharge Breath Weapon]]{1d4 rounds}.</p>"
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": ""
                },
                "rules": [],
                "slug": null,
                "traits": {
                    "rarity": "common",
                    "value": [
                        "primal"
                    ]
                }
            },
            "type": "action"
        },`
	expected := structs.Action{
		Name:     "Breath Weapon",
		Text:     stripHTMLUsingBluemonday("<p>The dragon unleashes a swarm of insects that deals @Damage[14d6[piercing]|options:area-damage] damage in a @Template[cone|distance:40] (@Check[reflex|dc:34|basic|options:area-effect] save) before dispersing.</p>\n<p>A creature that critically fails is @UUID[Compendium.pf2e.conditionitems.Item.Stunned]{Stunned 2} from the insects' venom; this is a poison effect.</p>\n<p>The dragon can't use Breath Weapon again for [[/gmr 1d4 #Recharge Breath Weapon]]{1d4 rounds}.</p>"),
		Traits:   []string{"primal"},
		Actions:  "2",
		Category: "offensive",
		Rarity:   "common",
	}
	result := ParseAction(jsonData)

	if result.Name != expected.Name {
		t.Errorf("Expected Name '%s', got '%s'", expected.Name, result.Name)
	}
	if result.Text != expected.Text {
		t.Errorf("Expected Text '%s', got '%s'", expected.Text, result.Text)
	}
	if result.Category != expected.Category {
		t.Errorf("Expected Category '%s', got '%s'", expected.Category, result.Category)
	}
	if result.Rarity != expected.Rarity {
		t.Errorf("Expected Rarity '%s', got '%s'", expected.Rarity, result.Rarity)
	}
	if len(result.Traits) != len(expected.Traits) {
		t.Fatalf("Expected %d traits, got %d", len(expected.Traits), len(result.Traits))
	}
	if expected.Actions != result.Actions {
		t.Errorf("Expected %s Actions, got %s'", expected.Actions, result.Actions)
	}

	for i, trait := range expected.Traits {
		if result.Traits[i] != trait {
			t.Errorf("Expected Trait '%s' at index %d, got '%s'", trait, i, result.Traits[i])
		}
	}
}

func TestIngestSpontaneousSpell(t *testing.T) {
	jsonData := `{
            "_id": "N5cIxpCa1E4SqZi7",
            "_stats": {
                "compendiumSource": "Compendium.pf2e.spells-srd.Item.9AAkVUCwF6WVNNY2"
            },
            "img": "icons/magic/lightning/bolt-strike-sparks-blue.webp",
            "name": "Lightning Bolt",
            "sort": 900000,
            "system": {
                "area": {
                    "type": "line",
                    "value": 120
                },
                "cost": {
                    "value": ""
                },
                "counteraction": false,
                "damage": {
                    "0": {
                        "applyMod": false,
                        "category": null,
                        "formula": "4d12",
                        "kinds": [
                            "damage"
                        ],
                        "materials": [],
                        "type": "electricity"
                    }
                },
                "defense": {
                    "save": {
                        "basic": true,
                        "statistic": "reflex"
                    }
                },
                "description": {
                    "value": "<p>A bolt of lightning strikes outward from your hand, dealing 4d12 electricity damage.</p>\n<hr />\n<p><strong>Heightened (+1)</strong> The damage increases by 1d12.</p>"
                },
                "duration": {
                    "sustained": false,
                    "value": ""
                },
                "heightening": {
                    "damage": {
                        "0": "1d12"
                    },
                    "interval": 1,
                    "type": "interval"
                },
                "level": {
                    "value": 3
                },
                "location": {
                    "heightenedLevel": 5,
                    "value": "6PZisICkQg9iEoQs"
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": "Pathfinder Core Rulebook"
                },
                "range": {
                    "value": ""
                },
                "requirements": "",
                "rules": [],
                "slug": "lightning-bolt",
                "target": {
                    "value": ""
                },
                "time": {
                    "value": "2"
                },
                "traits": {
                    "rarity": "common",
                    "traditions": [
                        "arcane",
                        "primal"
                    ],
                    "value": [
                        "concentrate",
                        "electricity",
                        "manipulate"
                    ]
                }
            },
            "type": "spell"
        },`
	// location == Null on rituals? Need a different mechanism for those.
	expected := structs.Spell{
		Name:           "Lightning Bolt",
		ID:             "N5cIxpCa1E4SqZi7",
		CastLevel:      "5",
		SpellBaseLevel: "3",
		Description:    stripHTMLUsingBluemonday("<p>A bolt of lightning strikes outward from your hand, dealing 4d12 electricity damage.</p>\n<hr />\n<p><strong>Heightened (+1)</strong> The damage increases by 1d12.</p>"),
		Range:          "",
		Duration: structs.DurationBlock{
			Sustained: false,
			Duration:  "",
		},
		Area: structs.SpellArea{
			Type:   "line",
			Value:  "120",
			Detail: "",
		},
		Targets: "",
		Traits:  []string{"concentrate", "electricity", "manipulate"},
		Defense: structs.DefenseBlock{
			Save:  "reflex",
			Basic: true,
		},
		CastTime:                    "2",
		CastRequirements:            "",
		Rarity:                      "common",
		AtWill:                      false,
		SpellCastingBlockLocationID: "6PZisICkQg9iEoQs",
		Uses:                        "1",
		Ritual:                      false,
		RitualData: structs.RitualData{
			PrimaryCheck:     "",
			SecondaryCasters: "",
			SecondaryCheck:   "",
		}}

	result := ParseSpell(jsonData)

	if result.Name != expected.Name {
		t.Errorf("Expected Name '%s', got '%s'", expected.Name, result.Name)
	}
	if result.ID != expected.ID {
		t.Errorf("Expected ID '%s', got '%s'", expected.ID, result.ID)
	}
	if result.CastLevel != expected.CastLevel {
		t.Errorf("Expected ID '%s', got '%s'", expected.CastLevel, result.CastLevel)
	}
	if result.SpellBaseLevel != expected.SpellBaseLevel {
		t.Errorf("Expected ID '%s', got '%s'", expected.SpellBaseLevel, result.SpellBaseLevel)
	}
	if result.SpellBaseLevel != expected.SpellBaseLevel {
		t.Errorf("Expected ID '%s', got '%s'", expected.SpellBaseLevel, result.SpellBaseLevel)
	}
	if result.Description != expected.Description {
		t.Errorf("Expected ID '%s', got '%s'", expected.Description, result.Description)
	}
	if result.Range != expected.Range {
		t.Errorf("Expected ID '%s', got '%s'", expected.Range, result.Range)
	}
	if result.Duration.Duration != expected.Duration.Duration {
		t.Errorf("Expected ID '%s', got '%s'", expected.Duration.Duration, result.Duration.Duration)
	}
	if result.Duration.Sustained != expected.Duration.Sustained {
		t.Errorf("Expected ID '%t', got '%t'", expected.Duration.Sustained, result.Duration.Sustained)
	}
	if result.Area.Detail != expected.Area.Detail {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Detail, result.Area.Detail)
	}
	if result.Area.Value != expected.Area.Value {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Value, result.Area.Value)
	}
	if result.Area.Type != expected.Area.Type {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Type, result.Area.Type)
	}
	if result.Targets != expected.Targets {
		t.Errorf("Expected ID '%s', got '%s'", expected.Targets, result.Targets)
	}
	if result.Defense.Save != expected.Defense.Save {
		t.Errorf("Expected ID '%s', got '%s'", expected.Defense.Save, result.Defense.Save)
	}
	if result.Defense.Basic != expected.Defense.Basic {
		t.Errorf("Expected ID '%t', got '%t'", expected.Defense.Basic, result.Defense.Basic)
	}
	if result.CastTime != expected.CastTime {
		t.Errorf("Expected ID '%s', got '%s'", expected.CastTime, result.CastTime)
	}
	if result.Rarity != expected.Rarity {
		t.Errorf("Expected Rarity '%s', got '%s'", expected.Rarity, result.Rarity)
	}
	if result.AtWill != expected.AtWill {
		t.Errorf("Expected AtWill '%t', got '%t'", expected.AtWill, result.AtWill)
	}
	if result.SpellCastingBlockLocationID != expected.SpellCastingBlockLocationID {
		t.Fatalf("Expected Spellcasting block ID '%s', got '%s'", expected.SpellCastingBlockLocationID, result.SpellCastingBlockLocationID)
	}
	if result.Uses != expected.Uses {
		t.Fatalf("Expected uses '%s', got '%s'", expected.Uses, result.Uses)
	}
	if result.Ritual != expected.Ritual {
		t.Fatalf("Expected ritualBool '%t', got '%t'", expected.Ritual, result.Ritual)
	}
	if result.RitualData.PrimaryCheck != expected.RitualData.PrimaryCheck {
		t.Fatalf("Expected ritualPrimaryCheck '%s', got '%s'", expected.RitualData.PrimaryCheck, result.RitualData.PrimaryCheck)
	}
	if result.RitualData.SecondaryCheck != expected.RitualData.SecondaryCheck {
		t.Fatalf("Expected ritual secondary check '%s', got '%s'", expected.RitualData.SecondaryCheck, result.RitualData.SecondaryCheck)
	}
	if result.RitualData.SecondaryCasters != expected.RitualData.SecondaryCasters {
		t.Fatalf("Expected ritual secondary casters '%s', got '%s'", expected.RitualData.SecondaryCasters, result.RitualData.SecondaryCasters)
	}
	if len(result.Traits) != len(expected.Traits) {
		t.Fatalf("Expected %d traits, got %d", len(expected.Traits), len(result.Traits))
	}
	for i, trait := range expected.Traits {
		if result.Traits[i] != trait {
			t.Errorf("Expected Trait '%s' at index %d, got '%s'", trait, i, result.Traits[i])
		}
	}
	// Level 5 spontaneous spell, (slots exist in the spellcsting Entry. We just have to tie it to the entry via location.value)
}

func TestIngestRitualInnateSpell(t *testing.T) {
	jsonData := `{
            "_id": "fts4AdQANVel1VuJ",
            "_stats": {
                "compendiumSource": "Compendium.pf2e.spells-srd.Item.tsKnoBuBbKMXkiz5"
            },
            "img": "icons/sundries/scrolls/scroll-writing-tan-grey.webp",
            "name": "Abyssal Pact",
            "sort": 200000,
            "system": {
                "area": null,
                "cost": {
                    "value": ""
                },
                "counteraction": false,
                "damage": {},
                "defense": null,
                "description": {
                    "value": "<p>You call in a favor from another demon whose level is no more than double <em>Abyssal pact's</em> spell rank, two demons whose levels are each at least 2 less than double the spell rank, or three demons whose levels are each at least 3 less than double the spell rank.</p>\n<hr />\n<p><strong>Critical Success</strong> You conjure the demon or demons. They are eager to pursue the task, so they don't ask for a favor.</p>\n<p><strong>Success</strong> You conjure the demon or demons. They are not eager to pursue the task, so they require a favor in return.</p>\n<p><strong>Failure</strong> You don't conjure any demons.</p>\n<p><strong>Critical Failure</strong> The demon or demons are angry that you disturbed them. They appear before you, but they immediately attack you.</p>"
                },
                "duration": {
                    "sustained": false,
                    "value": ""
                },
                "level": {
                    "value": 1
                },
                "location": {
                    "value": null
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": "Pathfinder Bestiary"
                },
                "range": {
                    "value": ""
                },
                "requirements": "",
                "ritual": {
                    "primary": {
                        "check": "Religion (expert; you must be a demon)"
                    },
                    "secondary": {
                        "casters": 0,
                        "checks": ""
                    }
                },
                "rules": [],
                "slug": "abyssal-pact",
                "target": {
                    "value": ""
                },
                "time": {
                    "value": "1 day"
                },
                "traits": {
                    "rarity": "uncommon",
                    "traditions": [],
                    "value": []
                }
            },
            "type": "spell"
        },`

	// location == Null on rituals? Need a different mechanism for those.
	expected := structs.Spell{
		Name:           "Abyssal Pact",
		ID:             "fts4AdQANVel1VuJ",
		CastLevel:      "1",
		SpellBaseLevel: "1",
		Description:    stripHTMLUsingBluemonday("<p>You call in a favor from another demon whose level is no more than double <em>Abyssal pact's</em> spell rank, two demons whose levels are each at least 2 less than double the spell rank, or three demons whose levels are each at least 3 less than double the spell rank.</p>\n<hr />\n<p><strong>Critical Success</strong> You conjure the demon or demons. They are eager to pursue the task, so they don't ask for a favor.</p>\n<p><strong>Success</strong> You conjure the demon or demons. They are not eager to pursue the task, so they require a favor in return.</p>\n<p><strong>Failure</strong> You don't conjure any demons.</p>\n<p><strong>Critical Failure</strong> The demon or demons are angry that you disturbed them. They appear before you, but they immediately attack you.</p>"),
		Range:          "",
		Duration: structs.DurationBlock{
			Sustained: false,
			Duration:  "",
		},
		Area: structs.SpellArea{
			Type:   "",
			Value:  "",
			Detail: "",
		},
		Targets: "",
		Traits:  []string{},
		Defense: structs.DefenseBlock{
			Save:  "",
			Basic: false,
		},
		CastTime:                    "1 day",
		CastRequirements:            "",
		Rarity:                      "uncommon",
		AtWill:                      false,
		SpellCastingBlockLocationID: "",
		Uses:                        "1",
		Ritual:                      true,
		RitualData: structs.RitualData{
			PrimaryCheck:     "Religion (expert; you must be a demon)",
			SecondaryCasters: "0",
			SecondaryCheck:   "",
		}}

	result := ParseSpell(jsonData)

	if result.Name != expected.Name {
		t.Errorf("Expected Name '%s', got '%s'", expected.Name, result.Name)
	}
	if result.ID != expected.ID {
		t.Errorf("Expected ID '%s', got '%s'", expected.ID, result.ID)
	}
	if result.CastLevel != expected.CastLevel {
		t.Errorf("Expected ID '%s', got '%s'", expected.CastLevel, result.CastLevel)
	}
	if result.SpellBaseLevel != expected.SpellBaseLevel {
		t.Errorf("Expected ID '%s', got '%s'", expected.SpellBaseLevel, result.SpellBaseLevel)
	}
	if result.SpellBaseLevel != expected.SpellBaseLevel {
		t.Errorf("Expected ID '%s', got '%s'", expected.SpellBaseLevel, result.SpellBaseLevel)
	}
	if result.Description != expected.Description {
		t.Errorf("Expected ID '%s', got '%s'", expected.Description, result.Description)
	}
	if result.Range != expected.Range {
		t.Errorf("Expected ID '%s', got '%s'", expected.Range, result.Range)
	}
	if result.Duration.Duration != expected.Duration.Duration {
		t.Errorf("Expected ID '%s', got '%s'", expected.Duration.Duration, result.Duration.Duration)
	}
	if result.Duration.Sustained != expected.Duration.Sustained {
		t.Errorf("Expected ID '%t', got '%t'", expected.Duration.Sustained, result.Duration.Sustained)
	}
	if result.Area.Detail != expected.Area.Detail {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Detail, result.Area.Detail)
	}
	if result.Area.Value != expected.Area.Value {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Value, result.Area.Value)
	}
	if result.Area.Type != expected.Area.Type {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Type, result.Area.Type)
	}
	if result.Targets != expected.Targets {
		t.Errorf("Expected ID '%s', got '%s'", expected.Targets, result.Targets)
	}
	if result.Defense.Save != expected.Defense.Save {
		t.Errorf("Expected ID '%s', got '%s'", expected.Defense.Save, result.Defense.Save)
	}
	if result.Defense.Basic != expected.Defense.Basic {
		t.Errorf("Expected ID '%t', got '%t'", expected.Defense.Basic, result.Defense.Basic)
	}
	if result.CastTime != expected.CastTime {
		t.Errorf("Expected ID '%s', got '%s'", expected.CastTime, result.CastTime)
	}
	if result.Rarity != expected.Rarity {
		t.Errorf("Expected Rarity '%s', got '%s'", expected.Rarity, result.Rarity)
	}
	if result.AtWill != expected.AtWill {
		t.Errorf("Expected AtWill '%t', got '%t'", expected.AtWill, result.AtWill)
	}
	if result.SpellCastingBlockLocationID != expected.SpellCastingBlockLocationID {
		t.Fatalf("Expected Spellcasting block ID '%s', got '%s'", expected.SpellCastingBlockLocationID, result.SpellCastingBlockLocationID)
	}
	if result.Uses != expected.Uses {
		t.Fatalf("Expected uses '%s', got '%s'", expected.Uses, result.Uses)
	}
	if result.Ritual != expected.Ritual {
		t.Fatalf("Expected ritualBool '%t', got '%t'", expected.Ritual, result.Ritual)
	}
	if result.RitualData.PrimaryCheck != expected.RitualData.PrimaryCheck {
		t.Fatalf("Expected ritualPrimaryCheck '%s', got '%s'", expected.RitualData.PrimaryCheck, result.RitualData.PrimaryCheck)
	}
	if result.RitualData.SecondaryCheck != expected.RitualData.SecondaryCheck {
		t.Fatalf("Expected ritual secondary check '%s', got '%s'", expected.RitualData.SecondaryCheck, result.RitualData.SecondaryCheck)
	}
	if result.RitualData.SecondaryCasters != expected.RitualData.SecondaryCasters {
		t.Fatalf("Expected ritual secondary casters '%s', got '%s'", expected.RitualData.SecondaryCasters, result.RitualData.SecondaryCasters)
	}
	if len(result.Traits) != len(expected.Traits) {
		t.Fatalf("Expected %d traits, got %d", len(expected.Traits), len(result.Traits))
	}
	for i, trait := range expected.Traits {
		if result.Traits[i] != trait {
			t.Errorf("Expected Trait '%s' at index %d, got '%s'", trait, i, result.Traits[i])
		}
	}
}

func TestIngestPreparedSpell(t *testing.T) {
	jsonData := `{
            "_id": "cgw07bSj0UprtiUE",
            "_stats": {
                "compendiumSource": "Compendium.pf2e.spells-srd.Item.gISYsBFby1TiXfBt"
            },
            "img": "icons/magic/acid/projectile-smoke-glowing.webp",
            "name": "Acid Splash",
            "sort": 2200000,
            "system": {
                "area": null,
                "cost": {
                    "value": ""
                },
                "counteraction": false,
                "damage": {
                    "0": {
                        "applyMod": false,
                        "category": null,
                        "formula": "1d6",
                        "kinds": [
                            "damage"
                        ],
                        "materials": [],
                        "type": "acid"
                    },
                    "gcovwqxwitqchoin": {
                        "applyMod": false,
                        "category": "splash",
                        "formula": "1",
                        "kinds": [
                            "damage"
                        ],
                        "materials": [],
                        "type": "acid"
                    }
                },
                "defense": null,
                "description": {
                    "value": "<p>You splash a glob of acid that splatters your target and nearby creatures. Make a spell attack. If you hit, you deal 1d6 acid damage plus 1 splash acid damage. On a critical success, the target also takes @Damage[(ceil(@item.level/2))[persistent,acid]] damage.</p><hr /><p><strong>Heightened (3rd)</strong> The initial damage increases to 2d6, and the persistent damage increases to 2.</p>\n<p><strong>Heightened (5th)</strong> The initial damage increases to 3d6, the persistent damage increases to 3, and the splash damage increases to 2.</p>\n<p><strong>Heightened (7th)</strong> The initial damage increases to 4d6, the persistent damage increases to 4, and the splash damage increases to 3.</p>\n<p><strong>Heightened (9th)</strong> The initial damage increases to 5d6, the persistent damage increases to 5, and the splash damage increases to 4.</p>"
                },
                "duration": {
                    "sustained": false,
                    "value": ""
                },
                "heightening": {
                    "levels": {
                        "3": {
                            "damage": {
                                "0": {
                                    "applyMod": false,
                                    "category": null,
                                    "formula": "2d6",
                                    "materials": [],
                                    "type": "acid"
                                },
                                "gcovwqxwitqchoin": {
                                    "applyMod": false,
                                    "category": "splash",
                                    "formula": "1",
                                    "materials": [],
                                    "type": "acid"
                                }
                            }
                        },
                        "5": {
                            "damage": {
                                "0": {
                                    "applyMod": false,
                                    "category": null,
                                    "formula": "3d6",
                                    "materials": [],
                                    "type": "acid"
                                },
                                "gcovwqxwitqchoin": {
                                    "applyMod": false,
                                    "category": "splash",
                                    "formula": "2",
                                    "materials": [],
                                    "type": "acid"
                                }
                            }
                        },
                        "7": {
                            "damage": {
                                "0": {
                                    "applyMod": false,
                                    "category": null,
                                    "formula": "4d6",
                                    "materials": [],
                                    "type": "acid"
                                },
                                "gcovwqxwitqchoin": {
                                    "applyMod": false,
                                    "category": "splash",
                                    "formula": "3",
                                    "materials": [],
                                    "type": "acid"
                                }
                            }
                        },
                        "9": {
                            "damage": {
                                "0": {
                                    "applyMod": false,
                                    "category": null,
                                    "formula": "5d6",
                                    "materials": [],
                                    "type": "acid"
                                },
                                "gcovwqxwitqchoin": {
                                    "applyMod": false,
                                    "category": "splash",
                                    "formula": "4",
                                    "materials": [],
                                    "type": "acid"
                                }
                            }
                        }
                    },
                    "type": "fixed"
                },
                "level": {
                    "value": 1
                },
                "location": {
                    "value": "9h6KJeGxzm8rEPaD"
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": "Pathfinder Core Rulebook"
                },
                "range": {
                    "value": "30 feet"
                },
                "requirements": "",
                "rules": [],
                "slug": "acid-splash",
                "target": {
                    "value": "1 creature"
                },
                "time": {
                    "value": "2"
                },
                "traits": {
                    "rarity": "common",
                    "traditions": [
                        "arcane",
                        "primal"
                    ],
                    "value": [
                        "acid",
                        "attack",
                        "cantrip",
                        "concentrate",
                        "manipulate"
                    ]
                }
            },
            "type": "spell"
        },`
	expected := structs.Spell{
		Name:           "Acid Splash",
		ID:             "cgw07bSj0UprtiUE",
		CastLevel:      "1",
		SpellBaseLevel: "1",
		Description:    stripHTMLUsingBluemonday("<p>You splash a glob of acid that splatters your target and nearby creatures. Make a spell attack. If you hit, you deal 1d6 acid damage plus 1 splash acid damage. On a critical success, the target also takes @Damage[(ceil(@item.level/2))[persistent,acid]] damage.</p><hr /><p><strong>Heightened (3rd)</strong> The initial damage increases to 2d6, and the persistent damage increases to 2.</p>\n<p><strong>Heightened (5th)</strong> The initial damage increases to 3d6, the persistent damage increases to 3, and the splash damage increases to 2.</p>\n<p><strong>Heightened (7th)</strong> The initial damage increases to 4d6, the persistent damage increases to 4, and the splash damage increases to 3.</p>\n<p><strong>Heightened (9th)</strong> The initial damage increases to 5d6, the persistent damage increases to 5, and the splash damage increases to 4.</p>"),
		Range:          "30 feet",
		Duration: structs.DurationBlock{
			Sustained: false,
			Duration:  "",
		},
		Area: structs.SpellArea{
			Type:   "",
			Value:  "",
			Detail: "",
		},
		Targets: "1 creature",
		Traits:  []string{"acid", "attack", "cantrip", "concentrate", "manipulate"},
		Defense: structs.DefenseBlock{
			Save:  "",
			Basic: false,
		},
		CastTime:                    "2",
		CastRequirements:            "",
		Rarity:                      "common",
		AtWill:                      false,
		SpellCastingBlockLocationID: "9h6KJeGxzm8rEPaD",
		Uses:                        "1",
		Ritual:                      false,
		RitualData: structs.RitualData{
			PrimaryCheck:     "",
			SecondaryCasters: "",
			SecondaryCheck:   "",
		},
	}

	result := ParseSpell(jsonData)

	if result.Name != expected.Name {
		t.Errorf("Expected Name '%s', got '%s'", expected.Name, result.Name)
	}
	if result.ID != expected.ID {
		t.Errorf("Expected ID '%s', got '%s'", expected.ID, result.ID)
	}
	if result.CastLevel != expected.CastLevel {
		t.Errorf("Expected ID '%s', got '%s'", expected.CastLevel, result.CastLevel)
	}
	if result.SpellBaseLevel != expected.SpellBaseLevel {
		t.Errorf("Expected ID '%s', got '%s'", expected.SpellBaseLevel, result.SpellBaseLevel)
	}
	if result.SpellBaseLevel != expected.SpellBaseLevel {
		t.Errorf("Expected ID '%s', got '%s'", expected.SpellBaseLevel, result.SpellBaseLevel)
	}
	if result.Description != expected.Description {
		t.Errorf("Expected ID '%s', got '%s'", expected.Description, result.Description)
	}
	if result.Range != expected.Range {
		t.Errorf("Expected ID '%s', got '%s'", expected.Range, result.Range)
	}
	if result.Duration.Duration != expected.Duration.Duration {
		t.Errorf("Expected ID '%s', got '%s'", expected.Duration.Duration, result.Duration.Duration)
	}
	if result.Duration.Sustained != expected.Duration.Sustained {
		t.Errorf("Expected ID '%t', got '%t'", expected.Duration.Sustained, result.Duration.Sustained)
	}
	if result.Area.Detail != expected.Area.Detail {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Detail, result.Area.Detail)
	}
	if result.Area.Value != expected.Area.Value {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Value, result.Area.Value)
	}
	if result.Area.Type != expected.Area.Type {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Type, result.Area.Type)
	}
	if result.Targets != expected.Targets {
		t.Errorf("Expected ID '%s', got '%s'", expected.Targets, result.Targets)
	}
	if result.Defense.Save != expected.Defense.Save {
		t.Errorf("Expected ID '%s', got '%s'", expected.Defense.Save, result.Defense.Save)
	}
	if result.Defense.Basic != expected.Defense.Basic {
		t.Errorf("Expected ID '%t', got '%t'", expected.Defense.Basic, result.Defense.Basic)
	}
	if result.CastTime != expected.CastTime {
		t.Errorf("Expected ID '%s', got '%s'", expected.CastTime, result.CastTime)
	}
	if result.Rarity != expected.Rarity {
		t.Errorf("Expected Rarity '%s', got '%s'", expected.Rarity, result.Rarity)
	}
	if result.AtWill != expected.AtWill {
		t.Errorf("Expected AtWill '%t', got '%t'", expected.AtWill, result.AtWill)
	}
	if result.SpellCastingBlockLocationID != expected.SpellCastingBlockLocationID {
		t.Fatalf("Expected Spellcasting block ID '%s', got '%s'", expected.SpellCastingBlockLocationID, result.SpellCastingBlockLocationID)
	}
	if result.Uses != expected.Uses {
		t.Fatalf("Expected uses '%s', got '%s'", expected.Uses, result.Uses)
	}
	if result.Ritual != expected.Ritual {
		t.Fatalf("Expected ritualBool '%t', got '%t'", expected.Ritual, result.Ritual)
	}
	if result.RitualData.PrimaryCheck != expected.RitualData.PrimaryCheck {
		t.Fatalf("Expected ritualPrimaryCheck '%s', got '%s'", expected.RitualData.PrimaryCheck, result.RitualData.PrimaryCheck)
	}
	if result.RitualData.SecondaryCheck != expected.RitualData.SecondaryCheck {
		t.Fatalf("Expected ritual secondary check '%s', got '%s'", expected.RitualData.SecondaryCheck, result.RitualData.SecondaryCheck)
	}
	if result.RitualData.SecondaryCasters != expected.RitualData.SecondaryCasters {
		t.Fatalf("Expected ritual secondary casters '%s', got '%s'", expected.RitualData.SecondaryCasters, result.RitualData.SecondaryCasters)
	}
	if len(result.Traits) != len(expected.Traits) {
		t.Fatalf("Expected %d traits, got %d", len(expected.Traits), len(result.Traits))
	}
	for i, trait := range expected.Traits {
		if result.Traits[i] != trait {
			t.Errorf("Expected Trait '%s' at index %d, got '%s'", trait, i, result.Traits[i])
		}
	}
}

//Level 1 prepared Primal Spell (forest-dragon-adult-spellcaster.json)
// location.value == spellcasting entry Level IS actually the slot it's prepped in.

func TestIngestInnateSpell(t *testing.T) {
	jsonData := `{
            "_id": "kBj0RqQnEELUYiNC",
            "_stats": {
                "compendiumSource": "Compendium.pf2e.spells-srd.Item.4koZzrnMXhhosn0D"
            },
            "img": "systems/pf2e/icons/spells/fear.webp",
            "name": "Fear",
            "sort": 300000,
            "system": {
                "area": null,
                "cost": {
                    "value": ""
                },
                "counteraction": false,
                "damage": {},
                "defense": {
                    "save": {
                        "basic": false,
                        "statistic": "will"
                    }
                },
                "description": {
                    "value": "<p>You plant fear in the target; it must attempt a Will save.</p>\n<hr />\n<p><strong>Critical Success</strong> The target is unaffected.</p>\n<p><strong>Success</strong> The target is @UUID[Compendium.pf2e.conditionitems.Item.Frightened]{Frightened 1}.</p>\n<p><strong>Failure</strong> The target is @UUID[Compendium.pf2e.conditionitems.Item.Frightened]{Frightened 2}.</p>\n<p><strong>Critical Failure</strong> The target is @UUID[Compendium.pf2e.conditionitems.Item.Frightened]{Frightened 3} and @UUID[Compendium.pf2e.conditionitems.Item.Fleeing] for 1 round.</p>\n<hr />\n<p><strong>Heightened (3rd)</strong> You can target up to five creatures.</p>"
                },
                "duration": {
                    "sustained": false,
                    "value": "varies"
                },
                "heightening": {
                    "levels": {
                        "3": {
                            "target": {
                                "value": "5 creatures"
                            }
                        }
                    },
                    "type": "fixed"
                },
                "level": {
                    "value": 1
                },
                "location": {
                    "heightenedLevel": 2,
                    "uses": {
                        "max": 2,
                        "value": 2
                    },
                    "value": "0jNl0jg5W1N5NrTS"
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": "Pathfinder Core Rulebook"
                },
                "range": {
                    "value": "30 feet"
                },
                "requirements": "",
                "rules": [],
                "slug": "fear",
                "target": {
                    "value": "1 creature"
                },
                "time": {
                    "value": "2"
                },
                "traits": {
                    "rarity": "common",
                    "traditions": [
                        "arcane",
                        "divine",
                        "occult",
                        "primal"
                    ],
                    "value": [
                        "concentrate",
                        "emotion",
                        "fear",
                        "manipulate",
                        "mental"
                    ]
                }
            },
            "type": "spell"
        },`
	expected := structs.Spell{
		Name:           "Fear",
		ID:             "kBj0RqQnEELUYiNC",
		CastLevel:      "2",
		SpellBaseLevel: "1",
		Description:    stripHTMLUsingBluemonday("<p>You plant fear in the target; it must attempt a Will save.</p>\n<hr />\n<p><strong>Critical Success</strong> The target is unaffected.</p>\n<p><strong>Success</strong> The target is @UUID[Compendium.pf2e.conditionitems.Item.Frightened]{Frightened 1}.</p>\n<p><strong>Failure</strong> The target is @UUID[Compendium.pf2e.conditionitems.Item.Frightened]{Frightened 2}.</p>\n<p><strong>Critical Failure</strong> The target is @UUID[Compendium.pf2e.conditionitems.Item.Frightened]{Frightened 3} and @UUID[Compendium.pf2e.conditionitems.Item.Fleeing] for 1 round.</p>\n<hr />\n<p><strong>Heightened (3rd)</strong> You can target up to five creatures.</p>"),
		Range:          "30 feet",
		Duration: structs.DurationBlock{
			Sustained: false,
			Duration:  "varies",
		},
		Area: structs.SpellArea{
			Type:   "",
			Value:  "",
			Detail: "",
		},
		Targets: "1 creature",
		Traits:  []string{"concentrate", "emotion", "fear", "manipulate", "mental"},
		Defense: structs.DefenseBlock{
			Save:  "will",
			Basic: false,
		},
		CastTime:                    "2",
		CastRequirements:            "",
		Rarity:                      "common",
		AtWill:                      false,
		SpellCastingBlockLocationID: "0jNl0jg5W1N5NrTS",
		Uses:                        "2",
		Ritual:                      false,
		RitualData: structs.RitualData{
			PrimaryCheck:     "",
			SecondaryCasters: "",
			SecondaryCheck:   "",
		},
	}
	result := ParseSpell(jsonData)

	if result.Name != expected.Name {
		t.Errorf("Expected Name '%s', got '%s'", expected.Name, result.Name)
	}
	if result.ID != expected.ID {
		t.Errorf("Expected ID '%s', got '%s'", expected.ID, result.ID)
	}
	if result.CastLevel != expected.CastLevel {
		t.Errorf("Expected CastLevel '%s', got '%s'", expected.CastLevel, result.CastLevel)
	}
	if result.SpellBaseLevel != expected.SpellBaseLevel {
		t.Errorf("Expected ID '%s', got '%s'", expected.SpellBaseLevel, result.SpellBaseLevel)
	}
	if result.SpellBaseLevel != expected.SpellBaseLevel {
		t.Errorf("Expected ID '%s', got '%s'", expected.SpellBaseLevel, result.SpellBaseLevel)
	}
	if result.Description != expected.Description {
		t.Errorf("Expected ID '%s', got '%s'", expected.Description, result.Description)
	}
	if result.Range != expected.Range {
		t.Errorf("Expected ID '%s', got '%s'", expected.Range, result.Range)
	}
	if result.Duration.Duration != expected.Duration.Duration {
		t.Errorf("Expected ID '%s', got '%s'", expected.Duration.Duration, result.Duration.Duration)
	}
	if result.Duration.Sustained != expected.Duration.Sustained {
		t.Errorf("Expected ID '%t', got '%t'", expected.Duration.Sustained, result.Duration.Sustained)
	}
	if result.Area.Detail != expected.Area.Detail {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Detail, result.Area.Detail)
	}
	if result.Area.Value != expected.Area.Value {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Value, result.Area.Value)
	}
	if result.Area.Type != expected.Area.Type {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Type, result.Area.Type)
	}
	if result.Targets != expected.Targets {
		t.Errorf("Expected ID '%s', got '%s'", expected.Targets, result.Targets)
	}
	if result.Defense.Save != expected.Defense.Save {
		t.Errorf("Expected ID '%s', got '%s'", expected.Defense.Save, result.Defense.Save)
	}
	if result.Defense.Basic != expected.Defense.Basic {
		t.Errorf("Expected ID '%t', got '%t'", expected.Defense.Basic, result.Defense.Basic)
	}
	if result.CastTime != expected.CastTime {
		t.Errorf("Expected ID '%s', got '%s'", expected.CastTime, result.CastTime)
	}
	if result.Rarity != expected.Rarity {
		t.Errorf("Expected Rarity '%s', got '%s'", expected.Rarity, result.Rarity)
	}
	if result.AtWill != expected.AtWill {
		t.Errorf("Expected AtWill '%t', got '%t'", expected.AtWill, result.AtWill)
	}
	if result.SpellCastingBlockLocationID != expected.SpellCastingBlockLocationID {
		t.Fatalf("Expected Spellcasting block ID '%s', got '%s'", expected.SpellCastingBlockLocationID, result.SpellCastingBlockLocationID)
	}
	if result.Uses != expected.Uses {
		t.Fatalf("Expected uses '%s', got '%s'", expected.Uses, result.Uses)
	}
	if result.Ritual != expected.Ritual {
		t.Fatalf("Expected ritualBool '%t', got '%t'", expected.Ritual, result.Ritual)
	}
	if result.RitualData.PrimaryCheck != expected.RitualData.PrimaryCheck {
		t.Fatalf("Expected ritualPrimaryCheck '%s', got '%s'", expected.RitualData.PrimaryCheck, result.RitualData.PrimaryCheck)
	}
	if result.RitualData.SecondaryCheck != expected.RitualData.SecondaryCheck {
		t.Fatalf("Expected ritual secondary check '%s', got '%s'", expected.RitualData.SecondaryCheck, result.RitualData.SecondaryCheck)
	}
	if result.RitualData.SecondaryCasters != expected.RitualData.SecondaryCasters {
		t.Fatalf("Expected ritual secondary casters '%s', got '%s'", expected.RitualData.SecondaryCasters, result.RitualData.SecondaryCasters)
	}
	if len(result.Traits) != len(expected.Traits) {
		t.Fatalf("Expected %d traits, got %d", len(expected.Traits), len(result.Traits))
	}
	for i, trait := range expected.Traits {
		if result.Traits[i] != trait {
			t.Errorf("Expected Trait '%s' at index %d, got '%s'", trait, i, result.Traits[i])
		}
	}
	//Location.value == spellcasting value
	// Heightened Level == Level? what if it says level?
	// Uses needs to be ingested into a use....
	// If no use set then theres 1 use?
}

func TestIngestInnateSpell1Use(t *testing.T) {
	jsonData := `{
            "_id": "6Dv8wIStddSP0cLP",
            "_stats": {
                "compendiumSource": "Compendium.pf2e.spells-srd.Item.3x6eUCm17n6ROzUa"
            },
            "img": "icons/magic/holy/prayer-hands-glowing-yellow-white.webp",
            "name": "Crisis of Faith",
            "sort": 300000,
            "system": {
                "area": null,
                "cost": {
                    "value": ""
                },
                "counteraction": false,
                "damage": {
                    "0": {
                        "category": null,
                        "formula": "6d6",
                        "kinds": [
                            "damage"
                        ],
                        "materials": [],
                        "type": "mental"
                    }
                },
                "defense": {
                    "save": {
                        "basic": false,
                        "statistic": "will"
                    }
                },
                "description": {
                    "value": "<p>You assault the target's faith, riddling the creature with doubt and mental turmoil that deal 6d6 mental damage, or 6d8 mental damage if it can cast divine spells. The effects are determined by its Will save.</p>\n<p>To many deities, casting this spell on a follower of your own deity without significant cause is anathema.</p>\n<hr />\n<p><strong>Critical Success</strong> The target is unaffected.</p>\n<p><strong>Success</strong> The target takes half damage.</p>\n<p><strong>Failure</strong> The target takes full damage; if the target can cast divine spells, it's @UUID[Compendium.pf2e.conditionitems.Item.Stupefied]{Stupefied 1} for 1 round.</p>\n<p><strong>Critical Failure</strong> The target takes double damage, is @UUID[Compendium.pf2e.conditionitems.Item.Stupefied]{Stupefied 1} for 1 round, and can't cast divine spells for 1 round.</p>\n<hr />\n<p><strong>Heightened (+1)</strong> The damage increases by 2d6 (or by 2d8 if the target is a divine spellcaster).</p>"
                },
                "duration": {
                    "sustained": false,
                    "value": ""
                },
                "heightening": {
                    "damage": {
                        "0": "2d6"
                    },
                    "interval": 1,
                    "type": "interval"
                },
                "level": {
                    "value": 3
                },
                "location": {
                    "heightenedLevel": 4,
                    "value": "p3v8D49u0adS76qw"
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": "Pathfinder Core Rulebook"
                },
                "range": {
                    "value": "30 feet"
                },
                "requirements": "",
                "rules": [],
                "slug": "crisis-of-faith",
                "target": {
                    "value": "1 creature"
                },
                "time": {
                    "value": "2"
                },
                "traits": {
                    "rarity": "common",
                    "traditions": [
                        "divine"
                    ],
                    "value": [
                        "concentrate",
                        "manipulate",
                        "mental"
                    ]
                }
            },
            "type": "spell"
        },`
	expected := structs.Spell{
		Name:           "Crisis of Faith",
		ID:             "6Dv8wIStddSP0cLP",
		CastLevel:      "4",
		SpellBaseLevel: "3",
		Description:    stripHTMLUsingBluemonday("<p>You assault the target's faith, riddling the creature with doubt and mental turmoil that deal 6d6 mental damage, or 6d8 mental damage if it can cast divine spells. The effects are determined by its Will save.</p>\n<p>To many deities, casting this spell on a follower of your own deity without significant cause is anathema.</p>\n<hr />\n<p><strong>Critical Success</strong> The target is unaffected.</p>\n<p><strong>Success</strong> The target takes half damage.</p>\n<p><strong>Failure</strong> The target takes full damage; if the target can cast divine spells, it's @UUID[Compendium.pf2e.conditionitems.Item.Stupefied]{Stupefied 1} for 1 round.</p>\n<p><strong>Critical Failure</strong> The target takes double damage, is @UUID[Compendium.pf2e.conditionitems.Item.Stupefied]{Stupefied 1} for 1 round, and can't cast divine spells for 1 round.</p>\n<hr />\n<p><strong>Heightened (+1)</strong> The damage increases by 2d6 (or by 2d8 if the target is a divine spellcaster).</p>"),
		Range:          "30 feet",
		Duration: structs.DurationBlock{
			Sustained: false,
			Duration:  "",
		},
		Area: structs.SpellArea{
			Type:   "",
			Value:  "",
			Detail: "",
		},
		Targets: "1 creature",
		Traits:  []string{"concentrate", "manipulate", "mental"},
		Defense: structs.DefenseBlock{
			Save:  "will",
			Basic: false,
		},
		CastTime:                    "2",
		CastRequirements:            "",
		Rarity:                      "common",
		AtWill:                      false,
		SpellCastingBlockLocationID: "p3v8D49u0adS76qw",
		Uses:                        "1",
		Ritual:                      false,
		RitualData: structs.RitualData{
			PrimaryCheck:     "",
			SecondaryCasters: "",
			SecondaryCheck:   "",
		},
	}
	result := ParseSpell(jsonData)

	if result.Name != expected.Name {
		t.Errorf("Expected Name '%s', got '%s'", expected.Name, result.Name)
	}
	if result.ID != expected.ID {
		t.Errorf("Expected ID '%s', got '%s'", expected.ID, result.ID)
	}
	if result.CastLevel != expected.CastLevel {
		t.Errorf("Expected CastLevel '%s', got '%s'", expected.CastLevel, result.CastLevel)
	}
	if result.SpellBaseLevel != expected.SpellBaseLevel {
		t.Errorf("Expected ID '%s', got '%s'", expected.SpellBaseLevel, result.SpellBaseLevel)
	}
	if result.SpellBaseLevel != expected.SpellBaseLevel {
		t.Errorf("Expected ID '%s', got '%s'", expected.SpellBaseLevel, result.SpellBaseLevel)
	}
	if result.Description != expected.Description {
		t.Errorf("Expected ID '%s', got '%s'", expected.Description, result.Description)
	}
	if result.Range != expected.Range {
		t.Errorf("Expected ID '%s', got '%s'", expected.Range, result.Range)
	}
	if result.Duration.Duration != expected.Duration.Duration {
		t.Errorf("Expected ID '%s', got '%s'", expected.Duration.Duration, result.Duration.Duration)
	}
	if result.Duration.Sustained != expected.Duration.Sustained {
		t.Errorf("Expected ID '%t', got '%t'", expected.Duration.Sustained, result.Duration.Sustained)
	}
	if result.Area.Detail != expected.Area.Detail {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Detail, result.Area.Detail)
	}
	if result.Area.Value != expected.Area.Value {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Value, result.Area.Value)
	}
	if result.Area.Type != expected.Area.Type {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Type, result.Area.Type)
	}
	if result.Targets != expected.Targets {
		t.Errorf("Expected ID '%s', got '%s'", expected.Targets, result.Targets)
	}
	if result.Defense.Save != expected.Defense.Save {
		t.Errorf("Expected ID '%s', got '%s'", expected.Defense.Save, result.Defense.Save)
	}
	if result.Defense.Basic != expected.Defense.Basic {
		t.Errorf("Expected ID '%t', got '%t'", expected.Defense.Basic, result.Defense.Basic)
	}
	if result.CastTime != expected.CastTime {
		t.Errorf("Expected ID '%s', got '%s'", expected.CastTime, result.CastTime)
	}
	if result.Rarity != expected.Rarity {
		t.Errorf("Expected Rarity '%s', got '%s'", expected.Rarity, result.Rarity)
	}
	if result.AtWill != expected.AtWill {
		t.Errorf("Expected AtWill '%t', got '%t'", expected.AtWill, result.AtWill)
	}
	if result.SpellCastingBlockLocationID != expected.SpellCastingBlockLocationID {
		t.Fatalf("Expected Spellcasting block ID '%s', got '%s'", expected.SpellCastingBlockLocationID, result.SpellCastingBlockLocationID)
	}
	if result.Uses != expected.Uses {
		t.Fatalf("Expected uses '%s', got '%s'", expected.Uses, result.Uses)
	}
	if result.Ritual != expected.Ritual {
		t.Fatalf("Expected ritualBool '%t', got '%t'", expected.Ritual, result.Ritual)
	}
	if result.RitualData.PrimaryCheck != expected.RitualData.PrimaryCheck {
		t.Fatalf("Expected ritualPrimaryCheck '%s', got '%s'", expected.RitualData.PrimaryCheck, result.RitualData.PrimaryCheck)
	}
	if result.RitualData.SecondaryCheck != expected.RitualData.SecondaryCheck {
		t.Fatalf("Expected ritual secondary check '%s', got '%s'", expected.RitualData.SecondaryCheck, result.RitualData.SecondaryCheck)
	}
	if result.RitualData.SecondaryCasters != expected.RitualData.SecondaryCasters {
		t.Fatalf("Expected ritual secondary casters '%s', got '%s'", expected.RitualData.SecondaryCasters, result.RitualData.SecondaryCasters)
	}
	if len(result.Traits) != len(expected.Traits) {
		t.Fatalf("Expected %d traits, got %d", len(expected.Traits), len(result.Traits))
	}
	for i, trait := range expected.Traits {
		if result.Traits[i] != trait {
			t.Errorf("Expected Trait '%s' at index %d, got '%s'", trait, i, result.Traits[i])
		}
	}
}

func TestIngestAtWillInnateSpellUse(t *testing.T) {
	jsonData := `        {
            "_id": "ZSNdsMrtDj0biKjg",
            "_stats": {
                "compendiumSource": "Compendium.pf2e.spells-srd.Item.9HpwDN4MYQJnW0LG"
            },
            "img": "systems/pf2e/icons/spells/dispel-magic.webp",
            "name": "Dispel Magic (At Will)",
            "sort": 1200000,
            "system": {
                "area": null,
                "cost": {
                    "value": ""
                },
                "counteraction": true,
                "damage": {},
                "defense": null,
                "description": {
                    "value": "<p>You unravel the magic behind a spell or effect. Attempt a counteract check against the target. If you successfully counteract a magic item, the item becomes a mundane item of its type for 10 minutes. This doesn't change the item's non-magical properties. If the target is an artifact or similar item, you automatically fail.</p>"
                },
                "duration": {
                    "sustained": false,
                    "value": ""
                },
                "level": {
                    "value": 2
                },
                "location": {
                    "heightenedLevel": 8,
                    "value": "IsRnnfl27oJF1UGY"
                },
                "publication": {
                    "license": "ORC",
                    "remaster": true,
                    "title": "Pathfinder Player Core"
                },
                "range": {
                    "value": "120 feet"
                },
                "requirements": "",
                "rules": [],
                "slug": "dispel-magic",
                "target": {
                    "value": "1 spell effect or unattended magic item"
                },
                "time": {
                    "value": "2"
                },
                "traits": {
                    "rarity": "common",
                    "traditions": [
                        "arcane",
                        "divine",
                        "occult",
                        "primal"
                    ],
                    "value": [
                        "concentrate",
                        "manipulate"
                    ]
                }
            },
            "type": "spell"
        },`
	expected := structs.Spell{
		ID:             "ZSNdsMrtDj0biKjg",
		Name:           "Dispel Magic (At Will)",
		CastLevel:      "8",
		SpellBaseLevel: "2",
		Description:    stripHTMLUsingBluemonday("<p>You unravel the magic behind a spell or effect. Attempt a counteract check against the target. If you successfully counteract a magic item, the item becomes a mundane item of its type for 10 minutes. This doesn't change the item's non-magical properties. If the target is an artifact or similar item, you automatically fail.</p>"),
		Range:          "120 feet",
		Duration: structs.DurationBlock{
			Sustained: false,
			Duration:  "",
		},
		Area: structs.SpellArea{
			Type:   "",
			Value:  "",
			Detail: "",
		},
		Targets: "1 spell effect or unattended magic item",
		Traits:  []string{"concentrate", "manipulate"},
		Defense: structs.DefenseBlock{
			Save:  "",
			Basic: false,
		},
		CastTime:                    "2",
		CastRequirements:            "",
		Rarity:                      "common",
		AtWill:                      true,
		SpellCastingBlockLocationID: "IsRnnfl27oJF1UGY",
		Uses:                        "unlimited",
		Ritual:                      false,
		RitualData: structs.RitualData{
			PrimaryCheck:     "",
			SecondaryCasters: "",
			SecondaryCheck:   "",
		},
	}

	result := ParseSpell(jsonData)

	if result.Name != expected.Name {
		t.Errorf("Expected Name '%s', got '%s'", expected.Name, result.Name)
	}
	if result.ID != expected.ID {
		t.Errorf("Expected ID '%s', got '%s'", expected.ID, result.ID)
	}
	if result.CastLevel != expected.CastLevel {
		t.Errorf("Expected CastLevel '%s', got '%s'", expected.CastLevel, result.CastLevel)
	}
	if result.SpellBaseLevel != expected.SpellBaseLevel {
		t.Errorf("Expected ID '%s', got '%s'", expected.SpellBaseLevel, result.SpellBaseLevel)
	}
	if result.SpellBaseLevel != expected.SpellBaseLevel {
		t.Errorf("Expected ID '%s', got '%s'", expected.SpellBaseLevel, result.SpellBaseLevel)
	}
	if result.Description != expected.Description {
		t.Errorf("Expected ID '%s', got '%s'", expected.Description, result.Description)
	}
	if result.Range != expected.Range {
		t.Errorf("Expected ID '%s', got '%s'", expected.Range, result.Range)
	}
	if result.Duration.Duration != expected.Duration.Duration {
		t.Errorf("Expected ID '%s', got '%s'", expected.Duration.Duration, result.Duration.Duration)
	}
	if result.Duration.Sustained != expected.Duration.Sustained {
		t.Errorf("Expected ID '%t', got '%t'", expected.Duration.Sustained, result.Duration.Sustained)
	}
	if result.Area.Detail != expected.Area.Detail {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Detail, result.Area.Detail)
	}
	if result.Area.Value != expected.Area.Value {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Value, result.Area.Value)
	}
	if result.Area.Type != expected.Area.Type {
		t.Errorf("Expected ID '%s', got '%s'", expected.Area.Type, result.Area.Type)
	}
	if result.Targets != expected.Targets {
		t.Errorf("Expected ID '%s', got '%s'", expected.Targets, result.Targets)
	}
	if result.Defense.Save != expected.Defense.Save {
		t.Errorf("Expected ID '%s', got '%s'", expected.Defense.Save, result.Defense.Save)
	}
	if result.Defense.Basic != expected.Defense.Basic {
		t.Errorf("Expected ID '%t', got '%t'", expected.Defense.Basic, result.Defense.Basic)
	}
	if result.CastTime != expected.CastTime {
		t.Errorf("Expected ID '%s', got '%s'", expected.CastTime, result.CastTime)
	}
	if result.Rarity != expected.Rarity {
		t.Errorf("Expected Rarity '%s', got '%s'", expected.Rarity, result.Rarity)
	}
	if result.AtWill != expected.AtWill {
		t.Errorf("Expected AtWill '%t', got '%t'", expected.AtWill, result.AtWill)
	}
	if result.SpellCastingBlockLocationID != expected.SpellCastingBlockLocationID {
		t.Fatalf("Expected Spellcasting block ID '%s', got '%s'", expected.SpellCastingBlockLocationID, result.SpellCastingBlockLocationID)
	}
	if result.Uses != expected.Uses {
		t.Fatalf("Expected uses '%s', got '%s'", expected.Uses, result.Uses)
	}
	if result.Ritual != expected.Ritual {
		t.Fatalf("Expected ritualBool '%t', got '%t'", expected.Ritual, result.Ritual)
	}
	if result.RitualData.PrimaryCheck != expected.RitualData.PrimaryCheck {
		t.Fatalf("Expected ritualPrimaryCheck '%s', got '%s'", expected.RitualData.PrimaryCheck, result.RitualData.PrimaryCheck)
	}
	if result.RitualData.SecondaryCheck != expected.RitualData.SecondaryCheck {
		t.Fatalf("Expected ritual secondary check '%s', got '%s'", expected.RitualData.SecondaryCheck, result.RitualData.SecondaryCheck)
	}
	if result.RitualData.SecondaryCasters != expected.RitualData.SecondaryCasters {
		t.Fatalf("Expected ritual secondary casters '%s', got '%s'", expected.RitualData.SecondaryCasters, result.RitualData.SecondaryCasters)
	}
	if len(result.Traits) != len(expected.Traits) {
		t.Fatalf("Expected %d traits, got %d", len(expected.Traits), len(result.Traits))
	}
	for i, trait := range expected.Traits {
		if result.Traits[i] != trait {
			t.Errorf("Expected Trait '%s' at index %d, got '%s'", trait, i, result.Traits[i])
		}
	}
}

func TestAssignSpellPrepared(t *testing.T) {
	spellData := `{
            "_id": "cgw07bSj0UprtiUE",
            "_stats": {
                "compendiumSource": "Compendium.pf2e.spells-srd.Item.gISYsBFby1TiXfBt"
            },
            "img": "icons/magic/acid/projectile-smoke-glowing.webp",
            "name": "Acid Splash",
            "sort": 2200000,
            "system": {
                "area": null,
                "cost": {
                    "value": ""
                },
                "counteraction": false,
                "damage": {
                    "0": {
                        "applyMod": false,
                        "category": null,
                        "formula": "1d6",
                        "kinds": [
                            "damage"
                        ],
                        "materials": [],
                        "type": "acid"
                    },
                    "gcovwqxwitqchoin": {
                        "applyMod": false,
                        "category": "splash",
                        "formula": "1",
                        "kinds": [
                            "damage"
                        ],
                        "materials": [],
                        "type": "acid"
                    }
                },
                "defense": null,
                "description": {
                    "value": "<p>You splash a glob of acid that splatters your target and nearby creatures. Make a spell attack. If you hit, you deal 1d6 acid damage plus 1 splash acid damage. On a critical success, the target also takes @Damage[(ceil(@item.level/2))[persistent,acid]] damage.</p><hr /><p><strong>Heightened (3rd)</strong> The initial damage increases to 2d6, and the persistent damage increases to 2.</p>\n<p><strong>Heightened (5th)</strong> The initial damage increases to 3d6, the persistent damage increases to 3, and the splash damage increases to 2.</p>\n<p><strong>Heightened (7th)</strong> The initial damage increases to 4d6, the persistent damage increases to 4, and the splash damage increases to 3.</p>\n<p><strong>Heightened (9th)</strong> The initial damage increases to 5d6, the persistent damage increases to 5, and the splash damage increases to 4.</p>"
                },
                "duration": {
                    "sustained": false,
                    "value": ""
                },
                "heightening": {
                    "levels": {
                        "3": {
                            "damage": {
                                "0": {
                                    "applyMod": false,
                                    "category": null,
                                    "formula": "2d6",
                                    "materials": [],
                                    "type": "acid"
                                },
                                "gcovwqxwitqchoin": {
                                    "applyMod": false,
                                    "category": "splash",
                                    "formula": "1",
                                    "materials": [],
                                    "type": "acid"
                                }
                            }
                        },
                        "5": {
                            "damage": {
                                "0": {
                                    "applyMod": false,
                                    "category": null,
                                    "formula": "3d6",
                                    "materials": [],
                                    "type": "acid"
                                },
                                "gcovwqxwitqchoin": {
                                    "applyMod": false,
                                    "category": "splash",
                                    "formula": "2",
                                    "materials": [],
                                    "type": "acid"
                                }
                            }
                        },
                        "7": {
                            "damage": {
                                "0": {
                                    "applyMod": false,
                                    "category": null,
                                    "formula": "4d6",
                                    "materials": [],
                                    "type": "acid"
                                },
                                "gcovwqxwitqchoin": {
                                    "applyMod": false,
                                    "category": "splash",
                                    "formula": "3",
                                    "materials": [],
                                    "type": "acid"
                                }
                            }
                        },
                        "9": {
                            "damage": {
                                "0": {
                                    "applyMod": false,
                                    "category": null,
                                    "formula": "5d6",
                                    "materials": [],
                                    "type": "acid"
                                },
                                "gcovwqxwitqchoin": {
                                    "applyMod": false,
                                    "category": "splash",
                                    "formula": "4",
                                    "materials": [],
                                    "type": "acid"
                                }
                            }
                        }
                    },
                    "type": "fixed"
                },
                "level": {
                    "value": 1
                },
                "location": {
                    "value": "9h6KJeGxzm8rEPaD"
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": "Pathfinder Core Rulebook"
                },
                "range": {
                    "value": "30 feet"
                },
                "requirements": "",
                "rules": [],
                "slug": "acid-splash",
                "target": {
                    "value": "1 creature"
                },
                "time": {
                    "value": "2"
                },
                "traits": {
                    "rarity": "common",
                    "traditions": [
                        "arcane",
                        "primal"
                    ],
                    "value": [
                        "acid",
                        "attack",
                        "cantrip",
                        "concentrate",
                        "manipulate"
                    ]
                }
            },
            "type": "spell"
        },`

	spellCastingBlock := `{
            "_id": "9h6KJeGxzm8rEPaD",
            "img": "systems/pf2e/icons/default-icons/spellcastingEntry.svg",
            "name": "Primal Prepared Spells",
            "sort": 100000,
            "system": {
                "autoHeightenLevel": {
                    "value": 6
                },
                "description": {
                    "value": ""
                },
                "prepared": {
                    "flexible": false,
                    "value": "prepared"
                },
                "proficiency": {
                    "value": 1
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": ""
                },
                "rules": [],
                "showSlotlessLevels": {
                    "value": false
                },
                "slots": {
                    "slot0": {
                        "max": 5,
                        "prepared": [
                            {
                                "id": "cgw07bSj0UprtiUE"
                            },
                            {
                                "id": "GeRqpkpFNtXrmbgm"
                            },
                            {
                                "id": "tLuFR0oqghOXKzbd"
                            },
                            {
                                "id": "wmqu97fbZeHaDCYh"
                            },
                            {
                                "id": "ELMWrZpjcRl1T4RG"
                            }
                        ]
                    },
                    "slot1": {
                        "max": 3,
                        "prepared": [
                            {
                                "id": "K2hzbKGlsnbs4Oim"
                            },
                            {
                                "id": "YfWayh8Vf56Z3brL"
                            },
                            {
                                "id": "ZiYYZgtUKyVmJTXf"
                            }
                        ]
                    },
                    "slot2": {
                        "max": 3,
                        "prepared": [
                            {
                                "id": "lyHhpzUmgixU51K3"
                            },
                            {
                                "id": "JxsY3WYSjn7MwRgz"
                            },
                            {
                                "id": "9YyN3ZnrZrlMGETw"
                            }
                        ]
                    },
                    "slot3": {
                        "max": 3,
                        "prepared": [
                            {
                                "id": "uu8jCMiKsmK3daVq"
                            },
                            {
                                "id": "kKqJb4vg5dRnYkWw"
                            },
                            {
                                "id": "gSRFsZkX8Qu19CEz"
                            }
                        ]
                    },
                    "slot4": {
                        "max": 3,
                        "prepared": [
                            {
                                "id": "Pc8OabeDh0D0QoNn"
                            },
                            {
                                "id": "T6VXVjgqGBXusSVY"
                            },
                            {
                                "id": "VVTdSugZYXwWMIqG"
                            }
                        ]
                    },
                    "slot5": {
                        "max": 3,
                        "prepared": [
                            {
                                "id": "Pr9Ih78tzMSQfxvf"
                            },
                            {
                                "id": "7YdPP01kBJ4BN5CS"
                            },
                            {
                                "id": "D5sHvAzd2vbdfA3E"
                            }
                        ]
                    },
                    "slot6": {
                        "max": 3,
                        "prepared": [
                            {
                                "id": "iUZaBJdkAt5wfkw9"
                            },
                            {
                                "id": "Qc0rR7NFVpIq7lgF"
                            },
                            {
                                "id": "ECGCJIVLGkNeDpoK"
                            }
                        ]
                    }
                },
                "slug": null,
                "spelldc": {
                    "dc": 34,
                    "mod": 0,
                    "value": 28
                },
                "tradition": {
                    "value": "primal"
                }
            },
            "type": "spellcastingEntry"
        }`

	spell2Good := ` {
            "_id": "K2hzbKGlsnbs4Oim",
            "_stats": {
                "compendiumSource": "Compendium.pf2e.spells-srd.Item.WzLKjSw6hsBhuklC"
            },
            "img": "systems/pf2e/icons/spells/create-water.webp",
            "name": "Create Water",
            "sort": 2300000,
            "system": {
                "area": null,
                "cost": {
                    "value": ""
                },
                "counteraction": false,
                "damage": {},
                "defense": null,
                "description": {
                    "value": "<p>As you cup your hands, water begins to flow forth from them. You create 2 gallons of water. If no one drinks it, it evaporates after 1 day.</p>"
                },
                "duration": {
                    "sustained": false,
                    "value": ""
                },
                "level": {
                    "value": 1
                },
                "location": {
                    "value": "9h6KJeGxzm8rEPaD"
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": "Pathfinder Core Rulebook"
                },
                "range": {
                    "value": "0 feet"
                },
                "requirements": "",
                "rules": [],
                "slug": "create-water",
                "target": {
                    "value": ""
                },
                "time": {
                    "value": "2"
                },
                "traits": {
                    "rarity": "common",
                    "traditions": [
                        "arcane",
                        "divine",
                        "primal"
                    ],
                    "value": [
                        "concentrate",
                        "manipulate",
                        "water"
                    ]
                }
            },
            "type": "spell"
        },`
	expectedSpell := ParseSpell(spell2Good)
	spellCasting := ParsePreparedSpellCasting(spellCastingBlock)
	spellList := []structs.Spell{ParseSpell(spellData), expectedSpell}

	demoSpellcasting := structs.SpellCasting{
		PreparedSpellCasting: []structs.PreparedSpellCasting{spellCasting},
	}
	AssignSpell(&spellList, &demoSpellcasting)

	//
	i := 0
	for i < len(demoSpellcasting.PreparedSpellCasting[0].Slots) {
		if demoSpellcasting.PreparedSpellCasting[0].Slots[i].Spell.Name != "Create Water" {
			i += 1
			fmt.Println("still lookin")
		} else if demoSpellcasting.PreparedSpellCasting[0].Slots[i].Spell.Name == "Create Water" {
			fmt.Println("Found it!")
			fmt.Println(demoSpellcasting.PreparedSpellCasting[0].Slots[i].Spell)
			break
		}
	}
	if i == len(demoSpellcasting.PreparedSpellCasting[0].Slots)-1 {
		t.Errorf("Unable to find Create water in prepared spellcasting block")
	}
}

func TestLoadJSON(t *testing.T) {
	data, err := LoadJSON("forest-dragon-adult-spellcaster.json")
	if err != nil {
		t.Errorf("Error on loading. %v", err)
	}
	fmt.Println(data)
}

func TestParseSenses(t *testing.T) {
	data, err := LoadJSON("forest-dragon-adult-spellcaster.json")
	if err != nil {
		t.Errorf("Error on loading. %v", err)
	}

	result := ParseSenses(data)
	expected := []structs.Sense{
		{
			Name: "darkvision",
		},
		{
			Name:   "scent",
			Range:  "60",
			Acuity: "imprecise",
		},
	}
	if result[0].Name != expected[0].Name {
		t.Errorf("Expected Name %s, got %s", expected[0].Name, result[0].Name)
	}
	if result[1].Name != expected[1].Name {
		t.Errorf("Expected Name %s, got %s", expected[1].Name, result[1].Name)
	}
	if result[1].Range != expected[1].Range {
		t.Errorf("Expected Range %s, got %s", expected[1].Range, result[1].Range)
	}
	if result[1].Acuity != expected[1].Acuity {
		t.Errorf("Expected Acuity %s, got %s", expected[1].Acuity, result[1].Acuity)
	}
}

func TestParseSaves(t *testing.T) {
	data, err := LoadJSON("forest-dragon-adult-spellcaster.json")
	if err != nil {
		t.Errorf("Error on loading. %v", err)
	}

	result := ParseSaves(data)
	expected := structs.Saves{
		Fort:       "25",
		FortDetail: "",
		Ref:        "22",
		RefDetail:  "",
		Will:       "27",
		WillDetail: "",
	}
	if result.Fort != expected.Fort {
		t.Errorf("Expected Fort %s, got %s", expected.Fort, result.Fort)
	}
	if result.FortDetail != expected.FortDetail {
		t.Errorf("Expected FortDetail %s, got %s", expected.FortDetail, result.FortDetail)
	}
	if result.Will != expected.Will {
		t.Errorf("Expected Will %s, got %s", expected.Will, result.Will)
	}
	if result.WillDetail != expected.WillDetail {
		t.Errorf("Expected WillDetail %s, got %s", expected.WillDetail, result.WillDetail)
	}
	if result.Ref != expected.Ref {
		t.Errorf("Expected Ref %s, got %s", expected.Ref, result.Ref)
	}
	if result.RefDetail != expected.RefDetail {
		t.Errorf("Expected RefDetail %s, got %s", expected.RefDetail, result.RefDetail)
	}
}

func TestParseItemShield(t *testing.T) {
	jsonData := `{
            "_id": "pU3Y57Kf8Ys0wWfG",
            "_stats": {
                "compendiumSource": "Compendium.pf2e.equipment-srd.Item.ezVp13Uw8cWW08Da"
            },
            "img": "icons/equipment/shield/round-wooden-boss-steel-yellow-blue.webp",
            "name": "Wooden Shield",
            "sort": 600000,
            "system": {
                "acBonus": 2,
                "baseItem": "wooden-shield",
                "bulk": {
                    "value": 1
                },
                "containerId": null,
                "description": {
                    "value": "<p>Though they come in a variety of shapes and sizes, the protection offered by wooden shields comes from the stoutness of their materials. While wooden shields are less expensive than steel shields, they break more easily.</p>\n<table class=\"pf2e\">\n<thead>\n<tr>\n<th>Hardness</th>\n<th>HP</th>\n<th>BT</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td>3</td>\n<td>12</td>\n<td>6</td>\n</tr>\n</tbody>\n</table>"
                },
                "equipped": {
                    "carryType": "held",
                    "handsHeld": 1,
                    "invested": null
                },
                "hardness": 3,
                "hp": {
                    "max": 12,
                    "value": 12
                },
                "level": {
                    "value": 0
                },
                "material": {
                    "grade": null,
                    "type": null
                },
                "price": {
                    "value": {
                        "gp": 1
                    }
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": ""
                },
                "quantity": 1,
                "rules": [],
                "runes": {
                    "reinforcing": 0
                },
                "size": "med",
                "slug": "wooden-shield",
                "speedPenalty": 0,
                "traits": {
                    "integrated": null,
                    "rarity": "common",
                    "value": []
                }
            },
            "type": "shield"
        },`
	expected := structs.Item{
		Name:        "Wooden Shield",
		ID:          "pU3Y57Kf8Ys0wWfG",
		Category:    "",
		Level:       "0",
		Description: stripHTMLUsingBluemonday("<p>Though they come in a variety of shapes and sizes, the protection offered by wooden shields comes from the stoutness of their materials. While wooden shields are less expensive than steel shields, they break more easily.</p>\n<table class=\"pf2e\">\n<thead>\n<tr>\n<th>Hardness</th>\n<th>HP</th>\n<th>BT</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td>3</td>\n<td>12</td>\n<td>6</td>\n</tr>\n</tbody>\n</table>"),
		Price: structs.PriceBlock{
			GP: 1,
		},
		Type:   "shield",
		Traits: []string{},
		Rarity: "common",
		Range:  "",
		Size:   "med",
		Reload: "",
		Bulk:   "1",
	}
	result := ParseItem(jsonData)

	if result.Name != expected.Name {
		t.Errorf("Expected name %s, got %s", expected.Name, result.Name)
	}
	if result.ID != expected.ID {
		t.Errorf("Expected ID %s, got %s", expected.ID, result.ID)
	}
	if result.Category != expected.Category {
		t.Errorf("Expected Category %s, got %s", expected.Category, result.Category)
	}
	if result.Description != expected.Description {
		t.Errorf("Expected Description %s, got %s", expected.Description, result.Description)
	}
	if result.Type != expected.Type {
		t.Errorf("Expected Type %s, got %s", expected.Type, result.Type)
	}
	if result.Rarity != expected.Rarity {
		t.Errorf("Expected Rarity %s, got %s", expected.Rarity, result.Rarity)
	}
	if result.Range != expected.Range {
		t.Errorf("Expected Range %s, got %s", expected.Range, result.Range)
	}
	if result.Reload != expected.Reload {
		t.Errorf("Expected Reload %s, got %s", expected.Reload, result.Reload)
	}
	if result.Price.CP != expected.Price.CP {
		t.Errorf("Expected CP Price %d, got %d", expected.Price.CP, result.Price.CP)
	}
	if result.Price.SP != expected.Price.SP {
		t.Errorf("Expected SP price %d, got %d", expected.Price.SP, result.Price.SP)
	}
	if result.Price.GP != expected.Price.GP {
		t.Errorf("Expected GP price %d, got %d", expected.Price.GP, result.Price.GP)
	}
	if result.Price.Per != expected.Price.Per {
		t.Errorf("Expected price per %d, got %d", expected.Price.Per, result.Price.Per)
	}

	for i, trait := range expected.Traits {
		if result.Traits[i] != trait {
			t.Errorf("Expected Trait '%s' at index %d, got '%s'", trait, i, result.Traits[i])
		}
	}

}

func TestParseItemWeapon(t *testing.T) {
	jsonData := `{
            "_id": "CvqNMEeYJFk8B5Uf",
            "_stats": {
                "compendiumSource": "Compendium.pf2e.equipment-srd.Item.tH5GirEy7YB3ZgCk"
            },
            "img": "systems/pf2e/icons/equipment/weapons/rapier.webp",
            "name": "Rapier",
            "sort": 500000,
            "system": {
                "baseItem": "rapier",
                "bonus": {
                    "value": 0
                },
                "bonusDamage": {
                    "value": 0
                },
                "bulk": {
                    "value": 1
                },
                "category": "martial",
                "containerId": null,
                "damage": {
                    "damageType": "piercing",
                    "dice": 1,
                    "die": "d6"
                },
                "description": {
                    "value": "<p>The rapier is a long and thin piercing blade with a basket hilt. It is prized among many as a dueling weapon.</p>"
                },
                "equipped": {
                    "carryType": "worn",
                    "handsHeld": 0,
                    "invested": null
                },
                "group": "sword",
                "hardness": 0,
                "hp": {
                    "max": 0,
                    "value": 0
                },
                "level": {
                    "value": 0
                },
                "material": {
                    "grade": null,
                    "type": null
                },
                "price": {
                    "value": {
                        "gp": 2
                    }
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": ""
                },
                "quantity": 1,
                "range": null,
                "reload": {
                    "value": ""
                },
                "rules": [],
                "runes": {
                    "potency": 0,
                    "property": [],
                    "striking": 0
                },
                "size": "med",
                "slug": "rapier",
                "splashDamage": {
                    "value": 0
                },
                "traits": {
                    "rarity": "common",
                    "value": [
                        "deadly-d8",
                        "disarm",
                        "finesse"
                    ]
                },
                "usage": {
                    "value": "held-in-one-hand"
                }
            },
            "type": "weapon"
        },`
	expected := structs.Item{
		Name:        "Rapier",
		ID:          "CvqNMEeYJFk8B5Uf",
		Category:    "martial",
		Level:       "0",
		Description: stripHTMLUsingBluemonday("<p>The rapier is a long and thin piercing blade with a basket hilt. It is prized among many as a dueling weapon.</p>"),
		Price: structs.PriceBlock{
			GP: 2,
		},
		Type:   "weapon",
		Traits: []string{"deadly-d8", "disarm", "finesse"},
		Rarity: "common",
		Range:  "",
		Size:   "med",
		Reload: "",
		Bulk:   "1",
	}
	result := ParseItem(jsonData)

	if result.Name != expected.Name {
		t.Errorf("Expected name %s, got %s", expected.Name, result.Name)
	}
	if result.ID != expected.ID {
		t.Errorf("Expected ID %s, got %s", expected.ID, result.ID)
	}
	if result.Category != expected.Category {
		t.Errorf("Expected Category %s, got %s", expected.Category, result.Category)
	}
	if result.Description != expected.Description {
		t.Errorf("Expected Description %s, got %s", expected.Description, result.Description)
	}
	if result.Type != expected.Type {
		t.Errorf("Expected Type %s, got %s", expected.Type, result.Type)
	}
	if result.Rarity != expected.Rarity {
		t.Errorf("Expected Rarity %s, got %s", expected.Rarity, result.Rarity)
	}
	if result.Range != expected.Range {
		t.Errorf("Expected Range %s, got %s", expected.Range, result.Range)
	}
	if result.Reload != expected.Reload {
		t.Errorf("Expected Reload %s, got %s", expected.Reload, result.Reload)
	}
	if result.Price.CP != expected.Price.CP {
		t.Errorf("Expected CP Price %d, got %d", expected.Price.CP, result.Price.CP)
	}
	if result.Price.SP != expected.Price.SP {
		t.Errorf("Expected SP price %d, got %d", expected.Price.SP, result.Price.SP)
	}
	if result.Price.GP != expected.Price.GP {
		t.Errorf("Expected GP price %d, got %d", expected.Price.GP, result.Price.GP)
	}
	if result.Price.Per != expected.Price.Per {
		t.Errorf("Expected price per %d, got %d", expected.Price.Per, result.Price.Per)
	}

	for i, trait := range expected.Traits {
		if result.Traits[i] != trait {
			t.Errorf("Expected Trait '%s' at index %d, got '%s'", trait, i, result.Traits[i])
		}
	}

}

// TODO left off here.
func TestParseItemArmor(t *testing.T) {
	jsonData := `{
            "_id": "CvqNMEeYJFk8B5Uf",
            "_stats": {
                "compendiumSource": "Compendium.pf2e.equipment-srd.Item.tH5GirEy7YB3ZgCk"
            },
            "img": "systems/pf2e/icons/equipment/weapons/rapier.webp",
            "name": "Rapier",
            "sort": 500000,
            "system": {
                "baseItem": "rapier",
                "bonus": {
                    "value": 0
                },
                "bonusDamage": {
                    "value": 0
                },
                "bulk": {
                    "value": 1
                },
                "category": "martial",
                "containerId": null,
                "damage": {
                    "damageType": "piercing",
                    "dice": 1,
                    "die": "d6"
                },
                "description": {
                    "value": "<p>The rapier is a long and thin piercing blade with a basket hilt. It is prized among many as a dueling weapon.</p>"
                },
                "equipped": {
                    "carryType": "worn",
                    "handsHeld": 0,
                    "invested": null
                },
                "group": "sword",
                "hardness": 0,
                "hp": {
                    "max": 0,
                    "value": 0
                },
                "level": {
                    "value": 0
                },
                "material": {
                    "grade": null,
                    "type": null
                },
                "price": {
                    "value": {
                        "gp": 2
                    }
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": ""
                },
                "quantity": 1,
                "range": null,
                "reload": {
                    "value": ""
                },
                "rules": [],
                "runes": {
                    "potency": 0,
                    "property": [],
                    "striking": 0
                },
                "size": "med",
                "slug": "rapier",
                "splashDamage": {
                    "value": 0
                },
                "traits": {
                    "rarity": "common",
                    "value": [
                        "deadly-d8",
                        "disarm",
                        "finesse"
                    ]
                },
                "usage": {
                    "value": "held-in-one-hand"
                }
            },
            "type": "weapon"
        },`
	expected := structs.Item{
		Name:        "Rapier",
		ID:          "CvqNMEeYJFk8B5Uf",
		Category:    "martial",
		Level:       "0",
		Description: stripHTMLUsingBluemonday("<p>The rapier is a long and thin piercing blade with a basket hilt. It is prized among many as a dueling weapon.</p>"),
		Price: structs.PriceBlock{
			GP: 2,
		},
		Type:   "weapon",
		Traits: []string{"deadly-d8", "disarm", "finesse"},
		Rarity: "common",
		Range:  "",
		Size:   "med",
		Reload: "",
		Bulk:   "1",
	}
	result := ParseItem(jsonData)

	if result.Name != expected.Name {
		t.Errorf("Expected name %s, got %s", expected.Name, result.Name)
	}
	if result.ID != expected.ID {
		t.Errorf("Expected ID %s, got %s", expected.ID, result.ID)
	}
	if result.Category != expected.Category {
		t.Errorf("Expected Category %s, got %s", expected.Category, result.Category)
	}
	if result.Description != expected.Description {
		t.Errorf("Expected Description %s, got %s", expected.Description, result.Description)
	}
	if result.Type != expected.Type {
		t.Errorf("Expected Type %s, got %s", expected.Type, result.Type)
	}
	if result.Rarity != expected.Rarity {
		t.Errorf("Expected Rarity %s, got %s", expected.Rarity, result.Rarity)
	}
	if result.Range != expected.Range {
		t.Errorf("Expected Range %s, got %s", expected.Range, result.Range)
	}
	if result.Reload != expected.Reload {
		t.Errorf("Expected Reload %s, got %s", expected.Reload, result.Reload)
	}
	if result.Price.CP != expected.Price.CP {
		t.Errorf("Expected CP Price %d, got %d", expected.Price.CP, result.Price.CP)
	}
	if result.Price.SP != expected.Price.SP {
		t.Errorf("Expected SP price %d, got %d", expected.Price.SP, result.Price.SP)
	}
	if result.Price.GP != expected.Price.GP {
		t.Errorf("Expected GP price %d, got %d", expected.Price.GP, result.Price.GP)
	}
	if result.Price.Per != expected.Price.Per {
		t.Errorf("Expected price per %d, got %d", expected.Price.Per, result.Price.Per)
	}

	for i, trait := range expected.Traits {
		if result.Traits[i] != trait {
			t.Errorf("Expected Trait '%s' at index %d, got '%s'", trait, i, result.Traits[i])
		}
	}

}
func TestParseItemTreasure(t *testing.T) {
	jsonData := `{
            "_id": "txASc5iIQvLV4Nxv",
            "img": "systems/pf2e/icons/equipment/worn-items/other-worn-items/aluum-charm.webp",
            "name": "Bejeweled Necklace featuring a Porpoise",
            "sort": 1700000,
            "system": {
                "baseItem": null,
                "bulk": {
                    "value": 0
                },
                "containerId": null,
                "description": {
                    "value": "<p>Ayla, My Beloved</p>"
                },
                "equipped": {
                    "carryType": "worn",
                    "handsHeld": 0
                },
                "hardness": 0,
                "hp": {
                    "max": 0,
                    "value": 0
                },
                "level": {
                    "value": 0
                },
                "material": {
                    "grade": null,
                    "type": null
                },
                "price": {
                    "value": {
                        "gp": 10
                    }
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": ""
                },
                "quantity": 1,
                "rules": [],
                "size": "med",
                "slug": null,
                "traits": {
                    "rarity": "common",
                    "value": []
                }
            },
            "type": "treasure"
        },`
	expected := structs.Item{
		Name:        "Bejeweled Necklace featuring a Porpoise",
		ID:          "txASc5iIQvLV4Nxv",
		Category:    "",
		Level:       "0",
		Description: stripHTMLUsingBluemonday("<p>Ayla, My Beloved</p>"),
		Price: structs.PriceBlock{
			GP: 10,
		},
		Type:   "treasure",
		Traits: []string{},
		Rarity: "common",
		Range:  "",
		Size:   "med",
		Reload: "",
		Bulk:   "0",
	}
	result := ParseItem(jsonData)

	if result.Name != expected.Name {
		t.Errorf("Expected name %s, got %s", expected.Name, result.Name)
	}
	if result.ID != expected.ID {
		t.Errorf("Expected ID %s, got %s", expected.ID, result.ID)
	}
	if result.Category != expected.Category {
		t.Errorf("Expected Category %s, got %s", expected.Category, result.Category)
	}
	if result.Description != expected.Description {
		t.Errorf("Expected Description %s, got %s", expected.Description, result.Description)
	}
	if result.Type != expected.Type {
		t.Errorf("Expected Type %s, got %s", expected.Type, result.Type)
	}
	if result.Rarity != expected.Rarity {
		t.Errorf("Expected Rarity %s, got %s", expected.Rarity, result.Rarity)
	}
	if result.Range != expected.Range {
		t.Errorf("Expected Range %s, got %s", expected.Range, result.Range)
	}
	if result.Reload != expected.Reload {
		t.Errorf("Expected Reload %s, got %s", expected.Reload, result.Reload)
	}
	if result.Price.CP != expected.Price.CP {
		t.Errorf("Expected CP Price %d, got %d", expected.Price.CP, result.Price.CP)
	}
	if result.Price.SP != expected.Price.SP {
		t.Errorf("Expected SP price %d, got %d", expected.Price.SP, result.Price.SP)
	}
	if result.Price.GP != expected.Price.GP {
		t.Errorf("Expected GP price %d, got %d", expected.Price.GP, result.Price.GP)
	}
	if result.Price.Per != expected.Price.Per {
		t.Errorf("Expected price per %d, got %d", expected.Price.Per, result.Price.Per)
	}

	for i, trait := range expected.Traits {
		if result.Traits[i] != trait {
			t.Errorf("Expected Trait '%s' at index %d, got '%s'", trait, i, result.Traits[i])
		}
	}

}

func TestParseItemConsumable(t *testing.T) {
	jsonData := `{
            "_id": "DSpF3QNsGXCQO5Re",
            "_stats": {
                "compendiumSource": "Compendium.pf2e.equipment-srd.Item.w2ENw2VMPcsbif8g"
            },
            "img": "systems/pf2e/icons/equipment/weapons/arrows.webp",
            "name": "Arrows",
            "sort": 500000,
            "system": {
                "baseItem": null,
                "bulk": {
                    "value": 0.1
                },
                "category": "ammo",
                "containerId": null,
                "damage": null,
                "description": {
                    "value": "<p>These projectiles are the ammunition for bows. The shaft of an arrow is made of wood. It is stabilized in flight by fletching at one end and bears a metal head on the other.</p>"
                },
                "equipped": {
                    "carryType": "held",
                    "handsHeld": 1
                },
                "hardness": 0,
                "hp": {
                    "max": 0,
                    "value": 0
                },
                "level": {
                    "value": 0
                },
                "material": {
                    "grade": null,
                    "type": null
                },
                "price": {
                    "per": 10,
                    "value": {
                        "sp": 1
                    }
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": "Pathfinder Core Rulebook"
                },
                "quantity": 20,
                "rules": [],
                "size": "med",
                "slug": "arrows",
                "stackGroup": "arrows",
                "traits": {
                    "rarity": "common",
                    "value": []
                },
                "usage": {
                    "value": "held-in-one-hand"
                },
                "uses": {
                    "autoDestroy": true,
                    "max": 1,
                    "value": 1
                }
            },
            "type": "consumable"
        },`
	expected := structs.Item{
		Name:        "Arrows",
		ID:          "DSpF3QNsGXCQO5Re",
		Category:    "ammo",
		Level:       "0",
		Description: stripHTMLUsingBluemonday("<p>These projectiles are the ammunition for bows. The shaft of an arrow is made of wood. It is stabilized in flight by fletching at one end and bears a metal head on the other.</p>"),
		Price: structs.PriceBlock{
			SP:  1,
			Per: 10,
		},
		Type:     "consumable",
		Traits:   []string{},
		Rarity:   "common",
		Range:    "",
		Size:     "med",
		Reload:   "",
		Bulk:     "0.1",
		Quantity: "20",
	}
	result := ParseItem(jsonData)

	if result.Name != expected.Name {
		t.Errorf("Expected name %s, got %s", expected.Name, result.Name)
	}
	if result.ID != expected.ID {
		t.Errorf("Expected ID %s, got %s", expected.ID, result.ID)
	}
	if result.Category != expected.Category {
		t.Errorf("Expected Category %s, got %s", expected.Category, result.Category)
	}
	if result.Description != expected.Description {
		t.Errorf("Expected Description %s, got %s", expected.Description, result.Description)
	}
	if result.Type != expected.Type {
		t.Errorf("Expected Type %s, got %s", expected.Type, result.Type)
	}
	if result.Rarity != expected.Rarity {
		t.Errorf("Expected Rarity %s, got %s", expected.Rarity, result.Rarity)
	}
	if result.Range != expected.Range {
		t.Errorf("Expected Range %s, got %s", expected.Range, result.Range)
	}
	if result.Reload != expected.Reload {
		t.Errorf("Expected Reload %s, got %s", expected.Reload, result.Reload)
	}
	if result.Price.CP != expected.Price.CP {
		t.Errorf("Expected CP Price %d, got %d", expected.Price.CP, result.Price.CP)
	}
	if result.Price.SP != expected.Price.SP {
		t.Errorf("Expected SP price %d, got %d", expected.Price.SP, result.Price.SP)
	}
	if result.Price.GP != expected.Price.GP {
		t.Errorf("Expected GP price %d, got %d", expected.Price.GP, result.Price.GP)
	}
	if result.Price.Per != expected.Price.Per {
		t.Errorf("Expected price per %d, got %d", expected.Price.Per, result.Price.Per)
	}
	if result.Quantity != expected.Quantity {
		t.Errorf("Expected Quantity %s, got %s", expected.Quantity, result.Quantity)
	}

	for i, trait := range expected.Traits {
		if result.Traits[i] != trait {
			t.Errorf("Expected Trait '%s' at index %d, got '%s'", trait, i, result.Traits[i])
		}
	}

}

func TestParseItemBackPack(t *testing.T) {
	jsonData := `{
            "_id": "hetb6HQzsfpikrYo",
            "_stats": {
                "compendiumSource": "Compendium.pf2e.equipment-srd.Item.jaEEvuQ32GjAa8jy"
            },
            "img": "systems/pf2e/icons/equipment/held-items/bag-of-holding.webp",
            "name": "Bag of Holding (Type I)",
            "sort": 1000000,
            "system": {
                "baseItem": null,
                "bulk": {
                    "capacity": 25,
                    "heldOrStowed": 1,
                    "ignored": 25,
                    "value": 1
                },
                "collapsed": false,
                "containerId": null,
                "description": {
                    "value": "<p>Though it appears to be a cloth sack decorated with panels of richly colored silk or stylish embroidery, a <em>bag of holding</em> opens into an extradimensional space larger than its outside dimensions. The Bulk held inside the bag doesn't change the Bulk of the <em>bag of holding</em> itself. The amount of Bulk the bag's extradimensional space can hold depends on its type.</p>\n<p>You can Interact with the <em>bag of holding</em> to put items in or remove them just like a mundane sack. Though the bag can hold a great amount of material, an object still needs to be able to fit through the opening of the sack to be stored inside.</p>\n<p>If the bag is overloaded or broken, it ruptures and is ruined, causing the items inside to be lost forever. If it's turned inside out, the items inside spill out unharmed, but the bag must be put right before it can be used again. A living creature placed inside the bag has enough air for 10 minutes before it begins to suffocate, and it can attempt to Escape against a DC of 13. An item inside the bag provides no benefits unless it's retrieved first. An item in the bag can't be detected by magic that detects only things on the same plane.</p>\n<p><strong>Capacity</strong> 25 Bulk</p>"
                },
                "equipped": {
                    "carryType": "worn",
                    "handsHeld": 0,
                    "invested": null
                },
                "hardness": 0,
                "hp": {
                    "max": 0,
                    "value": 0
                },
                "level": {
                    "value": 4
                },
                "material": {
                    "grade": null,
                    "type": null
                },
                "price": {
                    "value": {
                        "gp": 75
                    }
                },
                "publication": {
                    "license": "OGL",
                    "remaster": false,
                    "title": "Pathfinder Gamemastery Guide"
                },
                "quantity": 1,
                "rules": [],
                "size": "med",
                "slug": "bag-of-holding-type-i",
                "stowing": true,
                "traits": {
                    "rarity": "common",
                    "value": [
                        "extradimensional",
                        "magical"
                    ]
                },
                "usage": {
                    "value": "held-in-two-hands"
                }
            },
            "type": "backpack"
        },`

	expected := structs.Item{
		Name:        "Bag of Holding (Type I)",
		ID:          "hetb6HQzsfpikrYo",
		Category:    "",
		Level:       "4",
		Description: stripHTMLUsingBluemonday("<p>Though it appears to be a cloth sack decorated with panels of richly colored silk or stylish embroidery, a <em>bag of holding</em> opens into an extradimensional space larger than its outside dimensions. The Bulk held inside the bag doesn't change the Bulk of the <em>bag of holding</em> itself. The amount of Bulk the bag's extradimensional space can hold depends on its type.</p>\n<p>You can Interact with the <em>bag of holding</em> to put items in or remove them just like a mundane sack. Though the bag can hold a great amount of material, an object still needs to be able to fit through the opening of the sack to be stored inside.</p>\n<p>If the bag is overloaded or broken, it ruptures and is ruined, causing the items inside to be lost forever. If it's turned inside out, the items inside spill out unharmed, but the bag must be put right before it can be used again. A living creature placed inside the bag has enough air for 10 minutes before it begins to suffocate, and it can attempt to Escape against a DC of 13. An item inside the bag provides no benefits unless it's retrieved first. An item in the bag can't be detected by magic that detects only things on the same plane.</p>\n<p><strong>Capacity</strong> 25 Bulk</p>"),
		Price: structs.PriceBlock{
			GP: 75,
		},
		Type:     "backpack",
		Traits:   []string{"extradimensional", "magical"},
		Rarity:   "common",
		Range:    "",
		Size:     "med",
		Reload:   "",
		Bulk:     "1",
		Quantity: "1",
	}
	result := ParseItem(jsonData)

	if result.Name != expected.Name {
		t.Errorf("Expected name %s, got %s", expected.Name, result.Name)
	}
	if result.ID != expected.ID {
		t.Errorf("Expected ID %s, got %s", expected.ID, result.ID)
	}
	if result.Category != expected.Category {
		t.Errorf("Expected Category %s, got %s", expected.Category, result.Category)
	}
	if result.Description != expected.Description {
		t.Errorf("Expected Description %s, got %s", expected.Description, result.Description)
	}
	if result.Type != expected.Type {
		t.Errorf("Expected Type %s, got %s", expected.Type, result.Type)
	}
	if result.Rarity != expected.Rarity {
		t.Errorf("Expected Rarity %s, got %s", expected.Rarity, result.Rarity)
	}
	if result.Range != expected.Range {
		t.Errorf("Expected Range %s, got %s", expected.Range, result.Range)
	}
	if result.Reload != expected.Reload {
		t.Errorf("Expected Reload %s, got %s", expected.Reload, result.Reload)
	}
	if result.Price.CP != expected.Price.CP {
		t.Errorf("Expected CP Price %d, got %d", expected.Price.CP, result.Price.CP)
	}
	if result.Price.SP != expected.Price.SP {
		t.Errorf("Expected SP price %d, got %d", expected.Price.SP, result.Price.SP)
	}
	if result.Price.GP != expected.Price.GP {
		t.Errorf("Expected GP price %d, got %d", expected.Price.GP, result.Price.GP)
	}
	if result.Price.Per != expected.Price.Per {
		t.Errorf("Expected price per %d, got %d", expected.Price.Per, result.Price.Per)
	}
	if result.Quantity != expected.Quantity {
		t.Errorf("Expected Quantity %s, got %s", expected.Quantity, result.Quantity)
	}

	for i, trait := range expected.Traits {
		if result.Traits[i] != trait {
			t.Errorf("Expected Trait '%s' at index %d, got '%s'", trait, i, result.Traits[i])
		}
	}
}
