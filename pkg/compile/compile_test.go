package compile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/go-util/assertfile"
	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-morpherepo/internal/testutils"
	"github.com/kalo-build/plugin-morphe-morpherepo/pkg/compile"
	"github.com/kalo-build/plugin-morphe-morpherepo/pkg/compile/cfg"
)

type CompileTestSuite struct {
	assertfile.FileSuite

	TestDirPath            string
	TestGroundTruthDirPath string

	ModelsDirPath     string
	EnumsDirPath      string
	StructuresDirPath string
	EntitiesDirPath   string
}

func TestCompileTestSuite(t *testing.T) {
	suite.Run(t, new(CompileTestSuite))
}

func (suite *CompileTestSuite) SetupTest() {
	suite.TestDirPath = testutils.GetTestDirPath()
	suite.TestGroundTruthDirPath = filepath.Join(suite.TestDirPath, "ground-truth", "compile-minimal")

	suite.ModelsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "models")
	suite.EnumsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "enums")
	suite.StructuresDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "structures")
	suite.EntitiesDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "entities")
}

func (suite *CompileTestSuite) TearDownTest() {
	suite.TestDirPath = ""
}

func (suite *CompileTestSuite) TestMorpheToMorpheRepo() {
	workingDirPath := filepath.Join(suite.TestDirPath, "working")
	suite.Nil(os.Mkdir(workingDirPath, 0755))
	defer os.RemoveAll(workingDirPath)

	config := cfg.CompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      suite.EnumsDirPath,
			RegistryModelsDirPath:     suite.ModelsDirPath,
			RegistryStructuresDirPath: suite.StructuresDirPath,
			RegistryEntitiesDirPath:   suite.EntitiesDirPath,
		},
		OutputDirPath: workingDirPath,
	}

	compileErr := compile.MorpheToMorpheRepo(config)

	suite.NoError(compileErr)

	// Verify organization.repo
	orgPath := filepath.Join(workingDirPath, "organization.repo")
	gtOrgPath := filepath.Join(suite.TestGroundTruthDirPath, "organization.repo")
	suite.FileExists(orgPath)
	suite.FileEquals(orgPath, gtOrgPath)

	// Verify project.repo
	projectPath := filepath.Join(workingDirPath, "project.repo")
	gtProjectPath := filepath.Join(suite.TestGroundTruthDirPath, "project.repo")
	suite.FileExists(projectPath)
	suite.FileEquals(projectPath, gtProjectPath)

	// Verify task.repo
	taskPath := filepath.Join(workingDirPath, "task.repo")
	gtTaskPath := filepath.Join(suite.TestGroundTruthDirPath, "task.repo")
	suite.FileExists(taskPath)
	suite.FileEquals(taskPath, gtTaskPath)
}

func (suite *CompileTestSuite) TestMorpheToMorpheRepo_ExcludeModels() {
	workingDirPath := filepath.Join(suite.TestDirPath, "working-exclude")
	suite.Nil(os.Mkdir(workingDirPath, 0755))
	defer os.RemoveAll(workingDirPath)

	config := cfg.CompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      suite.EnumsDirPath,
			RegistryModelsDirPath:     suite.ModelsDirPath,
			RegistryStructuresDirPath: suite.StructuresDirPath,
			RegistryEntitiesDirPath:   suite.EntitiesDirPath,
		},
		ExcludeModels: []string{"Task"},
		OutputDirPath: workingDirPath,
	}

	compileErr := compile.MorpheToMorpheRepo(config)

	suite.NoError(compileErr)

	// Organization and Project should exist
	suite.FileExists(filepath.Join(workingDirPath, "organization.repo"))
	suite.FileExists(filepath.Join(workingDirPath, "project.repo"))

	// Task should NOT exist
	_, taskErr := os.Stat(filepath.Join(workingDirPath, "task.repo"))
	suite.True(os.IsNotExist(taskErr))
}
