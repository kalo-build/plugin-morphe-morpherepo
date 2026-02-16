package compile_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-morpherepo/pkg/compile"
)

type GenerateRepoTestSuite struct {
	suite.Suite
}

func TestGenerateRepoTestSuite(t *testing.T) {
	suite.Run(t, new(GenerateRepoTestSuite))
}

func (suite *GenerateRepoTestSuite) TestGenerateRepoYAML_SimpleModel() {
	info := compile.ExtractRepoInfo(yaml.Model{
		Name: "Organization",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: "UUID"},
			"Code": {Type: "String"},
			"Name": {Type: "String"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
			"code":    {Fields: []string{"Code"}},
		},
	})

	result := compile.GenerateRepoYAML(info)

	suite.Contains(result, "name: OrganizationRepository")
	suite.Contains(result, "model: Organization")
	suite.Contains(result, "primary:")
	suite.Contains(result, "- name: ID")
	suite.Contains(result, "type: UUID")
	suite.Contains(result, "code:")
	suite.Contains(result, "- name: Code")
	suite.Contains(result, "filters: []")
	suite.Contains(result, "list: true")
	suite.Contains(result, "get: true")
	suite.Contains(result, "create: true")
	suite.Contains(result, "update: true")
	suite.Contains(result, "delete: true")
}

func (suite *GenerateRepoTestSuite) TestGenerateRepoYAML_ModelWithForOne() {
	info := compile.ExtractRepoInfo(yaml.Model{
		Name: "Project",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: "UUID"},
			"Code": {Type: "String"},
			"Name": {Type: "String"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
			"code":    {Fields: []string{"Code"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Organization": {Type: "ForOne"},
		},
	})

	result := compile.GenerateRepoYAML(info)

	suite.Contains(result, "name: ProjectRepository")
	suite.Contains(result, "model: Project")
	suite.Contains(result, "filters:")
	suite.Contains(result, "- name: organizationID")
	suite.Contains(result, "type: UUID")
	suite.Contains(result, "relation: Organization")
}

func (suite *GenerateRepoTestSuite) TestGenerateRepoYAML_ModelNoSecondaryIdentifiers() {
	info := compile.ExtractRepoInfo(yaml.Model{
		Name: "Task",
		Fields: map[string]yaml.ModelField{
			"ID":     {Type: "UUID"},
			"Title":  {Type: "String"},
			"Status": {Type: "String"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Project": {Type: "ForOne"},
		},
	})

	result := compile.GenerateRepoYAML(info)

	suite.Contains(result, "name: TaskRepository")
	suite.Contains(result, "model: Task")
	suite.Contains(result, "primary:")
	suite.NotContains(result, "code:")
	suite.Contains(result, "- name: projectID")
	suite.Contains(result, "relation: Project")
}

func (suite *GenerateRepoTestSuite) TestGenerateRepoYAML_HasOneHasManyNotInFilters() {
	info := compile.ExtractRepoInfo(yaml.Model{
		Name: "Organization",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: "UUID"},
			"Name": {Type: "String"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Project": {Type: "HasMany"},
			"Profile": {Type: "HasOne"},
		},
	})

	result := compile.GenerateRepoYAML(info)

	suite.Contains(result, "filters: []")
}
