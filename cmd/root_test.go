// Package cmd_test contains unit tests for the command package.
package cmd

import (
	"testing"
)

// TestParseCSV tests the parseCSV function.
func TestParseCSV(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"title,date", []string{"title", "date"}},
		{" draft , category ", []string{"draft", "category"}},
		{"", nil},
		{",title,date,", []string{"", "title", "date", ""}},
		{"title,  date", []string{"title", "date"}},
	}

	for _, tt := range tests {
		got := parseCSV(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("parseCSV(%q) = %v; want %v", tt.input, got, tt.want)
		}
	}
}
