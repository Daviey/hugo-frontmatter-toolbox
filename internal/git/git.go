// Package git provides functionality for interacting with git repositories.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Daviey/hugo-frontmatter-toolbox/internal/report"
	"github.com/Daviey/hugo-frontmatter-toolbox/pkg/config"
)

// execCommand is a variable that holds the command execution function.
// It defaults to exec.Command but can be overridden for testing purposes.
var execCommand = exec.Command // 👈 allows test override

// CommitChanges commits the changes made to the files in the content directory.
func CommitChanges(cfg config.Config) error {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return fmt.Errorf("--gc enabled but no .git repo found")
	}
	args := append([]string{"add"}, report.ModifiedFiles...)
	if err := execCommand("git", args...).Run(); err != nil {
		return fmt.Errorf("git add failed: %v", err)
	}

	commitMsg := generateCommitMessage(cfg)

	if err := execCommand("git", "commit", "-m", commitMsg).Run(); err != nil {
		return fmt.Errorf("git commit failed: %v", err)
	}
	fmt.Printf("✅ Git commit created: %q\n", commitMsg)
	return nil
}

// generateCommitMessage generates a commit message based on the configuration.
func generateCommitMessage(cfg config.Config) string {
	if cfg.GcMsg != "" {
		return cfg.GcMsg
	}
	var parts []string
	if cfg.SetField != "" {
		parts = append(parts, fmt.Sprintf("set %s", cfg.SetField))
	}
	if cfg.Condition != "" {
		parts = append(parts, fmt.Sprintf("filtered on %q", cfg.Condition))
	}
	if cfg.Lint && cfg.Fix {
		parts = append(parts, "auto-fixed lint issues")
	} else if cfg.Lint {
		parts = append(parts, "ran lint")
	}
	if len(parts) == 0 {
		return "chore: batch update via hugo-frontmatter-toolbox"
	}
	return "chore: " + strings.Join(parts, ", ")
}
