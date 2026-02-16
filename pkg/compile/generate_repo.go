package compile

import (
	"fmt"
	"strings"
)

// GenerateRepoYAML generates a .repo YAML file content for a model.
func GenerateRepoYAML(info RepoInfo) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("name: %s\n", info.RepoName))
	b.WriteString(fmt.Sprintf("model: %s\n", info.ModelName))
	b.WriteString("\n")

	// Identifiers
	b.WriteString("identifiers:\n")
	for _, id := range info.Identifiers {
		b.WriteString(fmt.Sprintf("  %s:\n", id.Name))
		b.WriteString("    fields:\n")
		for _, field := range id.Fields {
			b.WriteString(fmt.Sprintf("      - name: %s\n", field.Name))
			b.WriteString(fmt.Sprintf("        type: %s\n", string(field.Type)))
		}
	}
	b.WriteString("\n")

	// Filters
	if len(info.Filters) == 0 {
		b.WriteString("filters: []\n")
	} else {
		b.WriteString("filters:\n")
		for _, filter := range info.Filters {
			b.WriteString(fmt.Sprintf("  - name: %s\n", filter.Name))
			b.WriteString(fmt.Sprintf("    type: %s\n", filter.Type))
			b.WriteString(fmt.Sprintf("    relation: %s\n", filter.Relation))
		}
	}
	b.WriteString("\n")

	// Operations
	b.WriteString("operations:\n")
	b.WriteString("  list: true\n")
	b.WriteString("  get: true\n")
	b.WriteString("  create: true\n")
	b.WriteString("  update: true\n")
	b.WriteString("  delete: true\n")

	return b.String()
}
