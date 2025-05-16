package utils

import "testing"

func TestGetXpBudget(t *testing.T) {
	tests := []struct {
		difficulty string
		pSize      int
		level      int
		expected   int
	}{
		{"easy", 1, 1, 100},
		{"medium", 2, 2, 200},
		{"hard", 3, 3, 300},
	}

	for _, test := range tests {
		result := getXpBudget(test.difficulty, test.pSize, test.level)
		if result != test.expected {
			t.Errorf("getXpBudget(%s, %d, %d) = %d; expected %d", test.difficulty, test.pSize, test.level, result, test.expected)
		}
	}
}
