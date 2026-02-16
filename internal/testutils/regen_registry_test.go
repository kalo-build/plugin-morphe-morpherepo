package testutils_test

import (
	"path/filepath"
	"testing"

	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-morpherepo/pkg/compile"
	"github.com/kalo-build/plugin-morphe-morpherepo/pkg/compile/cfg"
)

func TestRegenRegistry(t *testing.T) {
	t.Skip("Run manually to regenerate registry .repo files")

	registryRoot := filepath.Join("..", "..", "..", "kalo-plugin-registry", "morphe", "registry")
	outputPath := filepath.Join("..", "..", "..", "kalo-plugin-registry", "morphe", "repo")

	config := cfg.CompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      filepath.Join(registryRoot, "enums"),
			RegistryModelsDirPath:     filepath.Join(registryRoot, "models"),
			RegistryStructuresDirPath: filepath.Join(registryRoot, "structures"),
			RegistryEntitiesDirPath:   filepath.Join(registryRoot, "entities"),
		},
		ExcludeModels: []string{"User"},
		OutputDirPath: outputPath,
	}

	if err := compile.MorpheToMorpheRepo(config); err != nil {
		t.Fatal(err)
	}

	t.Log("Registry .repo files generated at:", outputPath)
}
