package core

import "github.com/go-playground/validator/v10"

func ValidateStruct(structToValidate interface{}) error {
	v := validator.New()
	return v.Struct(structToValidate)
}
