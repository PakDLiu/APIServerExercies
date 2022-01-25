package search

import (
	"APIServerExercise/core"
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"strings"
	"unicode"
)

//go:generate sh -c "test mocks/search/mock_search.go -nt $GOFILE && exit 0; mockgen -source=search.go -destination ../mock/search/mock_search.go"

type Indexer interface {
	AddToIndex(data interface{}, id uuid.UUID, prefix string)
	RemoveFromIndex(id uuid.UUID)
}

var _ Indexer = &Searcher{}

type Filterer interface {
	FilterMetadata(query map[string][]string, database *core.Database) ([]*core.Metadata, error)
}

var _ Filterer = &Searcher{}

type Searcher struct {
	// Map of field name -> Map of field value -> Set of metadata Ids
	Index             map[string]map[string]map[uuid.UUID]bool
	DisableIndexWords bool
}

// Adding data to index
// If data is a slice, will add children to index with data field name as prefix
func (s *Searcher) AddToIndex(data interface{}, id uuid.UUID, prefix string) {
	// For each field in data
	elements := reflect.ValueOf(data).Elem()
	for i := 0; i < elements.NumField(); i++ {
		fieldName := elements.Type().Field(i).Name
		// Add prefix to field name if provided
		if prefix != "" {
			fieldName = fmt.Sprintf("%s.%s", prefix, fieldName)
		}
		// Need to use lower case since the field names are lower case when returned to user
		fieldName = strings.ToLower(fieldName)
		fieldValueInterface := elements.Field(i).Interface()

		// Check to see if field is a slice
		rv := reflect.ValueOf(fieldValueInterface)
		if rv.Kind() == reflect.Slice {
			// If so, add children to index with field name as prefix
			for i := 0; i < rv.Len(); i++ {
				s.AddToIndex(rv.Index(i).Interface(), id, fieldName)
			}
			// skip adding the slice itself to the index
			continue
		}

		fieldValue := fmt.Sprintf("%v", fieldValueInterface)

		// Initialize index as needed
		if len(s.Index[fieldName]) == 0 {
			s.Index[fieldName] = map[string]map[uuid.UUID]bool{}
		}
		if len(s.Index[fieldName][fieldValue]) == 0 {
			s.Index[fieldName][fieldValue] = map[uuid.UUID]bool{}
		}

		// Add entire value to index
		s.Index[fieldName][fieldValue][id] = true

		if !s.DisableIndexWords {
			// Check to see if value has multiple words
			// If so, add each word to index as well
			fieldValueParts := strings.Fields(fieldValue)
			if len(fieldValueParts) <= 1 {
				continue
			}

			for _, part := range fieldValueParts {
				part, include := cleanWord(part)
				if !include {
					continue
				}

				// initialize map if needed
				if len(s.Index[fieldName][part]) == 0 {
					s.Index[fieldName][part] = map[uuid.UUID]bool{}
				}
				s.Index[fieldName][part][id] = true
			}
		}
	}
}

func (s *Searcher) RemoveFromIndex(id uuid.UUID) {
	for _, fieldValue := range s.Index {
		for key, uuids := range fieldValue {
			delete(uuids, id)
			// Check to see if there are any more ids for this field value
			// If not, remove the entry altogether
			if len(uuids) == 0 {
				delete(fieldValue, key)
			}
		}
	}
}

func (s *Searcher) FilterMetadata(query map[string][]string, database *core.Database) ([]*core.Metadata, error) {
	// Copy all metadatas
	results := make([]*core.Metadata, 0, len(database.Metadatas))
	for _, id := range database.Ordering {
		results = append(results, database.Metadatas[id])
	}

	// Filter by query parameters
	for queryKey, queryValues := range query {
		// If no more results are left, stop filtering
		if len(results) == 0 {
			break
		}

		fieldNameIndex, ok := s.Index[queryKey]
		if !ok {
			return nil, fmt.Errorf("no such field name %s", queryKey)
		}

		// Get the list of matched ids for the query
		// Only care about the first query value
		matchedIds := fieldNameIndex[queryValues[0]]
		newResult := make([]*core.Metadata, 0, len(matchedIds))

		// Craft new result list based on matching ids from index
		for index, metadata := range results {
			if _, ok := matchedIds[metadata.Id]; ok {
				newResult = append(newResult, results[index])
			}
		}
		results = newResult
	}

	return results, nil
}

// Trims the string
// Will removing characters from the beginning and the end of the string if it is not a letter or a number
// Returns the cleaned string and boolean whether to be indexed.
func cleanWord(s string) (string, bool) {
	s = strings.TrimFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	if len(s) == 0 {
		return "", false
	}
	return s, true
}
