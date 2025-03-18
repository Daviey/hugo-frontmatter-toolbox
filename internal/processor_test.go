// Package internal_test contains unit tests for the internal package.
package internal

import (
	"os"
	"testing"

	"github.com/Daviey/hugo-frontmatter-toolbox/pkg/config"
)

// TestRunTool_InvalidDir tests RunTool with an invalid content directory.
func TestRunTool_InvalidDir(t *testing.T) {
	cfg := config.Config{ContentDir: "invalid-dir"}
	err := RunTool(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRunTool_NotADir tests RunTool when the content directory is a file, not a directory.
func TestRunTool_NotADir(t *testing.T) {
	_ = os.WriteFile("file.md", []byte("test"), 0600)
	defer func() {
		if err := os.Remove("file.md"); err != nil {
			t.Logf("Failed to remove test file: %v", err)
		}
	}()

	cfg := config.Config{ContentDir: "file.md"}
	err := RunTool(cfg)
	if err == nil {
		t.Errorf("expected error for non-directory path")
	}
	if err.Error() != "'file.md' is not a directory" {
		t.Errorf("expected specific error message, got: %v", err)
	}
}

// TestRunTool_EmptyDir tests RunTool with an empty content directory.
func TestRunTool_EmptyDir(t *testing.T) {
	_ = os.Mkdir("testdir", 0700)
	defer func() {
		if err := os.Remove("testdir"); err != nil {
			t.Logf("Failed to remove test directory: %v", err)
		}
	}()

	cfg := config.Config{ContentDir: "testdir"}
	err := RunTool(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
