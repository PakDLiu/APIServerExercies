package core

import (
	"APIServerExercies/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

var testMetadata Metadata

func setupTest() {
	website, _ := url.Parse("https://website.com")
	source, _ := url.Parse("https://github.com/random/repo")
	testMetadata = Metadata{
		Id:      uuid.New(),
		Title:   "Valid App 1",
		Version: "0.0.1",
		Maintainers: []*Maintainer{
			{
				Name:  "firstmaintainer app1",
				Email: "firstmaintainer@hotmail.com",
			},
			{
				Name:  "secondmaintainer app1",
				Email: "secondmaintainer@gmail.com",
			},
		},
		Company:     "Random Inc.",
		Website:     util.Yamlurl{URL: website},
		Source:      util.Yamlurl{URL: source},
		License:     "Apache-2.0",
		Description: "### Interesting Title\n Some application content, and description",
	}
}

func TestValidateStruct(t *testing.T) {
	setupTest()
	err := ValidateStruct(testMetadata)
	assert.Nil(t, err)
}

func TestValidateStruct_WithMissingField(t *testing.T) {
	setupTest()
	testMetadata.Title = ""
	err := ValidateStruct(testMetadata)
	assert.Error(t, err)
	assert.Equal(t, "Key: 'Metadata.Title' Error:Field validation for 'Title' failed on the 'required' tag", err.Error())
}

func TestValidateStruct_InvalidEmail(t *testing.T) {
	setupTest()
	testMetadata.Maintainers[0].Email = "invalid@@@email.com"
	err := ValidateStruct(testMetadata)
	assert.Error(t, err)
	assert.Equal(t, "Key: 'Metadata.Maintainers[0].Email' Error:Field validation for 'Email' failed on the 'email' tag", err.Error())
}

func TestValidateStruct_NotEnoughMaintainers(t *testing.T) {
	setupTest()
	testMetadata.Maintainers = []*Maintainer{}
	err := ValidateStruct(testMetadata)
	assert.Error(t, err)
	assert.Equal(t, "Key: 'Metadata.Maintainers' Error:Field validation for 'Maintainers' failed on the 'gt' tag", err.Error())
}
