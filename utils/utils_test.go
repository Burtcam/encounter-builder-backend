package utils

import "testing"

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
