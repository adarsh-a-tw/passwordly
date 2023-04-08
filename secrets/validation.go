package secrets

import (
	"github.com/gin-gonic/gin/binding"
	validator "github.com/go-playground/validator/v10"
)

func validateSecretType(fl validator.FieldLevel) bool {
	secretType, ok := fl.Field().Interface().(SecretType)
	if !ok {
		return false
	}

	return secretType.IsValid()
}

func RegisterValidations() {

	validators := []struct {
		name      string
		validator validator.Func
	}{
		{
			name:      "secret_type",
			validator: validateSecretType,
		},
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for _, val := range validators {
			v.RegisterValidation(val.name, val.validator)
		}
	}
}
