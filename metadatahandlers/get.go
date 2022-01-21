package metadatahandlers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
	"net/http"
)

// GET /metadata/{id}
func (m *MetadataHandlerManager) HandleMetadataGetWithId(
	w http.ResponseWriter,
	req *http.Request) {
	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error parsing ID: %v\n", err.Error())))
		return
	}

	result := m.Database.Metadatas[id]
	r, err := yaml.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error marshalling metadata: Error: %v", err.Error())))
		return
	}
	w.Header().Set("Content-Type", "application/x-yaml")
	w.WriteHeader(http.StatusOK)
	w.Write(r)
}

// GET /metadata
func (m *MetadataHandlerManager) HandleMetadataGet(
	w http.ResponseWriter,
	req *http.Request) {
	results, err := m.Filterer.FilterMetadata(req.URL.Query(), m.Database)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	r, err := yaml.Marshal(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error marshalling metadata: Error: %v", err.Error())))
		return
	}
	w.Header().Set("Content-Type", "application/x-yaml")
	w.WriteHeader(http.StatusOK)
	w.Write(r)
}
