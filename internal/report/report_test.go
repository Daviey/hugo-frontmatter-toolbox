// Package report_test contains unit tests for the report package.
package report

import (
	"testing"
)

// TestPrint tests the Print function.
func TestPrint(t *testing.T) {
	Stats = struct {
		Processed int
		Matched   int
		Updated   int
		LintFails int
		LintFixed int
	}{Processed: 10, Matched: 8, Updated: 5, LintFails: 2, LintFixed: 1}
	ModifiedFiles = []string{"file1.md", "file2.md"}

	Print()
}
