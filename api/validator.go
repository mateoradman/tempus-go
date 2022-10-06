package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/mateoradman/tempus/util"
)

// validGender is a custom gender validator
var validGender validator.Func = func(fl validator.FieldLevel) bool {
	if gender, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedGender(gender)
	}
	return false
}
