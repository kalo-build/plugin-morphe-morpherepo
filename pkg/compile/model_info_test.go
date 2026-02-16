package compile_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-morpherepo/pkg/compile"
)

type ModelInfoTestSuite struct {
	suite.Suite
}

func TestModelInfoTestSuite(t *testing.T) {
	suite.Run(t, new(ModelInfoTestSuite))
}

func (suite *ModelInfoTestSuite) TestExtractRepoInfo_Name() {
	info := compile.ExtractRepoInfo(yaml.Model{
		Name: "Organization",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: "UUID"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	})

	suite.Equal("Organization", info.ModelName)
	suite.Equal("OrganizationRepository", info.RepoName)
}

func (suite *ModelInfoTestSuite) TestExtractRepoInfo_PrimaryIdentifier() {
	info := compile.ExtractRepoInfo(yaml.Model{
		Name: "Organization",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: "UUID"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	})

	suite.Require().Len(info.Identifiers, 1)
	suite.Equal("primary", info.Identifiers[0].Name)
	suite.Require().Len(info.Identifiers[0].Fields, 1)
	suite.Equal("ID", info.Identifiers[0].Fields[0].Name)
	suite.Equal(yaml.ModelFieldType("UUID"), info.Identifiers[0].Fields[0].Type)
}

func (suite *ModelInfoTestSuite) TestExtractRepoInfo_SecondaryIdentifier() {
	info := compile.ExtractRepoInfo(yaml.Model{
		Name: "Organization",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: "UUID"},
			"Code": {Type: "String"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
			"code":    {Fields: []string{"Code"}},
		},
	})

	suite.Require().Len(info.Identifiers, 2)
	suite.Equal("primary", info.Identifiers[0].Name)
	suite.Equal("code", info.Identifiers[1].Name)
	suite.Equal("Code", info.Identifiers[1].Fields[0].Name)
	suite.Equal(yaml.ModelFieldType("String"), info.Identifiers[1].Fields[0].Type)
}

func (suite *ModelInfoTestSuite) TestExtractRepoInfo_PrimaryAlwaysFirst() {
	info := compile.ExtractRepoInfo(yaml.Model{
		Name: "Organization",
		Fields: map[string]yaml.ModelField{
			"ID":    {Type: "UUID"},
			"Code":  {Type: "String"},
			"Email": {Type: "String"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"email":   {Fields: []string{"Email"}},
			"primary": {Fields: []string{"ID"}},
			"code":    {Fields: []string{"Code"}},
		},
	})

	suite.Equal("primary", info.Identifiers[0].Name)
}

func (suite *ModelInfoTestSuite) TestExtractRepoInfo_ForOneFilter() {
	info := compile.ExtractRepoInfo(yaml.Model{
		Name: "Project",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: "UUID"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Organization": {Type: "ForOne"},
		},
	})

	suite.Require().Len(info.Filters, 1)
	suite.Equal("organizationID", info.Filters[0].Name)
	suite.Equal("UUID", info.Filters[0].Type)
	suite.Equal("Organization", info.Filters[0].Relation)
}

func (suite *ModelInfoTestSuite) TestExtractRepoInfo_HasOneHasManyNoFilter() {
	info := compile.ExtractRepoInfo(yaml.Model{
		Name: "Organization",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: "UUID"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Project": {Type: "HasMany"},
			"Profile": {Type: "HasOne"},
		},
	})

	suite.Empty(info.Filters)
}

func (suite *ModelInfoTestSuite) TestExtractRepoInfo_NoRelations() {
	info := compile.ExtractRepoInfo(yaml.Model{
		Name: "Organization",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: "UUID"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	})

	suite.Empty(info.Filters)
}

func (suite *ModelInfoTestSuite) TestExtractRepoInfo_CompositeIdentifier() {
	info := compile.ExtractRepoInfo(yaml.Model{
		Name: "Person",
		Fields: map[string]yaml.ModelField{
			"ID":        {Type: "AutoIncrement"},
			"FirstName": {Type: "String"},
			"LastName":  {Type: "String"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
			"name":    {Fields: []string{"FirstName", "LastName"}},
		},
	})

	suite.Require().Len(info.Identifiers, 2)
	nameID := info.Identifiers[1]
	suite.Equal("name", nameID.Name)
	suite.Require().Len(nameID.Fields, 2)
	suite.Equal("FirstName", nameID.Fields[0].Name)
	suite.Equal("LastName", nameID.Fields[1].Name)
}
