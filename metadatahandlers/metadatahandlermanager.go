package metadatahandlers

import (
	"APIServerExercise/core"
	"APIServerExercise/search"
)

type MetadataHandlerManager struct {
	Database *core.Database
	Searcher *search.Searcher
}
