package metadatahandlers

import (
	"APIServerExercise/core"
	"APIServerExercise/search"
)

type MetadataHandlerManager struct {
	Database *core.Database
	Indexer  search.Indexer
	Filterer search.Filterer
}
