package metadatahandlers

import (
	"APIServerExercise/core"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
	"net/http"
)

// PUT /metadata/{id}
func (m *MetadataHandlerManager) HandleMetadataPutWithId(
	w http.ResponseWriter,
	req *http.Request) {
	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error parsing ID: %v\n", err.Error())))
		return
	}
	m.handleMetadataPutInner(w, req, id)
}

// PUT /metadata
func (m *MetadataHandlerManager) HandleMetadataPut(
	w http.ResponseWriter,
	req *http.Request) {
	m.handleMetadataPutInner(w, req, uuid.UUID{})
}

func (m *MetadataHandlerManager) handleMetadataPutInner(
	w http.ResponseWriter,
	req *http.Request,
	id uuid.UUID) {
	// Decode request body
	var metadata core.Metadata
	decoder := yaml.NewDecoder(req.Body)
	if err := decoder.Decode(&metadata); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Failed to decode body: %v\n", err.Error())))
		return
	}

	// Validate request metadata
	if err := core.ValidateStruct(metadata); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Validation failed: %v\n", err.Error())))
		return
	}

	if id.String() != (uuid.UUID{}).String() {
		// Id was passed in from url, takes precedence
		metadata.Id = id
	} else if metadata.Id.String() == (uuid.UUID{}).String() {
		// No id was passed in request body, generate new id
		metadata.Id = uuid.New()
	} // Else use Id that was passed in from request body

	// Check if there is an existing metadata with the same Id
	// If so, remove old metadata Id from indexes
	if _, ok := m.Database.Metadatas[metadata.Id]; ok {
		m.Indexer.RemoveFromIndex(metadata.Id)
	}

	// Save metadata in database and add to index
	m.Database.Metadatas[metadata.Id] = &metadata
	m.Indexer.AddToIndex(&metadata, metadata.Id, "")

	responseByte, _ := yaml.Marshal(&metadata)
	w.Header().Set("Content-Type", "application/x-yaml")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseByte)
}
