// Package helpers_test contains unit tests for the helpers package.
package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// TestSplitFrontmatter tests the SplitFrontmatter function.
func TestSplitFrontmatter(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		delimiter string
	}{
		{"YAML frontmatter", "---\ntitle: Test\n---\nBody text", "---"},
		{"TOML frontmatter", "+++\ntitle = 'Test'\n+++\nBody text", "+++"},
		{"JSON frontmatter", "{\n \"title\": \"Test\"\n}\nBody text", "{"},
		{"No frontmatter", "Body text without fm", ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			delim, fm, body := SplitFrontmatter([]byte(c.input))
			if delim != c.delimiter {
				t.Errorf("expected delimiter %q, got %q", c.delimiter, delim)
			}
			if c.delimiter != "" && len(fm) == 0 {
				t.Errorf("expected frontmatter block, got empty")
			}
			if len(body) == 0 {
				t.Errorf("expected body text, got empty")
			}
		})
	}
}

// TestIsMarkdownFile tests the IsMarkdownFile function.
func TestIsMarkdownFile(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"file.md", true},
		{"file.txt", false},
		{"FILE.MD", false},
		{"FILE.TXT", false},
		{"path/to/file.md", true},
		{"path/to/file.txt", false},
	}

	for _, tt := range tests {
		actual := IsMarkdownFile(tt.path)
		if actual != tt.expected {
			t.Errorf("IsMarkdownFile(%q) = %v; want %v", tt.path, actual, tt.expected)
		}
	}
}

// TestMarshalFrontmatter tests the MarshalFrontmatter function.
func TestMarshalFrontmatter(t *testing.T) {
	front := map[string]interface{}{
		"title": "Test Title",
		"draft": true,
		"tags":  []string{"test", "example"},
	}

	tests := []struct {
		delimiter string
	}{
		{"---"},
		{"+++"},
		{"{"},
	}

	for _, tt := range tests {
		marshaled, err := MarshalFrontmatter(tt.delimiter, front)
		if err != nil {
			t.Fatalf("MarshalFrontmatter(%q) error: %v", tt.delimiter, err)
		}

		// Check for inline arrays in TOML format
		if tt.delimiter == "+++" {
			if !strings.Contains(string(marshaled), "tags = [") || strings.Contains(string(marshaled), "tags = [\n") {
				t.Errorf("TOML arrays should be on a single line, got: %s", string(marshaled))
			}
		}

		// Basic verification that unmarshaling works
		var unmarshaled map[string]interface{}
		switch tt.delimiter {
		case "---":
			err = yaml.Unmarshal(marshaled, &unmarshaled)
		case "+++":
			err = toml.Unmarshal(marshaled, &unmarshaled)
		case "{":
			err = json.Unmarshal(marshaled, &unmarshaled)
		}
		if err != nil {
			t.Fatalf("Failed to unmarshal marshaled frontmatter: %v", err)
		}

		// Check that all keys exist
		for k := range front {
			if _, exists := unmarshaled[k]; !exists {
				t.Errorf("Key %q missing from unmarshaled frontmatter", k)
			}
		}
	}
}

// TestFlattenToStrings tests the FlattenToStrings function.
func TestFlattenToStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []string
	}{
		{
			name:     "Slice of interfaces",
			input:    []interface{}{1, "two", 3.0},
			expected: []string{"1", "two", "3"},
		},
		{
			name:     "Slice of strings",
			input:    []string{"one", "two", "three"},
			expected: []string{"one", "two", "three"},
		},
		{
			name:     "Empty slice",
			input:    []interface{}{},
			expected: []string{},
		},
		{
			name:     "Nil slice",
			input:    nil,
			expected: []string{},
		},
		{
			name:     "Single string",
			input:    "hello",
			expected: []string{"hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := flattenToStrings(tt.input)
			if tt.input == nil && len(actual) != 0 {
				t.Errorf("flattenToStrings(%v) = %v, want %v", tt.input, actual, tt.expected)
			}
			if tt.input != nil && !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("flattenToStrings(%v) = %v, want %v", tt.input, actual, tt.expected)
			}
		})
	}
}

// TestParseSet tests the ParseSet function.
func TestParseSet(t *testing.T) {
	tests := []struct {
		input string
		key   string
		value interface{}
	}{
		{"draft=true", "draft", true},
		{"published=false", "published", false},
		{"title=My Title", "title", "My Title"},
	}

	for _, tt := range tests {
		k, v := ParseSet(tt.input)
		if k != tt.key {
			t.Errorf("key mismatch: want %q got %q", tt.key, k)
		}
		if !reflect.DeepEqual(v, tt.value) {
			t.Errorf("value mismatch: want %v got %v", tt.value, v)
		}
	}
}

// TestCheckCondition tests the CheckCondition function.
func TestCheckCondition(t *testing.T) {
	front := map[string]interface{}{
		"draft": true,
		"tags":  []interface{}{"beta", 123, "release"},
		"date":  "2023-01-01",
	}

	tests := []struct {
		cond  string
		match bool
	}{
		{"draft=true", true},
		{"draft=false", false},
		{"tags contains 'beta'", true},
		{"tags contains 123", true},
		{"tags contains 'release'", true},
		{"tags contains 'missing'", false},
		{"date<2024-01-01", true},
	}

	for _, tt := range tests {
		match := CheckCondition(front, tt.cond)
		if match != tt.match {
			t.Errorf("CheckCondition(%q) = %v; want %v", tt.cond, match, tt.match)
		}
	}

	testFiles := []struct {
		path      string
		delimiter string
	}{
		{"testdata/test.yaml", "---"},
		{"testdata/test.toml", "+++"},
		{"testdata/test.json", "{"},
	}

	for _, tf := range testFiles {
		t.Run("TestFile_"+tf.path, func(t *testing.T) {
			data, err := os.ReadFile(tf.path)
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", tf.path, err)
			}

			delim, fmData, _ := SplitFrontmatter(data)
			fmt.Printf("🔍 TEST %s delim=%q raw_fm=%q\n", tf.path, delim, fmData)

			frontmatter, err := UnmarshalFrontmatter(delim, fmData)
			if err != nil {
				t.Fatalf("Failed to unmarshal frontmatter: %v", err)
			}

			if !CheckCondition(frontmatter, "tags contains 'test'") {
				t.Errorf("CheckCondition(tags contains 'test') = false; want true for %s", tf.path)
			}
		})
	}
}

// TestTomlArrayFormat tests that TOML arrays are correctly formatted inline.
func TestTomlArrayFormat(t *testing.T) {
	front := map[string]interface{}{
		"tags": []string{"one", "two", "three"},
	}

	marshaled, err := MarshalFrontmatter("+++", front)
	if err != nil {
		t.Fatalf("Failed to marshal TOML: %v", err)
	}

	content := string(marshaled)
	if !strings.Contains(content, "tags = [") {
		t.Errorf("Expected inline array format, got: %s", content)
	}

	// Check that the array is on a single line (no newline between brackets)
	if strings.Contains(content, "[\n") {
		t.Errorf("TOML array should be on a single line, got multi-line format: %s", content)
	}
}
