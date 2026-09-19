package command

import (
	"testing"
)

func TestSortedReplacementKeys(t *testing.T) {
	tests := []struct {
		name         string
		replacements map[string]string
		expected     []string
	}{
		{
			name:         "empty map",
			replacements: map[string]string{},
			expected:     []string{},
		},
		{
			name: "longest key first",
			replacements: map[string]string{
				"id":      "1",
				"user_id": "2",
			},
			expected: []string{"user_id", "id"},
		},
		{
			name: "same length keys lexicographic",
			replacements: map[string]string{
				"beta":  "b",
				"alpha": "a",
			},
			expected: []string{"alpha", "beta"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys := sortedKeys(tt.replacements)
			if len(keys) != len(tt.expected) {
				t.Fatalf("expected %d keys, got %d", len(tt.expected), len(keys))
			}
			for i, k := range keys {
				if k != tt.expected[i] {
					t.Errorf("position %d: expected %q, got %q", i, tt.expected[i], k)
				}
			}
		})
	}
}
