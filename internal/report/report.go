// Package report provides functionality for generating reports on the frontmatter processing.
package report

import "fmt"

// Stats holds the statistics for the frontmatter processing.
var Stats = struct {
	Processed int
	Matched   int
	Updated   int
	LintFails int
	LintFixed int
}{}

// ModifiedFiles is a slice of strings containing the paths of the files that were modified.
var ModifiedFiles []string

// Print prints the report to the console.
func Print() {
	fmt.Printf("\n📊 Report:\n")
	fmt.Printf("Processed: %d files\n", Stats.Processed)
	fmt.Printf("Matched condition: %d\n", Stats.Matched)
	fmt.Printf("Updated frontmatter: %d\n", Stats.Updated)
	if Stats.LintFails > 0 || Stats.LintFixed > 0 {
		fmt.Printf("Lint violations: %d\n", Stats.LintFails)
		fmt.Printf("Fields auto-fixed: %d\n", Stats.LintFixed)
	}
	fmt.Printf("Skipped: %d\n", Stats.Processed-Stats.Matched)

	if len(ModifiedFiles) > 0 {
		fmt.Printf("\nModified Files:\n")
		for _, file := range ModifiedFiles {
			fmt.Printf("- %s\n", file)
		}
	}
}
