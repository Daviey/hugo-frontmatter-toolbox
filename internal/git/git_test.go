// Package git_test contains unit tests for the git package.
package git

import (
	"os"
	"os/exec"
	"testing"

	"github.com/Daviey/hugo-frontmatter-toolbox/pkg/config"
)

// TestGenerateCommitMessage tests the generateCommitMessage function with default configuration.
func TestGenerateCommitMessage(t *testing.T) {
	cfg := config.Config{
		GcMsg: "",
	}
	msg := generateCommitMessage(cfg)
	if msg != "chore: batch update via hugo-frontmatter-toolbox" {
		t.Errorf("Expected default commit message, got: %s", msg)
	}
}

// TestGenerateCommitMessage_Custom tests the generateCommitMessage function with a custom commit message.
func TestGenerateCommitMessage_Custom(t *testing.T) {
	cfg := config.Config{
		GcMsg: "custom msg",
	}
	msg := generateCommitMessage(cfg)
	if msg != "custom msg" {
		t.Errorf("Expected custom commit message, got: %s", msg)
	}
}

func TestCommitChanges_Success(t *testing.T) {
	origExec := execCommand
	defer func() { execCommand = origExec }()

	// Safe mock for testing
	execCommand = func(name string, arg ...string) *exec.Cmd {
		testBin := os.Args[0]
		// #nosec G204 -- safe in controlled test environment
		cmd := exec.Command(testBin, "-test.run=TestHelperProcess", "--", name)
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
		return cmd
	}

	// Create a temp directory with a dummy .git folder
	tempDir := t.TempDir()
	if err := os.Mkdir(tempDir+"/.git", 0700); err != nil {
		t.Fatalf("failed to create dummy .git dir: %v", err)
	}

	// Change to tempDir and restore cwd after test
	oldCwd, _ := os.Getwd()
	defer func() {
		if err := os.Chdir(oldCwd); err != nil {
			t.Fatalf("failed to restore working dir: %v", err)
		}
	}()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}

	cfg := config.Config{
		GcMsg: "test commit",
	}

	if err := CommitChanges(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestHelperProcess is a helper process that simulates git commands for testing purposes.
// It checks for the GO_WANT_HELPER_PROCESS environment variable and exits if it's not set to "1".
// This function is not meant to be called directly.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}
