package metadatahandlers

import (
	"APIServerExercise/core"
	mock_search "APIServerExercise/mock/search"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"net/http"
	"net/http/httptest"
	"testing"
)

// region HandleMetadataGetWithId

func TestMetadataHandlerManager_HandleMetadataGetWithId(t *testing.T) {
	setupTest()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	database := &core.Database{Metadatas: map[uuid.UUID]*core.Metadata{
		testMetadata.Id: testMetadata,
	}}

	request := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/metadata/%s", testMetadata.Id.String()),
		nil)
	request = mux.SetURLVars(request, map[string]string{
		"id": testMetadata.Id.String(),
	})
	responseRecorder := httptest.NewRecorder()

	mockFilterer := mock_search.NewMockFilterer(ctrl)
	mockFilterer.
		EXPECT().
		FilterMetadata(gomock.Any(), database).
		DoAndReturn(func(query map[string][]string, database *core.Database) ([]*core.Metadata, error) {
			assert.Len(t, query["id"], 1)
			assert.Equal(t, testMetadata.Id.String(), query["id"][0])
			return []*core.Metadata{testMetadata}, nil
		}).
		Times(1)

	manager := MetadataHandlerManager{
		Database: database,
		Filterer: mockFilterer,
	}
	manager.HandleMetadataGetWithId(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var actual []core.Metadata
	err := yaml.Unmarshal(responseRecorder.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Len(t, actual, 1)
	assert.Equal(t, &actual[0], testMetadata)
}

func TestMetadataHandlerManager_HandleMetadataGetWithId_WithInvalidId(t *testing.T) {
	invalidId := "badId"

	request := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/metadata/%s", invalidId),
		nil)
	request = mux.SetURLVars(request, map[string]string{
		"id": invalidId,
	})
	responseRecorder := httptest.NewRecorder()

	manager := MetadataHandlerManager{}
	manager.HandleMetadataGetWithId(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "Error parsing ID: invalid UUID length:")
}

func TestMetadataHandlerManager_HandleMetadataGetWithId_WithFilterError(t *testing.T) {
	setupTest()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testError := fmt.Errorf("test error")

	database := &core.Database{Metadatas: map[uuid.UUID]*core.Metadata{
		testMetadata.Id: testMetadata,
	}}

	request := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/metadata/%s", testMetadata.Id.String()),
		nil)
	request = mux.SetURLVars(request, map[string]string{
		"id": testMetadata.Id.String(),
	})
	responseRecorder := httptest.NewRecorder()

	mockFilterer := mock_search.NewMockFilterer(ctrl)
	mockFilterer.
		EXPECT().
		FilterMetadata(gomock.Any(), database).
		DoAndReturn(func(query map[string][]string, database *core.Database) ([]core.Metadata, error) {
			assert.Len(t, query["id"], 1)
			assert.Equal(t, testMetadata.Id.String(), query["id"][0])
			return nil, testError
		}).
		Times(1)

	manager := MetadataHandlerManager{
		Database: database,
		Filterer: mockFilterer,
	}
	manager.HandleMetadataGetWithId(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), testError.Error())
}

// endregion

// region HandleMetadataGet

func TestMetadataHandlerManager_HandleMetadataGet(t *testing.T) {
	setupTest()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	database := &core.Database{Metadatas: map[uuid.UUID]*core.Metadata{
		testMetadata.Id: testMetadata,
	}}

	request := httptest.NewRequest(http.MethodGet, "/metadata", nil)
	responseRecorder := httptest.NewRecorder()

	mockFilterer := mock_search.NewMockFilterer(ctrl)
	mockFilterer.
		EXPECT().
		FilterMetadata(gomock.Any(), database).
		DoAndReturn(func(query map[string][]string, database *core.Database) ([]*core.Metadata, error) {
			assert.Empty(t, query)
			return []*core.Metadata{testMetadata}, nil
		}).
		Times(1)

	manager := MetadataHandlerManager{
		Database: database,
		Filterer: mockFilterer,
	}
	manager.HandleMetadataGet(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var actual []core.Metadata
	err := yaml.Unmarshal(responseRecorder.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Len(t, actual, 1)
	assert.Equal(t, &actual[0], testMetadata)
}

func TestMetadataHandlerManager_HandleMetadataGet_WithFilterError(t *testing.T) {
	setupTest()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testError := fmt.Errorf("test error")

	database := &core.Database{Metadatas: map[uuid.UUID]*core.Metadata{
		testMetadata.Id: testMetadata,
	}}

	request := httptest.NewRequest(http.MethodGet, "/metadata", nil)
	responseRecorder := httptest.NewRecorder()

	mockFilterer := mock_search.NewMockFilterer(ctrl)
	mockFilterer.
		EXPECT().
		FilterMetadata(gomock.Any(), database).
		DoAndReturn(func(query map[string][]string, database *core.Database) ([]core.Metadata, error) {
			assert.Empty(t, query)
			return nil, testError
		}).
		Times(1)

	manager := MetadataHandlerManager{
		Database: database,
		Filterer: mockFilterer,
	}
	manager.HandleMetadataGet(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), testError.Error())
}

// endregion
