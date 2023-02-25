package user_mocks

import (
	"github.com/gin-gonic/gin/binding"
	validator "github.com/go-playground/validator/v10"
)

func alwaysValid(validator.FieldLevel) bool {
	return true
}

func RegisterMockValidations() {

	validators := []struct {
		name      string
		validator validator.Func
	}{
		{
			name:      "username",
			validator: alwaysValid,
		},
		{
			name:      "password",
			validator: alwaysValid,
		},
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for _, val := range validators {
			v.RegisterValidation(val.name, val.validator)
		}
	}
}
