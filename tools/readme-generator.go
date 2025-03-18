// Package main provides a README.md generator for hugo-frontmatter-toolbox
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// Template for README.md
const readmeTemplate = `<!-- THIS FILE IS AUTO-GENERATED. DO NOT EDIT DIRECTLY. -->
<!-- To update this file, run: make readme -->

# hugo-frontmatter-toolbox

{{ .CoverageBadge }}

A CLI tool for batch editing Hugo frontmatter (YAML, TOML, JSON).

## Features

- 🔄 **Batch update frontmatter fields** - Easily modify fields like ` + "`draft: true`" + ` across multiple files
- 🧩 **Conditional filtering** - Target specific content with ` + "`--if`" + ` conditions
- 🧹 **Frontmatter linting** - Check for required or prohibited fields
- 🔧 **Automatic fixes** - Auto-fix lint issues with ` + "`--fix`" + `
- 🔍 **Diff visualization** - Preview changes with colorized diffs using ` + "`--dry-run`" + `
- 📊 **Summary reporting** - Get concise execution summaries with ` + "`--report`" + `
- 🔀 **Git integration** - Automatically commit changes with ` + "`--gc`" + `
- ✓ **Non-interactive mode** - Skip confirmation prompts with ` + "`--yes`" + ` or ` + "`-y`" + `

## Installation

### From Source

` + "```bash" + `
git clone https://github.com/Daviey/hugo-frontmatter-toolbox.git
cd hugo-frontmatter-toolbox
make install
` + "```" + `

### Using Go Install

` + "```bash" + `
go install github.com/Daviey/hugo-frontmatter-toolbox@latest
` + "```" + `

## Usage Examples

{{ .UsageExamples }}

## Understanding Conditions

You can use the ` + "`--if`" + ` flag to filter which markdown files to modify. The tool supports:

- **Simple comparison**: ` + "`--if \"draft=true\"`" + `
- **Date comparison**: ` + "`--if \"date<2022-01-01\"`" + `
- **List field checks**: ` + "`--if \"tags contains 'draft'\"`" + `
- **Boolean operators**: Use ` + "`AND`" + ` and ` + "`OR`" + ` to combine conditions

Examples with multiple conditions:

1. Find all posts from 2023 that are not drafts and tag them as featured:
` + "```bash" + `
hugo-frontmatter-toolbox --set featured=true --if "date>2023-01-01 AND tags contains 'important' AND draft=false"
` + "```" + `

2. Move all old posts (before 2020) from the 'news' category to draft status:
` + "```bash" + `
hugo-frontmatter-toolbox --set draft=true --if "date<2020-01-01 AND categories = 'news'"
` + "```" + `

3. Find posts that have neither a description nor a summary field:
` + "```bash" + `
hugo-frontmatter-toolbox --set description="Auto-generated description" --if "description=nil AND summary=nil"
` + "```" + `

## Advanced Usage

### Using with Different Frontmatter Formats

Hugo supports YAML, TOML, and JSON for frontmatter. This tool automatically detects and preserves the format:

- **YAML** (delimited by ` + "`---`" + `):
` + "```yaml" + `
---
title: "My Post"
draft: true
---
` + "```" + `

- **TOML** (delimited by ` + "`+++`" + `):
` + "```toml" + `
+++
title = "My Post"
draft = true
+++
` + "```" + `

- **JSON** (wrapped in curly braces):
` + "```json" + `
{
  "title": "My Post",
  "draft": true
}
` + "```" + `

### Bulk Migration Scenarios

#### Migrating from WordPress/Ghost/Jekyll

If you're migrating from another platform and need to add Hugo-specific frontmatter:

` + "```bash" + `
# Add layout and adjust categories
hugo-frontmatter-toolbox --set layout=post --set "hugo_categories=oldcategories" --yes
` + "```" + `

#### Handling Taxonomy Changes

When you need to rename or restructure taxonomies:

` + "```bash" + `
# Convert 'topics' to 'categories'
hugo-frontmatter-toolbox --if "topics contains 'technology'" --set categories=technology --yes
` + "```" + `

#### Batch Processing with Git Integration

Automatically create Git commits when making systematic updates:

` + "```bash" + `
# Update post format across all posts with Git commit
hugo-frontmatter-toolbox --set format=hugo --gc --gc-msg "chore: standardize post format to hugo" --yes
` + "```" + `

## Recipes

Here are some common use cases and recipes for solving specific problems with hugo-frontmatter-toolbox:

### Working with Dates

**Set publication date for drafts:**
` + "```bash" + `
hugo-frontmatter-toolbox --if "draft=true" --set "date=$(date +%Y-%m-%d)" --yes
` + "```" + `

**Update lastmod field for recently modified files:**
` + "```bash" + `
find content -name "*.md" -mtime -7 | xargs hugo-frontmatter-toolbox --set "lastmod=$(date +%Y-%m-%d)" --yes
` + "```" + `

### SEO Optimization

**Add description to posts missing it:**
` + "```bash" + `
hugo-frontmatter-toolbox --lint --required "description" --fix --yes
` + "```" + `

**Set canonical URL for all posts:**
` + "```bash" + `
hugo-frontmatter-toolbox --set "canonical_url=https://mysite.com/path" --yes
` + "```" + `

### Content Reorganization

**Move content to a different section:**
` + "```bash" + `
hugo-frontmatter-toolbox --if "section=blog" --set section=articles --yes
` + "```" + `

**Mark posts with specific tags as featured:**
` + "```bash" + `
hugo-frontmatter-toolbox --if "tags contains 'highlight'" --set featured=true --yes
` + "```" + `

### Extract values draft from frontmatter
` + "```bash" + `
hugo-frontmatter-toolbox --extract draft
` + "```" + `

## Flags Reference

{{ .FlagsTable }}

## Development

### Build and Test

Build the application:
` + "```bash" + `
make build
` + "```" + `

Run tests:
` + "```bash" + `
make test
` + "```" + `

Generate a test coverage report:
` + "```bash" + `
make cover
# Opens a browser with detailed coverage information
` + "```" + `

Current test coverage is reflected in the badge at the top of this README.

Run linters and format code:
` + "```bash" + `
make fmt lint
` + "```" + `

### Contributing

1. Fork the repository
2. Create a feature branch: ` + "`git checkout -b feature-name`" + `
3. Commit your changes: ` + "`git commit -am 'Add feature'`" + `
4. Push to the branch: ` + "`git push origin feature-name`" + `
5. Submit a pull request

## License
MIT
`

