package testutils_test

import (
	"path/filepath"
	"testing"

	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-morpherepo/internal/testutils"
	"github.com/kalo-build/plugin-morphe-morpherepo/pkg/compile"
	"github.com/kalo-build/plugin-morphe-morpherepo/pkg/compile/cfg"
)

func TestGenerateGroundTruth(t *testing.T) {
	t.Skip("Only run manually to regenerate ground truth files")

	testDirPath := testutils.GetTestDirPath()

	registryPath := filepath.Join(testDirPath, "registry", "minimal")
	outputPath := filepath.Join(testDirPath, "ground-truth", "compile-minimal")

	config := cfg.CompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      filepath.Join(registryPath, "enums"),
			RegistryModelsDirPath:     filepath.Join(registryPath, "models"),
			RegistryStructuresDirPath: filepath.Join(registryPath, "structures"),
			RegistryEntitiesDirPath:   filepath.Join(registryPath, "entities"),
		},
		OutputDirPath: outputPath,
	}

	if err := compile.MorpheToMorpheRepo(config); err != nil {
		t.Fatal(err)
	}

	t.Log("Ground truth files regenerated at:", outputPath)
}
