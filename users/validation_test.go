package users_test

import (
	"encoding/json"
	"testing"

	"github.com/adarsh-a-tw/passwordly/users"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

var validate *validator.Validate

func init() {
	users.RegisterValidations()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v
	}
}

func TestUserValidator_CreateUserRequest(t *testing.T) {
	testCases := []struct {
		input                string
		expectedErrorMessage string
	}{
		{
			`{"username": "test_123", "password": "P@ssword123", "email": "test@mail.com"}`,
			"",
		},
		{
			`{"username": "t", "password": "P@ssword123", "email": "test@mail.com"}`,
			"Key: 'CreateUserRequest.Username' Error:Field validation for 'Username' failed on the 'username' tag",
		},
		{
			`{"username": "test_12345", "password": "P@ssword123", "email": "@mail.com"}`,
			"Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag",
		},
		{
			`{"username": "test_123", "password": "P@", "email": "test@mail.com"}`,
			"Key: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'password' tag",
		},
		{
			`{"username": "test_123", "password": "12345678910", "email": "test@mail.com"}`,
			"Key: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'password' tag",
		},
		{
			`{"username": "test_123", "password": "abcdefghijkl", "email": "test@mail.com"}`,
			"Key: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'password' tag",
		},
		{
			`{"username": "test_123", "password": "ABCDEfghij", "email": "test@mail.com"}`,
			"Key: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'password' tag",
		},
		{
			`{"username": "test_123", "password": "ABCDEfghij1234", "email": "test@mail.com"}`,
			"Key: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'password' tag",
		},
	}

	for _, testCase := range testCases {
		var cur users.CreateUserRequest
		if parsingErr := parseJSON(testCase.input, &cur); parsingErr != nil {
			t.FailNow()
		}

		if err := validateRequestObject(cur); err != nil {
			assert.Equal(t, testCase.expectedErrorMessage, err.Error())
		} else {
			assert.Equal(t, testCase.expectedErrorMessage, "")
		}
	}
}

func TestUserValidator_ChangePasswordRequest(t *testing.T) {
	testCases := []struct {
		input                string
		expectedErrorMessage string
	}{
		{
			`{"current_password": "P@ssword123", "new_password": "P@ssword1234"}`,
			"",
		},
		{
			`{"current_password": "P@ssword123", "new_password": "P@ssword"}`,
			"Key: 'ChangePasswordRequest.NewPassword' Error:Field validation for 'NewPassword' failed on the 'password' tag",
		},
	}

	for _, testCase := range testCases {
		var cur users.ChangePasswordRequest
		if parsingErr := parseJSON(testCase.input, &cur); parsingErr != nil {
			t.FailNow()
		}

		if err := validateRequestObject(cur); err != nil {
			assert.Equal(t, testCase.expectedErrorMessage, err.Error())
		} else {
			assert.Equal(t, testCase.expectedErrorMessage, "")
		}
	}
}

func parseJSON(jsonString string, obj any) error {
	return json.Unmarshal([]byte(jsonString), &obj)
}

func validateRequestObject(obj any) error {
	return validate.Struct(obj)
}
