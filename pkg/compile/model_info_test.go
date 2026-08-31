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
	}, nil)

	suite.Equal("Organization", info.ModelName)
	suite.Equal("OrganizationRepository", info.RepoName)
}

func (suite *ModelInfoTestSuite) TestExtractRepoInfo_RelFieldForOne() {
	allModels := map[string]yaml.Model{
		"Task": {
			Name: "Task",
			Fields: map[string]yaml.ModelField{
				"ID": {Type: "AutoIncrement"},
			},
			Identifiers: map[string]yaml.ModelIdentifier{
				"primary": {Fields: []string{"ID"}},
			},
		},
	}
	info := compile.ExtractRepoInfo(yaml.Model{
		Name: "TaskTag",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: "UUID"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
			"taskTag": {Fields: []string{"rel:Task"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Task": {Type: "ForOne"},
		},
	}, allModels)

	suite.Require().Len(info.Identifiers, 2)
	taskTagID := info.Identifiers[1]
	suite.Equal("taskTag", taskTagID.Name)
	suite.Require().Len(taskTagID.Fields, 1)
	suite.Equal("TaskID", taskTagID.Fields[0].Name)
	suite.Equal(yaml.ModelFieldType("AutoIncrement"), taskTagID.Fields[0].Type)
}

func (suite *ModelInfoTestSuite) TestExtractRepoInfo_RelFieldComposite() {
	allModels := map[string]yaml.Model{
		"Task": {
			Name: "Task",
			Fields: map[string]yaml.ModelField{
				"ID": {Type: "UUID"},
			},
			Identifiers: map[string]yaml.ModelIdentifier{
				"primary": {Fields: []string{"ID"}},
			},
		},
		"Tag": {
			Name: "Tag",
			Fields: map[string]yaml.ModelField{
				"ID": {Type: "UUID"},
			},
			Identifiers: map[string]yaml.ModelIdentifier{
				"primary": {Fields: []string{"ID"}},
			},
		},
	}
	info := compile.ExtractRepoInfo(yaml.Model{
		Name: "TaskTag",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: "AutoIncrement"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
			"taskTag": {Fields: []string{"rel:Task", "rel:Tag"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Task": {Type: "ForOne"},
			"Tag":  {Type: "ForOne"},
		},
	}, allModels)

	suite.Require().Len(info.Identifiers, 2)
	taskTagID := info.Identifiers[1]
	suite.Equal("taskTag", taskTagID.Name)
	suite.Require().Len(taskTagID.Fields, 2)
	suite.Equal("TaskID", taskTagID.Fields[0].Name)
	suite.Equal(yaml.ModelFieldType("UUID"), taskTagID.Fields[0].Type)
	suite.Equal("TagID", taskTagID.Fields[1].Name)
	suite.Equal(yaml.ModelFieldType("UUID"), taskTagID.Fields[1].Type)
}

func (suite *ModelInfoTestSuite) TestExtractRepoInfo_RelFieldForOnePoly() {
	info := compile.ExtractRepoInfo(yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: "UUID"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary":   {Fields: []string{"ID"}},
			"commentable": {Fields: []string{"rel:Commentable"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {Type: "ForOnePoly", For: []string{"Task", "Project"}},
		},
	}, nil)

	suite.Require().Len(info.Identifiers, 2)
	commentableID := info.Identifiers[1]
	suite.Equal("commentable", commentableID.Name)
	suite.Require().Len(commentableID.Fields, 2)
	suite.Equal("CommentableType", commentableID.Fields[0].Name)
	suite.Equal(yaml.ModelFieldType("String"), commentableID.Fields[0].Type)
	suite.Equal("CommentableID", commentableID.Fields[1].Name)
	suite.Equal(yaml.ModelFieldType("String"), commentableID.Fields[1].Type)
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
	}, nil)

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
	}, nil)

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
	}, nil)

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
	}, nil)

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
	}, nil)

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
	}, nil)

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
	}, nil)

	suite.Require().Len(info.Identifiers, 2)
	nameID := info.Identifiers[1]
	suite.Equal("name", nameID.Name)
	suite.Require().Len(nameID.Fields, 2)
	suite.Equal("FirstName", nameID.Fields[0].Name)
	suite.Equal("LastName", nameID.Fields[1].Name)
}
