package api

import (
	validator "github.com/go-playground/validator/v10"
	"github.com/mateoradman/tempus/internal/types"
)

// validGender is a custom gender validator
var validGender validator.Func = func(fl validator.FieldLevel) bool {
	if gender, ok := fl.Field().Interface().(string); ok {
		return types.IsSupportedGender(gender)
	}
	return false
}

// validLanguage is a custom language validator
var validLanguage validator.Func = func(fl validator.FieldLevel) bool {
	if lang, ok := fl.Field().Interface().(string); ok {
		return types.IsSupportedLanguage(lang)
	}
	return false
}
