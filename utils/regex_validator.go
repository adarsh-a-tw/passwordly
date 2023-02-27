package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func RegexValidator(pattern string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		s, ok := fl.Field().Interface().(string)
		if ok {
			match, _ := regexp.MatchString(pattern, s)
			return match
		}
		return false
	}
}
