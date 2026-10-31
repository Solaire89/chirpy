package main

import "testing"

func TestCensorWord(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "I've had it with these fornax snakes in my code!",
			expected: "I've had it with these **** snakes in my code!",
		},
	}
	for _, c := range cases {
		actual := censorBadWord(c.input)
		if actual != c.expected {
			t.Errorf("Sentences don't match: %s vs %s", actual, c.expected)
		}
	}
}