// Command example with description template
const exampleTemplate = `### {{.Title}}
{{.Description}}

` + "```bash" + `
hugo-frontmatter-toolbox {{.Command}}
` + "```" + `
`

// Flag structure for template
type Flag struct {
	Name        string
	Description string
}

// Example structure for usage examples
type Example struct {
	Title       string
	Description string
	Command     string
}

// Helper function to extract command help
func getCommandHelp() (string, error) {
	// Check if binary exists in current dir
	_, err := os.Stat("./hugo-frontmatter-toolbox")
	if err == nil {
		cmd := exec.Command("./hugo-frontmatter-toolbox", "--help")
		output, err := cmd.CombinedOutput()
		if err == nil {
			return string(output), nil
		}
	}

	// Try building it first if needed
	buildCmd := exec.Command("go", "build", "-o", "hugo-frontmatter-toolbox")
	err = buildCmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to build binary: %v", err)
	}

	cmd := exec.Command("./hugo-frontmatter-toolbox", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get help text: %v", err)
	}

	return string(output), nil
}

// Parse flags from help output
func parseFlags(helpText string) ([]Flag, error) {
	var flags []Flag

	lines := strings.Split(helpText, "\n")
	capturingFlags := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Look for the flags section
		if strings.Contains(trimmed, "Flags:") {
			capturingFlags = true
			continue
		}

		if capturingFlags && strings.HasPrefix(trimmed, "--") {
			parts := strings.SplitN(trimmed, "   ", 2)
			if len(parts) >= 2 {
				name := strings.TrimSpace(parts[0])
				// Extract just the flag name without shortcuts
				flagName := strings.TrimSpace(strings.Split(name, ",")[0])
				description := strings.TrimSpace(parts[1])

				flags = append(flags, Flag{
					Name:        flagName,
					Description: description,
				})
			}
		}
	}

	return flags, nil
}

