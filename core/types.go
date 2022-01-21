package core

import (
	"APIServerExercise/util"
	"github.com/google/uuid"
)

type Database struct {
	Metadatas map[uuid.UUID]*Metadata `yaml:"metadatas"`
}

type Metadata struct {
	Id          uuid.UUID     `yaml:"id"`
	Title       string        `yaml:"title" validate:"required"`
	Version     string        `yaml:"version" validate:"required"`
	Maintainers []*Maintainer `yaml:"maintainers" validate:"required,gt=0,dive"`
	Company     string        `yaml:"company" validate:"required"`
	Website     util.Yamlurl  `yaml:"website" validate:"required"`
	Source      util.Yamlurl  `yaml:"source" validate:"required"`
	License     string        `yaml:"license" validate:"required"`
	Description string        `yaml:"description" validate:"required"`
}

type Maintainer struct {
	Name  string `yaml:"name" validate:"required"`
	Email string `yaml:"email" validate:"required,email"`
}
