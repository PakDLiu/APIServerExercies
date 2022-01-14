package metadatahandlers

import (
	"APIServerExercise/core"
	mock_search "APIServerExercise/mock/search"
	"APIServerExercise/util"
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var testMetadata core.Metadata

func setupTest() {
	website, _ := url.Parse("https://website.com")
	source, _ := url.Parse("https://github.com/random/repo")
	testMetadata = core.Metadata{
		Id:      uuid.New(),
		Title:   "Valid App 1",
		Version: "0.0.1",
		Maintainers: []*core.Maintainer{
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

// region HandleMetadataPutWithId

func TestMetadataHandlerManager_HandleMetadataPutWithId(t *testing.T) {
	setupTest()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var buf bytes.Buffer
	err := yaml.NewEncoder(&buf).Encode(testMetadata)
	assert.Nil(t, err)

	expectedBytes, err := yaml.Marshal(testMetadata)
	assert.Nil(t, err)

	request := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/metadata/%s", testMetadata.Id.String()),
		&buf)
	request = mux.SetURLVars(request, map[string]string{
		"id": testMetadata.Id.String(),
	})
	responseRecorder := httptest.NewRecorder()

	mockIndexer := mock_search.NewMockIndexer(ctrl)
	mockIndexer.
		EXPECT().
		AddToIndex(
			gomock.Eq(&testMetadata),
			testMetadata.Id,
			"").
		Times(1)

	manager := MetadataHandlerManager{
		Database: &core.Database{Metadatas: map[uuid.UUID]core.Metadata{}},
		Indexer:  mockIndexer,
	}
	manager.HandleMetadataPutWithId(responseRecorder, request)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code)
	assert.Equal(t, string(expectedBytes), responseRecorder.Body.String())
	assert.Len(t, manager.Database.Metadatas, 1)
	assert.Equal(t, manager.Database.Metadatas[testMetadata.Id], testMetadata)
}

func TestMetadataHandlerManager_HandleMetadataPutWithId_DifferentIds(t *testing.T) {
	setupTest()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pathId := uuid.New()

	var buf bytes.Buffer
	err := yaml.NewEncoder(&buf).Encode(testMetadata)
	assert.Nil(t, err)

	testMetadata.Id = pathId
	expectedBytes, err := yaml.Marshal(testMetadata)
	assert.Nil(t, err)

	request := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/metadata/%s", pathId.String()),
		&buf)
	request = mux.SetURLVars(request, map[string]string{
		"id": pathId.String(),
	})
	responseRecorder := httptest.NewRecorder()

	mockIndexer := mock_search.NewMockIndexer(ctrl)
	mockIndexer.
		EXPECT().
		AddToIndex(
			gomock.Eq(&testMetadata),
			pathId,
			"").
		Times(1)

	manager := MetadataHandlerManager{
		Database: &core.Database{Metadatas: map[uuid.UUID]core.Metadata{}},
		Indexer:  mockIndexer,
	}
	manager.HandleMetadataPutWithId(responseRecorder, request)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code)
	assert.Equal(t, string(expectedBytes), responseRecorder.Body.String())
	assert.Len(t, manager.Database.Metadatas, 1)
	assert.Equal(t, manager.Database.Metadatas[pathId], testMetadata)
}

func TestMetadataHandlerManager_HandleMetadataPutWithId_UpdateMetadata(t *testing.T) {
	setupTest()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	oldWebsite, _ := url.Parse("https://www.microsoft.com")
	oldMetadata := core.Metadata{
		Id:          testMetadata.Id,
		Title:       "old title",
		Version:     "0.0.0",
		Maintainers: nil,
		Company:     "old company",
		Website:     util.Yamlurl{oldWebsite},
		Source:      util.Yamlurl{oldWebsite},
		License:     "old license",
		Description: "old description",
	}

	var buf bytes.Buffer
	err := yaml.NewEncoder(&buf).Encode(testMetadata)
	assert.Nil(t, err)

	expectedBytes, err := yaml.Marshal(testMetadata)
	assert.Nil(t, err)

	request := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/metadata/%s", testMetadata.Id.String()),
		&buf)
	request = mux.SetURLVars(request, map[string]string{
		"id": testMetadata.Id.String(),
	})
	responseRecorder := httptest.NewRecorder()

	mockIndexer := mock_search.NewMockIndexer(ctrl)
	mockIndexer.
		EXPECT().
		AddToIndex(
			gomock.Eq(&testMetadata),
			testMetadata.Id,
			"").
		Times(1)
	mockIndexer.EXPECT().RemoveFromIndex(testMetadata.Id).Times(1)

	manager := MetadataHandlerManager{
		Database: &core.Database{Metadatas: map[uuid.UUID]core.Metadata{
			testMetadata.Id: oldMetadata,
		}},
		Indexer: mockIndexer,
	}
	manager.HandleMetadataPutWithId(responseRecorder, request)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code)
	assert.Equal(t, string(expectedBytes), responseRecorder.Body.String())
	assert.Len(t, manager.Database.Metadatas, 1)
	assert.Equal(t, manager.Database.Metadatas[testMetadata.Id], testMetadata)
}

