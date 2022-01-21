package metadatahandlers

import (
	"APIServerExercise/core"
	mock_search "APIServerExercise/mock/search"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetadataHandlerManager_HandleMetadataDeleteWithId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()

	request := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/metadata/%s", id.String()),
		nil)
	request = mux.SetURLVars(request, map[string]string{
		"id": id.String(),
	})
	responseRecorder := httptest.NewRecorder()

	mockIndexer := mock_search.NewMockIndexer(ctrl)
	mockIndexer.EXPECT().RemoveFromIndex(id).Times(1)

	manager := MetadataHandlerManager{
		Database: &core.Database{Metadatas: map[uuid.UUID]*core.Metadata{
			id: {},
		}},
		Indexer: mockIndexer,
	}
	manager.HandleMetadataDeleteWithId(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Len(t, manager.Database.Metadatas, 0)
}

func TestMetadataHandlerManager_HandleMetadataDeleteWithId_WithInvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := "invalidId"

	request := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/metadata/%s", id),
		nil)
	request = mux.SetURLVars(request, map[string]string{
		"id": id,
	})
	responseRecorder := httptest.NewRecorder()

	manager := MetadataHandlerManager{}
	manager.HandleMetadataDeleteWithId(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "Error parsing ID: invalid UUID length:")
}

func TestMetadataHandlerManager_HandleMetadataDeleteWithId_WithNonExistentId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()

	request := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/metadata/%s", id.String()),
		nil)
	request = mux.SetURLVars(request, map[string]string{
		"id": id.String(),
	})
	responseRecorder := httptest.NewRecorder()

	manager := MetadataHandlerManager{
		Database: &core.Database{Metadatas: map[uuid.UUID]*core.Metadata{
			uuid.New(): {},
		}},
	}
	manager.HandleMetadataDeleteWithId(responseRecorder, request)

	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
	assert.Len(t, manager.Database.Metadatas, 1)
}
