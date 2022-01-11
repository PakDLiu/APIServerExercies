package metadatahandlers

import (
	"APIServerExercies/core"
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

	query := map[string][]string{
		"id": {
			id.String(),
		},
	}
	results, err := m.Searcher.FilterMetadata(query, m.Database)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	m.returnResults(w, results)
}

// GET /metadata
func (m *MetadataHandlerManager) HandleMetadataGet(
	w http.ResponseWriter,
	req *http.Request) {
	results, err := m.Searcher.FilterMetadata(req.URL.Query(), m.Database)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	m.returnResults(w, results)
}

func (m *MetadataHandlerManager) returnResults(w http.ResponseWriter, results []core.Metadata) {
	r, err := yaml.Marshal(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error unmarshalling metadata: Error: %v", err.Error())))
		return
	}
	w.Header().Set("Content-Type", "application/x-yaml")
	w.WriteHeader(http.StatusOK)
	w.Write(r)
}
