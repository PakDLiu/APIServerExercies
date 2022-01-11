package metadatahandlers

import (
	"APIServerExercies/core"
	"APIServerExercies/search"
)

type MetadataHandlerManager struct {
	Database *core.Database
	Searcher *search.Searcher
}
