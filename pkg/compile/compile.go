package compile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/plugin-morphe-morpherepo/pkg/compile/cfg"
)

// MorpheToMorpheRepo is the main compilation entrypoint.
// It reads Morphe models and generates .repo YAML files.
func MorpheToMorpheRepo(config cfg.CompileConfig) error {
	// Load Morphe registry
	r, rErr := registry.LoadMorpheRegistry(registry.LoadMorpheRegistryHooks{}, config.MorpheLoadRegistryConfig)
	if rErr != nil {
		return fmt.Errorf("failed to load morphe registry: %w", rErr)
	}

	// Build exclude set
	excludeSet := make(map[string]bool)
	for _, name := range config.ExcludeModels {
		excludeSet[name] = true
	}

	// Get all models sorted
	allModels := r.GetAllModels()
	modelNames := make([]string, 0, len(allModels))
	for name := range allModels {
		modelNames = append(modelNames, name)
	}
	sort.Strings(modelNames)

	// Create output directory
	if err := os.MkdirAll(config.OutputDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate .repo file for each model
	generated := 0
	for _, name := range modelNames {
		if excludeSet[name] {
			continue
		}
		model := allModels[name]
		info := ExtractRepoInfo(model)
		content := GenerateRepoYAML(info)

		fileName := toSnakeCase(name) + ".repo"
		filePath := filepath.Join(config.OutputDirPath, fileName)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write repo file for %s: %w", name, err)
		}
		generated++
	}

	if generated == 0 {
		return fmt.Errorf("no models to generate (all excluded or none found)")
	}

	return nil
}

// toSnakeCase converts PascalCase to snake_case.
func toSnakeCase(s string) string {
	if s == "" {
		return s
	}
	var words []string
	wordStart := 0
	runes := []rune(s)
	for i := 1; i < len(runes); i++ {
		if runes[i] >= 'A' && runes[i] <= 'Z' {
			if runes[i-1] >= 'a' && runes[i-1] <= 'z' {
				words = append(words, string(runes[wordStart:i]))
				wordStart = i
			} else if i+1 < len(runes) && runes[i+1] >= 'a' && runes[i+1] <= 'z' {
				words = append(words, string(runes[wordStart:i]))
				wordStart = i
			}
		}
	}
	words = append(words, string(runes[wordStart:]))
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, "_")
}