func TestMetadataHandlerManager_HandleMetadataPutWithId_InvalidPathId(t *testing.T) {
	invalidId := "badId"

	request := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/metadata/%s", invalidId),
		nil)
	request = mux.SetURLVars(request, map[string]string{
		"id": invalidId,
	})
	responseRecorder := httptest.NewRecorder()

	manager := MetadataHandlerManager{}
	manager.HandleMetadataPutWithId(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "Error parsing ID: invalid UUID length:")
}

func TestMetadataHandlerManager_HandleMetadataPutWithId_InvalidPayload(t *testing.T) {
	setupTest()
	testMetadata.Title = ""

	var buf bytes.Buffer
	err := yaml.NewEncoder(&buf).Encode(testMetadata)
	assert.Nil(t, err)

	request := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/metadata/%s", testMetadata.Id.String()),
		&buf)
	request = mux.SetURLVars(request, map[string]string{
		"id": testMetadata.Id.String(),
	})
	responseRecorder := httptest.NewRecorder()

	manager := MetadataHandlerManager{}
	manager.HandleMetadataPutWithId(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "Validation failed")
}

// endregion

// region HandleMetadataPut

func TestMetadataHandlerManager_HandleMetadataPut_WithIdInPayload(t *testing.T) {
	setupTest()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var buf bytes.Buffer
	err := yaml.NewEncoder(&buf).Encode(testMetadata)
	assert.Nil(t, err)

	expectedBytes, err := yaml.Marshal(testMetadata)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, "/metadata", &buf)
	responseRecorder := httptest.NewRecorder()

	mockIndexer := mock_search.NewMockIndexer(ctrl)
	mockIndexer.
		EXPECT().
		AddToIndex(
			gomock.Eq(&testMetadata),
			testMetadata.Id,
			"").
		Times(1)

	manager := MetadataHandlerManager{
		Database: &core.Database{Metadatas: map[uuid.UUID]core.Metadata{}},
		Indexer:  mockIndexer,
	}
	manager.HandleMetadataPut(responseRecorder, request)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code)
	assert.Equal(t, string(expectedBytes), responseRecorder.Body.String())
	assert.Len(t, manager.Database.Metadatas, 1)
	assert.Equal(t, manager.Database.Metadatas[testMetadata.Id], testMetadata)
}

func TestMetadataHandlerManager_HandleMetadataPut_WithoutId(t *testing.T) {
	setupTest()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testMetadata.Id = uuid.UUID{}

	var buf bytes.Buffer
	err := yaml.NewEncoder(&buf).Encode(testMetadata)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, "/metadata", &buf)
	responseRecorder := httptest.NewRecorder()

	mockIndexer := mock_search.NewMockIndexer(ctrl)
	mockIndexer.
		EXPECT().
		AddToIndex(
			gomock.Any(),
			gomock.Any(),
			"").
		Do(func(actual *core.Metadata, id uuid.UUID, prefix string) {
			testMetadata.Id = id
			assert.Equal(t, &testMetadata, actual)
		}).
		Times(1)

	manager := MetadataHandlerManager{
		Database: &core.Database{Metadatas: map[uuid.UUID]core.Metadata{}},
		Indexer:  mockIndexer,
	}
	manager.HandleMetadataPut(responseRecorder, request)

	expectedBytes, err := yaml.Marshal(testMetadata)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code)
	assert.Equal(t, string(expectedBytes), responseRecorder.Body.String())
	assert.Len(t, manager.Database.Metadatas, 1)
	assert.Equal(t, manager.Database.Metadatas[testMetadata.Id], testMetadata)
}

func TestMetadataHandlerManager_HandleMetadataPut_InvalidPayload(t *testing.T) {
	setupTest()
	testMetadata.Title = ""

	var buf bytes.Buffer
	err := yaml.NewEncoder(&buf).Encode(testMetadata)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, "/metadata", &buf)
	responseRecorder := httptest.NewRecorder()

	manager := MetadataHandlerManager{}
	manager.HandleMetadataPut(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "Validation failed")
}

// endregion
