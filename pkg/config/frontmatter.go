package config

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Options holds agent options parsed from YAML frontmatter in agent files.
type Options struct {
	Model     string `yaml:"model"`
	AgentType string `yaml:"agent"`
}

// validModels contains accepted full model IDs for agent frontmatter.
var validModels = map[string]bool{
	"claude-opus-4.6":      true,
	"claude-opus-4.6-fast": true,
	"claude-opus-4.5":      true,
	"claude-sonnet-4.6":    true,
	"claude-sonnet-4.5":    true,
	"claude-sonnet-4":      true,
	"claude-haiku-4.5":     true,
	"gpt-5.4":              true,
	"gpt-5.3-codex":        true,
	"gpt-5.2-codex":        true,
	"gpt-5.2":              true,
	"gpt-5.1-codex-max":    true,
	"gpt-5.1-codex":        true,
	"gpt-5.1":              true,
	"gpt-5.1-codex-mini":   true,
	"gpt-5-mini":           true,
	"gpt-4.1":              true,
	"gemini-3-pro-preview": true,
}

// String returns a human-readable summary of the options for logging.
func (o Options) String() string {
	model := o.Model
	if model == "" {
		model = "default"
	}
	subagent := o.AgentType
	if subagent == "" {
		subagent = "general-purpose"
	}
	return fmt.Sprintf("model=%s, subagent=%s", model, subagent)
}

// Validate returns warnings for invalid option values.
func (o Options) Validate() []string {
	var warnings []string
	if o.Model != "" && !validModels[o.Model] {
		names := make([]string, 0, len(validModels))
		for k := range validModels {
			names = append(names, k)
		}
		sort.Strings(names)
		warnings = append(warnings, fmt.Sprintf("unknown model %q, must be one of: %s", o.Model, strings.Join(names, ", ")))
	}
	return warnings
}

// parseOptions extracts agent options from YAML frontmatter delimited by "---".
// we only support YAML with "---" delimiters because agent files are our own format —
// no need for TOML/JSON/multi-format support that libraries like adrg/frontmatter provide.
// CutPrefix + Cut handle delimiter splitting without index arithmetic.
// returns parsed options and body. if no frontmatter, returns zero value and original content.
func parseOptions(content string) (Options, string) {
	after, found := strings.CutPrefix(content, "---\n")
	if !found {
		return Options{}, content
	}

	header, body, found := strings.Cut(after, "\n---")
	if !found {
		return Options{}, content
	}
	// closing delimiter must be on its own line
	if body != "" && body[0] != '\n' {
		return Options{}, content
	}

	var opts Options
	if err := yaml.Unmarshal([]byte(header), &opts); err != nil {
		return Options{}, content // malformed YAML → treat as no frontmatter
	}

	return opts, strings.TrimSpace(body)
}
