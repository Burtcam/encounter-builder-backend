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
	result := gjson.Parse(jsonData)

	// Call parsePassives function
	passive := ParsePassives(result)

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
