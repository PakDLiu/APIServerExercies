package main

import (
	"APIServerExercise/core"
	"APIServerExercise/metadatahandlers"
	"APIServerExercise/search"
	"flag"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var database *core.Database
var searcher *search.Searcher

func handleMetadata(w http.ResponseWriter, req *http.Request) {
	manager := metadatahandlers.MetadataHandlerManager{
		Database: database,
		Indexer:  searcher,
		Filterer: searcher,
	}

	switch req.Method {
	case http.MethodGet:
		manager.HandleMetadataGet(w, req)
	case http.MethodPut:
		manager.HandleMetadataPut(w, req)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleMetadataWithId(w http.ResponseWriter, req *http.Request) {
	manager := metadatahandlers.MetadataHandlerManager{
		Database: database,
		Indexer:  searcher,
		Filterer: searcher,
	}

	switch req.Method {
	case http.MethodGet:
		manager.HandleMetadataGetWithId(w, req)
	case http.MethodPut:
		manager.HandleMetadataPutWithId(w, req)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	disableIndexWordsFlag := flag.Bool(
		"disableIndexWords",
		false,
		"Disable indexing part of values. IE: Do not index each word in description field")
	flag.Parse()

	database = &core.Database{Metadatas: map[uuid.UUID]core.Metadata{}}
	searcher = &search.Searcher{
		Index:             map[string]map[string]map[uuid.UUID]bool{},
		DisableIndexWords: *disableIndexWordsFlag,
	}

	r := mux.NewRouter()
	r.HandleFunc("/metadata", handleMetadata)
	r.HandleFunc("/metadata/{id}", handleMetadataWithId)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
