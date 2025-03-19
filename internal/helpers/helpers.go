// Package helpers provides helper functions for the hugo-frontmatter-toolbox.
package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
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
	// Make a copy of the map to avoid modifying the original
	frontCopy := make(map[string]interface{})
	for k, v := range front {
		frontCopy[k] = v
	}

	switch delimiter {
	case YamlDelimiter:
		// For YAML, we'll create a custom ordered output to ensure consistent field order
		var orderedYAML strings.Builder

		// Define the order of frontmatter fields
		orderedFields := []string{"title", "date", "draft", "series", "categories", "tags"}

		// First, add the ordered fields if they exist
		for _, field := range orderedFields {
			if value, exists := frontCopy[field]; exists {
				addYAMLField(&orderedYAML, field, value)
				delete(frontCopy, field) // Remove so we don't add it again
			}
		}

		// Now add any remaining fields in alphabetical order
		var remainingFields []string
		for field := range frontCopy {
			remainingFields = append(remainingFields, field)
		}
		sort.Strings(remainingFields)

		for _, field := range remainingFields {
			addYAMLField(&orderedYAML, field, frontCopy[field])
		}

		return []byte(orderedYAML.String()), nil

	case TomlDelimiter:
		// For TOML, we'll build output with arrays explicitly formatted inline
		var buf bytes.Buffer

		// Sort keys for consistent output
		var keys []string
		for k := range frontCopy {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			value := frontCopy[key]

			switch v := value.(type) {
			case []string:
				// Format string arrays inline
				var parts []string
				for _, s := range v {
					parts = append(parts, fmt.Sprintf("%q", s))
				}
				fmt.Fprintf(&buf, "%s = [%s]\n", key, strings.Join(parts, ", "))

			case []interface{}:
				// Format interface arrays inline
				var parts []string
				for _, item := range v {
					switch val := item.(type) {
					case string:
						parts = append(parts, fmt.Sprintf("%q", val))
					default:
						parts = append(parts, fmt.Sprintf("%v", val))
					}
				}
				fmt.Fprintf(&buf, "%s = [%s]\n", key, strings.Join(parts, ", "))

			case string:
				fmt.Fprintf(&buf, "%s = %q\n", key, v)

			case bool:
				fmt.Fprintf(&buf, "%s = %t\n", key, v)

			case int, int64, float64:
				fmt.Fprintf(&buf, "%s = %v\n", key, v)

			case time.Time:
				fmt.Fprintf(&buf, "%s = %q\n", key, v.Format(time.RFC3339))

			default:
				// For complex types, use standard TOML marshaling
				singleField := map[string]interface{}{key: value}
				fieldData, err := toml.Marshal(singleField)
				if err == nil {
					// Clean up any multiline arrays in the output
					fieldStr := string(fieldData)
					fieldStr = strings.ReplaceAll(fieldStr, "[\n", "[")
					fieldStr = strings.ReplaceAll(fieldStr, "\n]", "]")

					// Clean up whitespace between array elements
					fieldStr = strings.ReplaceAll(fieldStr, ",\n", ", ")

					buf.WriteString(fieldStr)
				} else {
					// Fallback for unsupported types
					fmt.Fprintf(&buf, "%s = %v\n", key, value)
				}
			}
		}

		return buf.Bytes(), nil

	case JsonDelimiter:
		return json.MarshalIndent(frontCopy, "", "  ")

	default:
		return nil, fmt.Errorf("unsupported frontmatter format: %s", delimiter)
	}
}

// addYAMLField adds a field to the YAML builder with proper formatting
func addYAMLField(builder *strings.Builder, field string, value interface{}) {
	// Handle arrays specially for inline format
	switch v := value.(type) {
	case []interface{}, []string:
		items := formatArrayItems(v)
		builder.WriteString(fmt.Sprintf("%s: [%s]\n", field, strings.Join(items, ", ")))
	default:
		// For non-array values, use standard formatting
		builder.WriteString(fmt.Sprintf("%s: %v\n", field, formatYAMLValue(value)))
	}
}

// formatArrayItems converts array items to properly formatted strings
func formatArrayItems(arr interface{}) []string {
	var items []string

	switch v := arr.(type) {
	case []interface{}:
		for _, item := range v {
			items = append(items, formatYAMLValue(item))
		}
	case []string:
		for _, item := range v {
			// Quote string items
			items = append(items, fmt.Sprintf("%q", item))
		}
	}

	return items
}

// formatYAMLValue formats a value for YAML output
func formatYAMLValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		// Check if string needs quotes
		if needsQuotes(v) {
			return fmt.Sprintf("%q", v)
		}
		return v
	case bool, int, float64:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// needsQuotes checks if a string needs quotes in YAML
func needsQuotes(s string) bool {
	// Add quotes if string contains special characters or could be interpreted as another type
	if s == "" {
		return true
	}

	// Check for common values that need quotes
	if s == "true" || s == "false" || s == "yes" || s == "no" || s == "null" || s == "~" {
		return true
	}

	// Check for numbers
	if _, err := fmt.Sscanf(s, "%f", new(float64)); err == nil {
		return true
	}

	// Check for special characters
	for _, ch := range s {
		if strings.ContainsRune("{}[]#&*!|>'\"%@`, ", ch) {
			return true
		}
	}

	// Check if string starts with special characters
	if strings.HasPrefix(s, "-") || strings.HasPrefix(s, ":") || strings.HasPrefix(s, "?") {
		return true
	}

	return false
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
