package main

import (
	"reflect"
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
			input:    "",
			expected: []string{},
		},
		{
			input:    "a   b   c",
			expected: []string{"a", "b", "c"},
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			actual := cleanInput(c.input)

			if !reflect.DeepEqual(actual, c.expected) {
				t.Fatalf("input %q: expected %q, got %q", c.input, c.expected, actual)
			}
		})
	}
}
