package metadatahandlers

import (
	"APIServerExercise/core"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
	"net/http"
	"strconv"
)

const (
	offsetParameter   = "offset"
	pageSizeParameter = "pageSize"
	defaultOffset     = 0
	defaultPageSize   = 10
	// Hardcoded http as we are using http.ListenAndServe and NOT http.ListenAndServeTLS
	// Will need to update this if we use https
	scheme = "http://"
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
	var err error

	var query map[string][]string = req.URL.Query()

	offset, pageSize, err := parsePagingParameters(query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Need to remove the paging parameters as they are not going to be in the index
	delete(query, offsetParameter)
	delete(query, pageSizeParameter)

	results, err := m.Filterer.FilterMetadata(query, m.Database)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	page := pageResults(results, offset, pageSize, req)

	p, err := yaml.Marshal(page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error marshalling metadata: Error: %v", err.Error())))
		return
	}
	w.Header().Set("Content-Type", "application/x-yaml")
	w.WriteHeader(http.StatusOK)
	w.Write(p)
}

func pageResults(
	results []*core.Metadata,
	offset int,
	pageSize int,
	req *http.Request) *core.ResultPage {
	begin := offset
	end := offset + pageSize
	addNextLink := true

	if end >= len(results) {
		end = len(results)
		addNextLink = false
	}
	if begin > len(results) {
		begin = len(results)
	}

	var nextLink string

	if addNextLink {
		originalUrl := req.URL
		query := originalUrl.Query()
		query.Set(offsetParameter, fmt.Sprintf("%d", end))
		query.Set(pageSizeParameter, fmt.Sprintf("%d", pageSize))
		originalUrl.RawQuery = query.Encode()
		nextLink = fmt.Sprintf("%s%s%s", scheme, req.Host, originalUrl.String())
	}

	return &core.ResultPage{
		Resources: results[begin:end],
		NextLink:  nextLink,
	}
}

func parsePagingParameters(query map[string][]string) (int, int, error) {
	var err error
	offset := defaultOffset
	pageSize := defaultPageSize

	if o, ok := query[offsetParameter]; ok {
		offset, err = strconv.Atoi(o[0])
		if err != nil {
			return 0, 0, fmt.Errorf("offset must be numeric")
		}
	}
	if o, ok := query[pageSizeParameter]; ok {
		pageSize, err = strconv.Atoi(o[0])
		if err != nil {
			return 0, 0, fmt.Errorf("pageSize must be numeric")
		}
	}

	if offset < 0 {
		return 0, 0, fmt.Errorf("offset must be greater than or equal to 0")
	}
	if pageSize < 1 {
		return 0, 0, fmt.Errorf("pageSize must be greater than 0")
	}

	return offset, pageSize, nil
}
