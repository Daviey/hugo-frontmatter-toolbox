// Package cmd implements the command-line interface for the hugo-frontmatter-toolbox.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Daviey/hugo-frontmatter-toolbox/internal"
	"github.com/Daviey/hugo-frontmatter-toolbox/pkg/config"
	"github.com/spf13/cobra"
)

var (
	contentDir    string
	setField      string
	condition     string
	dryRun        bool
	report        bool
	diffContext   int
	lint          bool
	fix           bool
	requiredStr   string
	prohibitedStr string
	gitCommit     bool
	gcMsg         string
	yes           bool
	extractKey    string
	extractFormat string
	version       = "v1.0.0"
	exitFunc      = os.Exit // 👈 overrideable for tests
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	rootCmd := &cobra.Command{
		Use:   "hugo-frontmatter-toolbox",
		Short: "Batch edit Hugo frontmatter (YAML, TOML, JSON)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return internal.RunTool(config.Config{
				ContentDir:       contentDir,
				SetField:         setField,
				Condition:        condition,
				DryRun:           dryRun,
				Report:           report,
				DiffContext:      diffContext,
				Lint:             lint,
				Fix:              fix,
				RequiredFields:   parseCSV(requiredStr),
				ProhibitedFields: parseCSV(prohibitedStr),
				GitCommit:        gitCommit,
				GcMsg:            gcMsg,
				Yes:              yes,
				ExtractKey:       extractKey,
				ExtractFormat:    extractFormat,
			})
		},
	}

	rootCmd.PersistentFlags().StringVarP(&contentDir, "content-dir", "c", "content", "Path to Hugo content directory")
	rootCmd.PersistentFlags().StringVarP(&setField, "set", "s", "", "Set frontmatter field, e.g. draft=true")
	rootCmd.PersistentFlags().StringVarP(&condition, "if", "i", "", "Condition, e.g. date<2023-01-01 AND draft=false")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "n", false, "Show diff but don't write changes")
	rootCmd.PersistentFlags().BoolVar(&report, "report", false, "Show report summary after execution")
	rootCmd.PersistentFlags().BoolVar(&lint, "lint", false, "Lint for required/prohibited fields")
	rootCmd.PersistentFlags().BoolVar(&fix, "fix", false, "Fix linting issues (add/remove fields)")
	rootCmd.PersistentFlags().StringVar(&requiredStr, "required", "", "Comma-separated required fields")
	rootCmd.PersistentFlags().StringVar(&prohibitedStr, "prohibited", "", "Comma-separated prohibited fields")
	rootCmd.PersistentFlags().BoolVar(&gitCommit, "gc", false, "Auto git commit modified files")
	rootCmd.PersistentFlags().StringVar(&gcMsg, "gc-msg", "", "Override commit message for --gc")
	rootCmd.PersistentFlags().IntVar(&diffContext, "diff-context", 2, "Lines of unchanged context around diffs")
	rootCmd.PersistentFlags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompts and proceed with changes")
	rootCmd.PersistentFlags().StringVar(&extractKey, "extract", "", "Extract value of specified frontmatter key across all files")
	rootCmd.PersistentFlags().StringVar(&extractFormat, "extract-format", "plain", "Output format for --extract: plain, csv, or json")
	rootCmd.PersistentFlags().Bool("version", false, "Print version info")

	// PersistentPreRun is executed before any command and is used to display help or version information.
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if len(os.Args) == 1 {
			_ = cmd.Help()
			exitFunc(0)
		}
		if v, _ := cmd.Flags().GetBool("version"); v {
			fmt.Printf("hugo-frontmatter-toolbox %s\n", version)
			exitFunc(0)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		exitFunc(1)
	}
}

func parseCSV(input string) []string {
	if input == "" {
		return nil
	}
	parts := strings.Split(input, ",")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts
}