// Generate example usage
func generateExamples() string {
	examples := []Example{
		{
			Title:       "Basic frontmatter update",
			Description: "Set the `draft` field to `true` for all markdown files in the content directory:",
			Command:     "--set draft=true",
		},
		{
			Title:       "Conditional update",
			Description: "Set `draft=true` only for posts with a date before 2022 that are currently not drafts:",
			Command:     "--set draft=true --if \"date<2022-01-01 AND draft=false\"",
		},
		{
			Title:       "Conditional by tags or categories",
			Description: "Mark posts as draft if they have the 'beta' tag or belong to the 'drafts' category:",
			Command:     "--set draft=true --if \"tags contains 'beta' OR categories = 'drafts'\"",
		},
		{
			Title:       "Lint frontmatter fields",
			Description: "Check if all posts have the required 'title' and 'date' fields, and ensure no post has the deprecated 'obsolete_field':",
			Command:     "--lint --required \"title,date\" --prohibited \"obsolete_field\"",
		},
		{
			Title:       "Lint and autofix",
			Description: "Check for required/prohibited fields and automatically fix issues by adding missing fields and removing prohibited ones:",
			Command:     "--lint --fix --required \"title,date\" --prohibited \"obsolete_field\"",
		},
		{
			Title:       "Dry-run diff mode",
			Description: "Preview changes without modifying files, showing a colorized diff of what would change:",
			Command:     "--set draft=true --dry-run",
		},
		{
			Title:       "Diff context control",
			Description: "Adjust the amount of context shown in diff output to 5 lines (default is 2):",
			Command:     "--set draft=true --dry-run --diff-context 5",
		},
		{
			Title:       "Git auto-commit",
			Description: "Automatically commit changes to git after updating frontmatter:",
			Command:     "--set draft=true --gc",
		},
		{
			Title:       "Git auto-commit with custom message",
			Description: "Automatically commit changes with a custom commit message:",
			Command:     "--set draft=true --gc --gc-msg \"chore: mark old posts as draft\"",
		},
		{
			Title:       "Reporting summary",
			Description: "Generate a summary report after execution showing stats about processed files:",
			Command:     "--set draft=true --report",
		},
		{
			Title:       "Non-interactive mode",
			Description: "Skip all confirmation prompts and apply changes automatically:",
			Command:     "--set draft=true --yes",
		},
		{
			Title:       "Custom content directory",
			Description: "Process markdown files in a directory other than the default 'content':",
			Command:     "--content-dir=\"my-custom-content\" --set draft=true",
		},
		{
			Title:       "Extract draft",
			Description: "Extract all values of the 'draft' field:",
			Command:     "--extract draft",
		},
	}

	var result strings.Builder
	tmpl, err := template.New("example").Parse(exampleTemplate)
	if err != nil {
		return "Error generating examples: " + err.Error()
	}

	for _, example := range examples {
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, example)
		if err != nil {
			continue
		}
		result.WriteString(buf.String() + "\n")
	}

	return result.String()
}

// Generate flags table in markdown format
func generateFlagsTable(flags []Flag) string {
	var result strings.Builder

	result.WriteString("| Flag | Description |\n")
	result.WriteString("|------|-------------|\n")

	for _, flag := range flags {
		result.WriteString(fmt.Sprintf("| `%s` | %s |\n", flag.Name, flag.Description))
	}

	return result.String()
}

// Collect test coverage information
func getTestCoverage() (float64, error) {
	// Run tests with coverage
	cmd := exec.Command("go", "test", "./...", "-cover")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("failed to run tests with coverage: %v", err)
	}

	// Parse the output to find coverage percentage
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	var totalCoverage float64
	var packages int

	for _, line := range lines {
		if strings.Contains(line, "coverage:") {
			parts := strings.Split(line, "coverage: ")
			if len(parts) >= 2 {
				percentStr := strings.Split(parts[1], "%")[0]
				percent, err := strconv.ParseFloat(percentStr, 64)
				if err == nil {
					totalCoverage += percent
					packages++
				}
			}
		}
	}

	if packages == 0 {
		return 0, fmt.Errorf("no package coverage information found")
	}

	// Return average coverage across packages
	return totalCoverage / float64(packages), nil
}

func main() {
	helpText, err := getCommandHelp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting command help: %v\n", err)
		os.Exit(1)
	}

	flags, err := parseFlags(helpText)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Get test coverage
	coverage, err := getTestCoverage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Couldn't get test coverage: %v\n", err)
		// Continue without coverage info
	}

	coverageBadge := ""
	if err == nil {
		// Generate coverage badge
		color := "red"
		if coverage >= 80 {
			color = "brightgreen"
		} else if coverage >= 70 {
			color = "green"
		} else if coverage >= 60 {
			color = "yellowgreen"
		} else if coverage >= 50 {
			color = "yellow"
		} else if coverage >= 40 {
			color = "orange"
		}

		coverageBadge = fmt.Sprintf("![Test Coverage](https://img.shields.io/badge/coverage-%.1f%%25-%s)", coverage, color)
	}

	tmpl, err := template.New("readme").Parse(readmeTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing template: %v\n", err)
		os.Exit(1)
	}

	data := struct {
		UsageExamples  string
		FlagsTable     string
		GenerationTime string
		CoverageBadge  string
	}{
		UsageExamples:  generateExamples(),
		FlagsTable:     generateFlagsTable(flags),
		GenerationTime: time.Now().Format("2006-01-02 15:04:05"),
		CoverageBadge:  coverageBadge,
	}

	var output bytes.Buffer
	err = tmpl.Execute(&output, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template: %v\n", err)
		os.Exit(1)
	}

	// Fix: G306 - Changed file permission from 0644 to 0600
	err = os.WriteFile("README.md", output.Bytes(), 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing README.md: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ README.md successfully generated")
}
