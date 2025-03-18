// Package internal implements the core logic of the hugo-frontmatter-toolbox.
package internal

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Daviey/hugo-frontmatter-toolbox/internal/git"
	"github.com/Daviey/hugo-frontmatter-toolbox/internal/helpers"
	"github.com/Daviey/hugo-frontmatter-toolbox/internal/report"
	"github.com/Daviey/hugo-frontmatter-toolbox/pkg/config"
)

var extractedData []map[string]string

func RunTool(cfg config.Config) error {
	info, err := os.Stat(cfg.ContentDir)
	if os.IsNotExist(err) {
		fmt.Printf("⚠️  Directory '%s' does not exist. Nothing to process.\n", cfg.ContentDir)
		return nil
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("'%s' is not a directory", cfg.ContentDir)
	}

	err = filepath.Walk(cfg.ContentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && helpers.IsMarkdownFile(path) {
			report.Stats.Processed++
			return processFile(cfg, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if cfg.ExtractKey != "" {
		return outputExtract(cfg)
	}

	if cfg.Report {
		report.Print()
	}

	if cfg.GitCommit && !cfg.DryRun && len(report.ModifiedFiles) > 0 {
		return git.CommitChanges(cfg)
	}

	return nil
}

func processFile(cfg config.Config, path string) error {
	// #nosec G304 - Path is from filepath.Walk and has been validated as markdown file
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	delimiter, fmData, body := helpers.SplitFrontmatter(data)
	if delimiter == "" {
		return nil
	}

	front, err := helpers.UnmarshalFrontmatter(delimiter, fmData)
	if err != nil {
		return err
	}

	if cfg.ExtractKey != "" {
		val := "<missing>"
		if v, ok := front[cfg.ExtractKey]; ok {
			val = fmt.Sprintf("%v", v)
		}
		extractedData = append(extractedData, map[string]string{
			"file":  path,
			"key":   cfg.ExtractKey,
			"value": val,
		})
		return nil
	}

	if cfg.Condition != "" && !helpers.EvaluateConditions(front, cfg.Condition) {
		return nil
	}
	report.Stats.Matched++

	if cfg.Lint {
		lintAndFix(cfg, front)
	}

	if cfg.SetField != "" {
		k, v := helpers.ParseSet(cfg.SetField)
		front[k] = v
		report.Stats.Updated++
	}

	updatedFront, err := helpers.MarshalFrontmatter(delimiter, front)
	if err != nil {
		return err
	}

	hasChanges := string(fmData) != string(updatedFront)
	if hasChanges {
		if cfg.DryRun || (!cfg.Yes && !cfg.DryRun) {
			if err := helpers.ShowFrontmatterDiff(path, fmData, updatedFront, delimiter, cfg.DiffContext); err != nil {
				return err
			}
		}

		if !cfg.Yes && !cfg.DryRun {
			ok, err := confirm(path)
			if err != nil {
				return err
			}
			if !ok {
				fmt.Printf("Skipping %s\n", path)
				return nil
			}
		}
	} else {
		return nil
	}

	if cfg.DryRun {
		return nil
	}

	var buf bytes.Buffer
	if delimiter == "{" {
		buf.Write(updatedFront)
	} else {
		buf.WriteString(delimiter + "\n")
		buf.Write(updatedFront)
		buf.WriteString(delimiter + "\n")
	}
	buf.Write(body)

	if err := os.WriteFile(path, buf.Bytes(), 0600); err != nil {
		return err
	}
	report.ModifiedFiles = append(report.ModifiedFiles, path)
	return nil
}

func confirm(path string) (bool, error) {
	fmt.Printf("Apply changes to %s? (y/N): ", path)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false, err
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes", nil
}

func lintAndFix(cfg config.Config, front map[string]interface{}) {
	hasIssue := false
	for _, req := range cfg.RequiredFields {
		if _, ok := front[req]; !ok {
			hasIssue = true
			if cfg.Fix {
				front[req] = ""
				report.Stats.LintFixed++
			}
		}
	}
	for _, block := range cfg.ProhibitedFields {
		if _, ok := front[block]; ok {
			hasIssue = true
			if cfg.Fix {
				delete(front, block)
				report.Stats.LintFixed++
			}
		}
	}
	if hasIssue {
		report.Stats.LintFails++
	}
}

func outputExtract(cfg config.Config) error {
	switch cfg.ExtractFormat {
	case "json":
		out, _ := json.MarshalIndent(extractedData, "", "  ")
		fmt.Println(string(out))
	case "csv":
		writer := csv.NewWriter(os.Stdout)
		_ = writer.Write([]string{"file", "key", "value"})
		for _, row := range extractedData {
			_ = writer.Write([]string{row["file"], row["key"], row["value"]})
		}
		writer.Flush()
	default:
		for _, row := range extractedData {
			fmt.Printf("%s: %s = %s\n", row["file"], row["key"], row["value"])
		}
	}
	return nil
}
