// Package helpers provides helper functions for the hugo-frontmatter-toolbox.
package helpers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v2"
)

// IsMarkdownFile checks if a file path has a .md extension.
func IsMarkdownFile(path string) bool {
	return strings.HasSuffix(path, ".md")
}

// SplitFrontmatter attempts to split a byte slice into frontmatter delimiter, frontmatter content, and body content.
func SplitFrontmatter(data []byte) (string, []byte, []byte) {
	content := string(data)
	if strings.HasPrefix(content, "---\n") {
		parts := strings.SplitN(content[4:], "---", 2)
		return "---", []byte(parts[0]), []byte(parts[1])
	}
	if strings.HasPrefix(content, "+++\n") {
		parts := strings.SplitN(content[4:], "+++", 2)
		return "+++", []byte(parts[0]), []byte(parts[1])
	}
	if strings.HasPrefix(content, "{") {
		idx := strings.Index(content, "}\n")
		if idx > 0 {
			return "{", []byte(content[:idx+1]), []byte(content[idx+2:])
		}
	}
	return "", nil, data
}

// UnmarshalFrontmatter unmarshals frontmatter data based on the specified delimiter (---, +++, or {).
func UnmarshalFrontmatter(delimiter string, data []byte) (map[string]interface{}, error) {
	front := make(map[string]interface{})
	switch delimiter {
	case YamlDelimiter:
		err := yaml.Unmarshal(data, &front)
		return front, err
	case TomlDelimiter:
		err := toml.Unmarshal(data, &front)
		return front, err
	case JsonDelimiter:
		err := json.Unmarshal(data, &front)
		return front, err
	}
	return front, fmt.Errorf("unknown frontmatter format: %s", delimiter)
}

// MarshalFrontmatter marshals frontmatter data based on the specified delimiter (---, +++, or {).
func MarshalFrontmatter(delimiter string, front map[string]interface{}) ([]byte, error) {
	switch delimiter {
	case YamlDelimiter:
		return yaml.Marshal(front)
	case TomlDelimiter:
		data, err := toml.Marshal(front)
		if err != nil {
			return nil, err
		}
		// Post-process TOML arrays to be inline
		re := regexp.MustCompile(`(?m)^(\s*\w+\s*=\s*)\[\s*\n(\s*)(.+?)\s*\n(\s*)]`)
		processed := re.ReplaceAllString(string(data), "$1[$3]")
		return []byte(processed), nil
	case JsonDelimiter:
		return json.MarshalIndent(front, "", "  ")
	default:
		return nil, fmt.Errorf("unsupported frontmatter format: %s", delimiter)
	}
}

// ParseSet parses a string in the format "key=value" and returns the key and value as a string and interface{}.
func ParseSet(input string) (string, interface{}) {
	parts := strings.SplitN(input, "=", 2)
	key := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])
	if val == "true" {
		return key, true
	}
	if val == "false" {
		return key, false
	}
	return key, val
}

// EvaluateConditions evaluates a complex condition string against the frontmatter data.
func EvaluateConditions(front map[string]interface{}, cond string) bool {
	// Replace AND/OR with their symbolic equivalents
	cond = strings.ReplaceAll(cond, " AND ", " && ")
	cond = strings.ReplaceAll(cond, " OR ", " || ")
	// Split the condition into clauses separated by OR
	clauses := strings.Split(cond, "||")
	// Iterate through each OR clause
	for _, clause := range clauses {
		all := true
		// Split the clause into AND conditions
		ands := strings.Split(clause, "&&")
		// Iterate through each AND condition
		for _, andClause := range ands {
			andClause = strings.TrimSpace(andClause)
			// Check if the individual condition is met
			if !CheckCondition(front, andClause) {
				all = false
				break
			}
		}
		// If all AND conditions are true, the OR clause is true
		if all {
			return true
		}
	}
	// If none of the OR clauses are true, the entire condition is false
	return false
}

