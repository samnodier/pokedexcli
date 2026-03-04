package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Unit Testing   in Go",
			expected: []string{"unit", "testing", "in", "go"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput failed. Slices length do not match. \nGot: %d, Expected: %d", len(actual), len(c.expected))
		}
		// Check the length of the actual slice against the expected slice
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice and if they don't match return error
			if word != expectedWord {
				t.Errorf("cleanInput failed. Words don't match.\nGot: %s, Expected: %s", word, expectedWord)
			}
		}
	}
}
