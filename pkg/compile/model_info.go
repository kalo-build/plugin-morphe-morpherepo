package compile

import (
	"sort"
	"strings"

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
// allModels is required to resolve rel: prefixed identifier fields to their
// target model's primary key type.
func ExtractRepoInfo(model yaml.Model, allModels map[string]yaml.Model) RepoInfo {
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
			resolved := resolveIdentifierField(fieldName, model, allModels)
			repoID.Fields = append(repoID.Fields, resolved...)
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

func resolveIdentifierField(fieldName string, model yaml.Model, allModels map[string]yaml.Model) []RepoIdentifierField {
	if !strings.HasPrefix(fieldName, "rel:") {
		fieldType := yaml.ModelFieldType("String")
		if field, ok := model.Fields[fieldName]; ok {
			fieldType = field.Type
		}
		return []RepoIdentifierField{{Name: fieldName, Type: fieldType}}
	}

	relationName := strings.TrimPrefix(fieldName, "rel:")
	relation, hasRelation := model.Related[relationName]
	if !hasRelation {
		return []RepoIdentifierField{{Name: relationName + "ID", Type: "String"}}
	}

	if isPolyFor(relation.Type) {
		return []RepoIdentifierField{
			{Name: relationName + "Type", Type: "String"},
			{Name: relationName + "ID", Type: "String"},
		}
	}

	targetName := relationName
	if relation.Aliased != "" {
		targetName = relation.Aliased
	}

	fkType := resolveTargetPrimaryType(targetName, allModels)
	return []RepoIdentifierField{{Name: relationName + "ID", Type: fkType}}
}

func resolveTargetPrimaryType(targetModelName string, allModels map[string]yaml.Model) yaml.ModelFieldType {
	targetModel, exists := allModels[targetModelName]
	if !exists {
		return "UUID"
	}
	primaryID, hasPrimary := targetModel.Identifiers["primary"]
	if !hasPrimary || len(primaryID.Fields) == 0 {
		return "UUID"
	}
	primaryFieldName := primaryID.Fields[0]
	if strings.HasPrefix(primaryFieldName, "rel:") {
		return "UUID"
	}
	if field, ok := targetModel.Fields[primaryFieldName]; ok {
		return field.Type
	}
	return "UUID"
}

func isPolyFor(relationType string) bool {
	lower := strings.ToLower(relationType)
	return strings.HasPrefix(lower, "for") && strings.HasSuffix(lower, "poly")
}

