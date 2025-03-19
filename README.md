<!-- THIS FILE IS AUTO-GENERATED. DO NOT EDIT DIRECTLY. -->
<!-- To update this file, run: make readme -->

# hugo-frontmatter-toolbox

![Test Coverage](https://img.shields.io/badge/coverage-47.5%25-orange)

A CLI tool for batch editing Hugo frontmatter (YAML, TOML, JSON).

## Features

- 🔄 **Batch update frontmatter fields** - Easily modify fields like `draft: true` across multiple files
- 🧩 **Conditional filtering** - Target specific content with `--if` conditions
- 🧹 **Frontmatter linting** - Check for required or prohibited fields
- 🔧 **Automatic fixes** - Auto-fix lint issues with `--fix`
- 🔍 **Diff visualization** - Preview changes with colorized diffs using `--dry-run`
- 📊 **Summary reporting** - Get concise execution summaries with `--report`
- 🔀 **Git integration** - Automatically commit changes with `--gc`
- ✓ **Non-interactive mode** - Skip confirmation prompts with `--yes` or `-y`

## Installation

### From Source

```bash
git clone https://github.com/Daviey/hugo-frontmatter-toolbox.git
cd hugo-frontmatter-toolbox
make install
```

### Using Go Install

```bash
go install github.com/Daviey/hugo-frontmatter-toolbox@latest
```

## Usage Examples

### Basic frontmatter update
Set the `draft` field to `true` for all markdown files in the content directory:

```bash
hugo-frontmatter-toolbox --set draft=true
```

### Conditional update
Set `draft=true` only for posts with a date before 2022 that are currently not drafts:

```bash
hugo-frontmatter-toolbox --set draft=true --if "date<2022-01-01 AND draft=false"
```

### Conditional by tags or categories
Mark posts as draft if they have the 'beta' tag or belong to the 'drafts' category:

```bash
hugo-frontmatter-toolbox --set draft=true --if "tags contains 'beta' OR categories = 'drafts'"
```

### Lint frontmatter fields
Check if all posts have the required 'title' and 'date' fields, and ensure no post has the deprecated 'obsolete_field':

```bash
hugo-frontmatter-toolbox --lint --required "title,date" --prohibited "obsolete_field"
```

### Lint and autofix
Check for required/prohibited fields and automatically fix issues by adding missing fields and removing prohibited ones:

```bash
hugo-frontmatter-toolbox --lint --fix --required "title,date" --prohibited "obsolete_field"
```

### Dry-run diff mode
Preview changes without modifying files, showing a colorized diff of what would change:

```bash
hugo-frontmatter-toolbox --set draft=true --dry-run
```

### Diff context control
Adjust the amount of context shown in diff output to 5 lines (default is 2):

```bash
hugo-frontmatter-toolbox --set draft=true --dry-run --diff-context 5
```

### Git auto-commit
Automatically commit changes to git after updating frontmatter:

```bash
hugo-frontmatter-toolbox --set draft=true --gc
```

### Git auto-commit with custom message
Automatically commit changes with a custom commit message:

```bash
hugo-frontmatter-toolbox --set draft=true --gc --gc-msg "chore: mark old posts as draft"
```

### Reporting summary
Generate a summary report after execution showing stats about processed files:

```bash
hugo-frontmatter-toolbox --set draft=true --report
```

### Non-interactive mode
Skip all confirmation prompts and apply changes automatically:

```bash
hugo-frontmatter-toolbox --set draft=true --yes
```

### Custom content directory
Process markdown files in a directory other than the default 'content':

```bash
hugo-frontmatter-toolbox --content-dir="my-custom-content" --set draft=true
```

### Extract draft
Extract all values of the 'draft' field:

```bash
hugo-frontmatter-toolbox --extract draft
```



## Understanding Conditions

You can use the `--if` flag to filter which markdown files to modify. The tool supports:

- **Simple comparison**: `--if "draft=true"`
- **Date comparison**: `--if "date<2022-01-01"`
- **List field checks**: `--if "tags contains 'draft'"`
- **Boolean operators**: Use `AND` and `OR` to combine conditions

Examples with multiple conditions:

1. Find all posts from 2023 that are not drafts and tag them as featured:
```bash
hugo-frontmatter-toolbox --set featured=true --if "date>2023-01-01 AND tags contains 'important' AND draft=false"
```

2. Move all old posts (before 2020) from the 'news' category to draft status:
```bash
hugo-frontmatter-toolbox --set draft=true --if "date<2020-01-01 AND categories = 'news'"
```

3. Find posts that have neither a description nor a summary field:
```bash
hugo-frontmatter-toolbox --set description="Auto-generated description" --if "description=nil AND summary=nil"
```

## Advanced Usage

### Using with Different Frontmatter Formats

Hugo supports YAML, TOML, and JSON for frontmatter. This tool automatically detects and preserves the format:

- **YAML** (delimited by `---`):
```yaml
---
title: "My Post"
draft: true
---
```

- **TOML** (delimited by `+++`):
```toml
+++
title = "My Post"
draft = true
+++
```

- **JSON** (wrapped in curly braces):
```json
{
  "title": "My Post",
  "draft": true
}
```

### Bulk Migration Scenarios

#### Migrating from WordPress/Ghost/Jekyll

If you're migrating from another platform and need to add Hugo-specific frontmatter:

```bash
# Add layout and adjust categories
hugo-frontmatter-toolbox --set layout=post --set "hugo_categories=oldcategories" --yes
```

#### Handling Taxonomy Changes

When you need to rename or restructure taxonomies:

```bash
# Convert 'topics' to 'categories'
hugo-frontmatter-toolbox --if "topics contains 'technology'" --set categories=technology --yes
```

#### Batch Processing with Git Integration

Automatically create Git commits when making systematic updates:

```bash
# Update post format across all posts with Git commit
hugo-frontmatter-toolbox --set format=hugo --gc --gc-msg "chore: standardize post format to hugo" --yes
```

## Recipes

Here are some common use cases and recipes for solving specific problems with hugo-frontmatter-toolbox:

### Working with Dates

**Set publication date for drafts:**
```bash
hugo-frontmatter-toolbox --if "draft=true" --set "date=$(date +%Y-%m-%d)" --yes
```

**Update lastmod field for recently modified files:**
```bash
find content -name "*.md" -mtime -7 | xargs hugo-frontmatter-toolbox --set "lastmod=$(date +%Y-%m-%d)" --yes
```

### SEO Optimization

**Add description to posts missing it:**
```bash
hugo-frontmatter-toolbox --lint --required "description" --fix --yes
```

**Set canonical URL for all posts:**
```bash
hugo-frontmatter-toolbox --set "canonical_url=https://mysite.com/path" --yes
```

### Content Reorganization

**Move content to a different section:**
```bash
hugo-frontmatter-toolbox --if "section=blog" --set section=articles --yes
```

**Mark posts with specific tags as featured:**
```bash
hugo-frontmatter-toolbox --if "tags contains 'highlight'" --set featured=true --yes
```

### Extract values draft from frontmatter
```bash
hugo-frontmatter-toolbox --extract draft
```

## Flags Reference

| Flag | Description |
|------|-------------|
| `--diff-context int` | Lines of unchanged context around diffs (default 2) |
| `--extract string` | Extract value of specified frontmatter key across all files |
| `--extract-format string` | Output format for --extract: plain, csv, or json (default "plain") |
| `--fix` | Fix linting issues (add/remove fields) |
| `--gc` | Auto git commit modified files |
| `--gc-msg string` | Override commit message for --gc |
| `--lint` | Lint for required/prohibited fields |
| `--prohibited string` | Comma-separated prohibited fields |
| `--report` | Show report summary after execution |
| `--required string` | Comma-separated required fields |
| `--version` | Print version info |


## Development

### Build and Test

Build the application:
```bash
make build
```

Run tests:
```bash
make test
```

Generate a test coverage report:
```bash
make cover
# Opens a browser with detailed coverage information
```

Current test coverage is reflected in the badge at the top of this README.

Run linters and format code:
```bash
make fmt lint
```

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Commit your changes: `git commit -am 'Add feature'`
4. Push to the branch: `git push origin feature-name`
5. Submit a pull request

## License
MIT
