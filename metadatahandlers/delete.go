package metadatahandlers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

// DELETE /metadata/{id}
func (m *MetadataHandlerManager) HandleMetadataDeleteWithId(
	w http.ResponseWriter,
	req *http.Request) {
	vars := mux.Vars(req)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error parsing ID: %v\n", err.Error())))
		return
	}

	if _, ok := m.Database.Metadatas[id]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	m.Indexer.RemoveFromIndex(id)
	delete(m.Database.Metadatas, id)
	for index, existingId := range m.Database.Ordering {
		if id == existingId {
			// remove id from Ordering
			m.Database.Ordering = append(m.Database.Ordering[:index], m.Database.Ordering[index+1:]...)
			break
		}
	}

	w.WriteHeader(http.StatusOK)
}
