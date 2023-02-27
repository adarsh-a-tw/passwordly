package users

import (
	"regexp"

	"github.com/adarsh-a-tw/passwordly/utils"
	"github.com/gin-gonic/gin/binding"
	validator "github.com/go-playground/validator/v10"
)

func validatePassword(fl validator.FieldLevel) bool {
	password, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if len(password) < 8 {
		return false
	}
	done, err := regexp.MatchString("([a-z])+", password)
	if err != nil || !done {
		return false
	}

	done, err = regexp.MatchString("([A-Z])+", password)
	if err != nil || !done {
		return false
	}
	done, err = regexp.MatchString("([0-9])+", password)
	if err != nil || !done {
		return false
	}

	done, err = regexp.MatchString("([!@#$%^&*.?-])+", password)
	if err != nil || !done {
		return false
	}

	return true
}

func RegisterValidations() {

	usernamePattern := "^[a-zA-Z0-9_-]{5,20}$"

	validators := []struct {
		name      string
		validator validator.Func
	}{
		{
			name:      "username",
			validator: utils.RegexValidator(usernamePattern),
		},
		{
			name:      "password",
			validator: validatePassword,
		},
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for _, val := range validators {
			v.RegisterValidation(val.name, val.validator)
		}
	}
}
