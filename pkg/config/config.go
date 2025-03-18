// Package config defines the configuration options for the hugo-frontmatter-toolbox.
package config

type Config struct {
	ContentDir       string
	SetField         string
	Condition        string
	DryRun           bool
	Report           bool
	DiffContext      int
	Lint             bool
	Fix              bool
	RequiredFields   []string
	ProhibitedFields []string
	GitCommit        bool
	GcMsg            string
	Yes              bool
	ExtractKey       string
	ExtractFormat    string
}
