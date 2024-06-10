package main

import (
	"testing"
)

func TestParseFields(t *testing.T) {
	tests := []struct {
		input    string
		expected []int
	}{
		{"1,2,3", []int{0, 1, 2}},
		{"4,5,6", []int{3, 4, 5}},
		{"1", []int{0}},
	}

	for _, test := range tests {
		result, err := parseFields(test.input)
		if err != nil {
			t.Errorf("parseFields(%s) returned an error: %v", test.input, err)
		}

		if len(result) != len(test.expected) {
			t.Errorf("parseFields(%s) returned %d fields, expected %d", test.input, len(result), len(test.expected))
		}

		for i := range result {
			if result[i] != test.expected[i] {
				t.Errorf("parseFields(%s) returned %v, expected %v", test.input, result, test.expected)
				break
			}
		}
	}
}

func TestProcessLine(t *testing.T) {
	args := &Args{
		f: "1,2,3",
		d: "\t",
		s: true,
	}

	fields := []int{0, 1, 2}

	tests := []struct {
		line     string
		expected string
	}{
		{"a\tb\tc", "a\tb\tc"},
		{"a\tb", ""},
		{"a", ""},
	}

	for _, test := range tests {
		result := processLine(test.line, args, fields)
		if result != test.expected {
			t.Errorf("processLine(%s) returned %s, expected %s", test.line, result, test.expected)
		}
	}
}
