package compile

import (
	"sort"

	"github.com/kalo-build/morphe-go/pkg/yaml"
)

// RepoIdentifierField describes a single field within an identifier.
type RepoIdentifierField struct {
	Name string
	Type yaml.ModelFieldType
}

// RepoIdentifier describes an identifier (primary or secondary).
type RepoIdentifier struct {
	Name   string // e.g., "primary", "code", "email"
	Fields []RepoIdentifierField
}

// RepoFilter describes a filter parameter for GetAll.
type RepoFilter struct {
	Name     string // camelCase parameter name, e.g., "organizationID"
	Type     string // Morphe type string, e.g., "UUID"
	Relation string // source relationship name, e.g., "Organization"
}

// RepoInfo holds all derived repository information for a model.
type RepoInfo struct {
	ModelName   string
	RepoName    string
	Identifiers []RepoIdentifier
	Filters     []RepoFilter
}

// ExtractRepoInfo derives repository contract information from a Morphe model.
func ExtractRepoInfo(model yaml.Model) RepoInfo {
	info := RepoInfo{
		ModelName: model.Name,
		RepoName:  model.Name + "Repository",
	}

	// Extract identifiers (sorted for deterministic output)
	idNames := sortedKeys(model.Identifiers)
	for _, idName := range idNames {
		id := model.Identifiers[idName]
		repoID := RepoIdentifier{Name: idName}
		for _, fieldName := range id.Fields {
			fieldType := yaml.ModelFieldType("String") // default
			if field, ok := model.Fields[fieldName]; ok {
				fieldType = field.Type
			}
			repoID.Fields = append(repoID.Fields, RepoIdentifierField{
				Name: fieldName,
				Type: fieldType,
			})
		}
		info.Identifiers = append(info.Identifiers, repoID)
	}

	// Sort identifiers: primary first, then alphabetical
	sort.SliceStable(info.Identifiers, func(i, j int) bool {
		if info.Identifiers[i].Name == "primary" {
			return true
		}
		if info.Identifiers[j].Name == "primary" {
			return false
		}
		return info.Identifiers[i].Name < info.Identifiers[j].Name
	})

	// Extract filters from ForOne/ForOnePoly relationships
	relNames := sortedKeys(model.Related)
	for _, relName := range relNames {
		rel := model.Related[relName]
		switch rel.Type {
		case "ForOne":
			info.Filters = append(info.Filters, RepoFilter{
				Name:     lowerFirst(relName) + "ID",
				Type:     "UUID", // FK references are UUID by convention
				Relation: relName,
			})
		case "ForOnePoly":
			filterName := relName
			if rel.Through != "" {
				filterName = rel.Through
			}
			info.Filters = append(info.Filters, RepoFilter{
				Name:     lowerFirst(filterName) + "ID",
				Type:     "UUID",
				Relation: filterName,
			})
		}
	}

	return info
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	if runes[0] >= 'A' && runes[0] <= 'Z' {
		runes[0] = runes[0] + 32
	}
	return string(runes)
}
