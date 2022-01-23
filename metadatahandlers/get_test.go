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
	"net/url"
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

	manager := MetadataHandlerManager{
		Database: database,
	}
	manager.HandleMetadataGetWithId(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var actual core.Metadata
	err := yaml.Unmarshal(responseRecorder.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, testMetadata, &actual)
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

	var actual core.ResultPage
	err := yaml.Unmarshal(responseRecorder.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Len(t, actual.Resources, 1)
	assert.Equal(t, testMetadata, actual.Resources[0])
	assert.Empty(t, actual.NextLink)
}

func TestMetadataHandlerManager_HandleMetadataGet_WithQuery(t *testing.T) {
	setupTest()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	database := &core.Database{Metadatas: map[uuid.UUID]*core.Metadata{
		testMetadata.Id: testMetadata,
	}}

	request := httptest.NewRequest(http.MethodGet, "/metadata", nil)
	query := url.Values{}
	query["key"] = []string{"value"}
	request.URL.RawQuery = query.Encode()
	responseRecorder := httptest.NewRecorder()

	mockFilterer := mock_search.NewMockFilterer(ctrl)
	mockFilterer.
		EXPECT().
		FilterMetadata(gomock.Any(), database).
		DoAndReturn(func(query map[string][]string, database *core.Database) ([]*core.Metadata, error) {
			assert.Len(t, query, 1)
			assert.Equal(t, "value", query["key"][0])
			return []*core.Metadata{testMetadata}, nil
		}).
		Times(1)

	manager := MetadataHandlerManager{
		Database: database,
		Filterer: mockFilterer,
	}
	manager.HandleMetadataGet(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var actual core.ResultPage
	err := yaml.Unmarshal(responseRecorder.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Len(t, actual.Resources, 1)
	assert.Equal(t, testMetadata, actual.Resources[0])
	assert.Empty(t, actual.NextLink)
}

func TestMetadataHandlerManager_HandleMetadataGet_WithPaging(t *testing.T) {
	setupTest()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	database := &core.Database{Metadatas: map[uuid.UUID]*core.Metadata{
		testMetadata.Id: testMetadata,
	}}

	request := httptest.NewRequest(http.MethodGet, "/metadata", nil)
	query := url.Values{}
	query[offsetParameter] = []string{"0"}
	query[pageSizeParameter] = []string{"1"}
	request.URL.RawQuery = query.Encode()
	responseRecorder := httptest.NewRecorder()

	mockFilterer := mock_search.NewMockFilterer(ctrl)
	mockFilterer.
		EXPECT().
		FilterMetadata(gomock.Any(), database).
		DoAndReturn(func(query map[string][]string, database *core.Database) ([]*core.Metadata, error) {
			assert.Empty(t, query)
			return []*core.Metadata{testMetadata, testMetadata}, nil
		}).
		Times(1)

	manager := MetadataHandlerManager{
		Database: database,
		Filterer: mockFilterer,
	}
	manager.HandleMetadataGet(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var actual core.ResultPage
	err := yaml.Unmarshal(responseRecorder.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Len(t, actual.Resources, 1)
	assert.Equal(t, testMetadata, actual.Resources[0])

	nextLink, err := url.Parse(actual.NextLink)
	assert.Nil(t, err)
	assert.Equal(t, "/metadata", nextLink.Path)
	assert.Equal(t, "1", nextLink.Query().Get(offsetParameter))
	assert.Equal(t, "1", nextLink.Query().Get(pageSizeParameter))
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

// region parsePagingParameters

func TestParsePagingParameters_WithDefaults(t *testing.T) {
	offset, pageSize, err := parsePagingParameters(map[string][]string{})
	assert.Nil(t, err)
	assert.Equal(t, defaultOffset, offset)
	assert.Equal(t, defaultPageSize, pageSize)
}

func TestParsePagingParameters_WithOffset(t *testing.T) {
	expectedOffset := 3
	offset, pageSize, err := parsePagingParameters(
		map[string][]string{
			offsetParameter: {fmt.Sprintf("%d", expectedOffset)},
		})
	assert.Nil(t, err)
	assert.Equal(t, expectedOffset, offset)
	assert.Equal(t, defaultPageSize, pageSize)
}

func TestParsePagingParameters_WithPageSize(t *testing.T) {
	expectedPageSize := 5
	offset, pageSize, err := parsePagingParameters(
		map[string][]string{
			pageSizeParameter: {fmt.Sprintf("%d", expectedPageSize)},
		})
	assert.Nil(t, err)
	assert.Equal(t, defaultOffset, offset)
	assert.Equal(t, expectedPageSize, pageSize)
}

func TestParsePagingParameters_WithOffsetAndPageSize(t *testing.T) {
	expectedOffset := 23
	expectedPageSize := 30
	offset, pageSize, err := parsePagingParameters(
		map[string][]string{
			offsetParameter:   {fmt.Sprintf("%d", expectedOffset)},
			pageSizeParameter: {fmt.Sprintf("%d", expectedPageSize)},
		})
	assert.Nil(t, err)
	assert.Equal(t, expectedOffset, offset)
	assert.Equal(t, expectedPageSize, pageSize)
}

func TestParsePagingParameters_WithNonNumericOffset(t *testing.T) {
	_, _, err := parsePagingParameters(
		map[string][]string{
			offsetParameter: {"NotANumber"},
		})
	assert.Error(t, err)
	assert.Equal(t, "offset must be numeric", err.Error())
}

func TestParsePagingParameters_WithNonNumericPageSize(t *testing.T) {
	_, _, err := parsePagingParameters(
		map[string][]string{
			pageSizeParameter: {"NotANumber"},
		})
	assert.Error(t, err)
	assert.Equal(t, "pageSize must be numeric", err.Error())
}

func TestParsePagingParameters_WithInvalidOffset(t *testing.T) {
	_, _, err := parsePagingParameters(
		map[string][]string{
			offsetParameter: {"-1"},
		})
	assert.Error(t, err)
	assert.Equal(t, "offset must be greater than or equal to 0", err.Error())
}

func TestParsePagingParameters_WithInvalidPageSize(t *testing.T) {
	_, _, err := parsePagingParameters(
		map[string][]string{
			pageSizeParameter: {"0"},
		})
	assert.Error(t, err)
	assert.Equal(t, "pageSize must be greater than 0", err.Error())
}

// endregion

// region pageResults

func TestPageResults(t *testing.T) {
	results := []*core.Metadata{
		{Id: uuid.New()},
		{Id: uuid.New()},
		{Id: uuid.New()},
	}
	offset := 0
	pageSize := 2
	request := httptest.NewRequest(http.MethodGet, "/metadata", nil)

	actual := pageResults(results, offset, pageSize, request)

	assert.Len(t, actual.Resources, pageSize)
	assert.Equal(t, results[0], actual.Resources[0])
	assert.Equal(t, results[1], actual.Resources[1])

	nextLink, err := url.Parse(actual.NextLink)
	assert.Nil(t, err)
	assert.Equal(t, "/metadata", nextLink.Path)
	assert.Equal(t, fmt.Sprintf("%d", offset+pageSize), nextLink.Query().Get(offsetParameter))
	assert.Equal(t, fmt.Sprintf("%d", pageSize), nextLink.Query().Get(pageSizeParameter))
}

func TestPageResults_EndOfPage(t *testing.T) {
	results := []*core.Metadata{
		{Id: uuid.New()},
		{Id: uuid.New()},
		{Id: uuid.New()},
	}
	offset := 2
	pageSize := 2
	request := httptest.NewRequest(http.MethodGet, "/metadata", nil)

	actual := pageResults(results, offset, pageSize, request)

	assert.Len(t, actual.Resources, 1)
	assert.Equal(t, results[2], actual.Resources[0])
	assert.Empty(t, actual.NextLink)
}

func TestPageResults_OffPage(t *testing.T) {
	results := []*core.Metadata{
		{Id: uuid.New()},
		{Id: uuid.New()},
		{Id: uuid.New()},
	}
	offset := 4
	pageSize := 2
	request := httptest.NewRequest(http.MethodGet, "/metadata", nil)

	actual := pageResults(results, offset, pageSize, request)

	assert.Empty(t, actual.Resources)
	assert.Empty(t, actual.NextLink)
}

func TestPageResults_BigPageSize(t *testing.T) {
	results := []*core.Metadata{
		{Id: uuid.New()},
		{Id: uuid.New()},
		{Id: uuid.New()},
	}
	offset := 0
	pageSize := 20
	request := httptest.NewRequest(http.MethodGet, "/metadata", nil)

	actual := pageResults(results, offset, pageSize, request)

	assert.Len(t, actual.Resources, len(results))
	assert.Equal(t, results[0], actual.Resources[0])
	assert.Equal(t, results[1], actual.Resources[1])
	assert.Equal(t, results[2], actual.Resources[2])
	assert.Empty(t, actual.NextLink)
}

// endregion