// CheckCondition checks if a single condition is met within the frontmatter data.
func CheckCondition(front map[string]interface{}, cond string) bool {
	// Date comparison
	if strings.Contains(cond, "<") && strings.HasPrefix(cond, "date") {
		cutoff := strings.TrimSpace(strings.Split(cond, "<")[1])
		cutoffTime, _ := time.Parse("2006-01-02", cutoff)
		if val, ok := front["date"]; ok {
			dateStr := fmt.Sprintf("%v", val)
			postTime, err := time.Parse("2006-01-02", dateStr)
			if err == nil && postTime.Before(cutoffTime) {
				return true
			}
		}
		return false
	}

	// Contains check for arrays
	if strings.Contains(cond, "contains") {
		parts := strings.Split(cond, "contains")
		key := strings.TrimSpace(parts[0])
		rawNeedle := strings.TrimSpace(parts[1])
		rawNeedle = strings.Trim(rawNeedle, "\"'")

		if val, ok := front[key]; ok {
			flat := flattenToStrings(val)
			for _, item := range flat {
				if item == rawNeedle {
					return true
				}
			}
		}
		return false
	}

	// Equality check
	if strings.Contains(cond, "=") {
		parts := strings.Split(cond, "=")
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if val, ok := front[key]; ok {
			return fmt.Sprintf("%v", val) == value
		}
		return false
	}
	return false
}

// flattenToStrings converts a slice of interfaces to a slice of strings.
func flattenToStrings(v interface{}) []string {
	if v == nil {
		return []string{}
	}
	out := []string{}

	// Fix: S1011 - Use direct append for string slices
	switch arr := v.(type) {
	case []interface{}:
		for _, item := range arr {
			out = append(out, fmt.Sprintf("%v", item))
		}
	case []string:
		// Direct append instead of loop
		out = append(out, arr...)
	default:
		out = append(out, fmt.Sprintf("%v", v))
	}
	return out
}

const (
	YamlDelimiter = "---"
	TomlDelimiter = "+++"
	JsonDelimiter = "{"
)

// ShowFrontmatterDiff calculates and prints the diff between the original and updated frontmatter.
func ShowFrontmatterDiff(file string, original, updated []byte, delim string, context int) error {
	fmt.Printf("\n🔍 Frontmatter Diff for: %s\n", file)

	// Instead of line-by-line comparison, we'll unmarshal both frontmatters
	// to compare the actual data structures, then show the original lines
	// that have changed
	origFront, err := UnmarshalFrontmatter(delim, original)
	if err != nil {
		return err
	}

	updatedFront, err := UnmarshalFrontmatter(delim, updated)
	if err != nil {
		return err
	}

	// Get all lines from original frontmatter
	origLines := strings.Split(string(original), "\n")

	// Track which keys have changed
	changedKeys := make(map[string]bool)

	// Find changed keys
	for k, v := range updatedFront {
		origVal, exists := origFront[k]
		if !exists || fmt.Sprintf("%v", origVal) != fmt.Sprintf("%v", v) {
			changedKeys[k] = true
		}
	}

	// Also check for keys that were in original but removed in updated
	for k := range origFront {
		if _, exists := updatedFront[k]; !exists {
			changedKeys[k] = true
		}
	}

	// Print unchanged lines and identify changed sections
	inChangedSection := false

	for _, line := range origLines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}

		// Check if this line starts a new key
		if !strings.HasPrefix(trimmedLine, "-") && !strings.HasPrefix(trimmedLine, " ") && strings.Contains(trimmedLine, ":") {
			parts := strings.SplitN(trimmedLine, ":", 2)
			key := strings.TrimSpace(parts[0])
			inChangedSection = changedKeys[key]
		}

		// If we're in a changed section, collect the lines to show as removed
		if inChangedSection {
			_, _ = color.New(color.FgRed).Printf("- %s\n", line)
		} else {
			fmt.Printf("  %s\n", line)
		}
	}

	// Now show the new values for changed keys
	for k := range changedKeys {
		if newVal, exists := updatedFront[k]; exists {
			// Format the value based on its type
			switch newVal.(type) {
			case []interface{}, []string:
				// For arrays, show the new format exactly as it would be in the updated frontmatter
				for i, line := range strings.Split(string(updated), "\n") {
					trimmedLine := strings.TrimSpace(line)
					if strings.HasPrefix(trimmedLine, k+":") ||
						(i > 0 && strings.HasPrefix(trimmedLine, "-") && strings.Contains(string(updated), k+":")) {
						_, _ = color.New(color.FgGreen).Printf("+ %s\n", line)
					}
				}
			default:
				// For simple values, just show the key-value pair
				_, _ = color.New(color.FgGreen).Printf("+ %s: %v\n", k, newVal)
			}
		}
	}

	fmt.Println()
	return nil
}
