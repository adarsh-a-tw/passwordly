package users_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/adarsh-a-tw/passwordly/users"
	user_mocks "github.com/adarsh-a-tw/passwordly/users/mocks"
	utils_mocks "github.com/adarsh-a-tw/passwordly/utils/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TODO: Need to check if there is a better way to do this
func init() {
	user_mocks.RegisterMockValidations()
}

func TestUserHandler_Login_ShouldLoginUserSuccessfully(t *testing.T) {
	mockAPIToken := "MOCK_API_TOKEN"
	expectedResponse := users.LoginUserSuccessResponse{Token: mockAPIToken}

	lur := users.LoginUserRequest{
		Username: "mock_username",
		Password: "mockPassword@123",
	}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/users/login", "POST", lur)

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}

	repo.On("Find", "mock_username", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = mu.Password
	})
	ap.On("GenerateToken", "mock_id").Return(mockAPIToken, nil).Once()

	uh := users.UserHandler{
		Repo:         repo,
		AuthProvider: ap,
	}

	uh.Login(ctx)

	var actualResponse users.LoginUserSuccessResponse
	decodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "Find", "mock_username", mock.AnythingOfType("*users.User"))
	ap.AssertCalled(t, "GenerateToken", "mock_id")

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_Login_ShouldThrowErrorForInternalServerError(t *testing.T) {
	mockAPIToken := "MOCK_API_TOKEN"
	expectedResponse := common.ErrorResponse{Message: "Something went wrong. Try again."}

	lur := users.LoginUserRequest{
		Username: "mock_username",
		Password: "mockPassword@123",
	}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/users/login", "POST", lur)

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}

	repo.On("Find", "mock_username", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = mu.Password
	})
	ap.On("GenerateToken", "mock_id").Return(mockAPIToken, nil).Return(
		"", errors.New("Something went wrong. Try again."),
	)

	uh := users.UserHandler{
		Repo:         repo,
		AuthProvider: ap,
	}

	uh.Login(ctx)

	var actualResponse common.ErrorResponse
	decodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "Find", "mock_username", mock.AnythingOfType("*users.User"))
	ap.AssertCalled(t, "GenerateToken", "mock_id")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_Login_ShouldThrowErrorForInvalidPassword(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Invalid Credentials"}

	lur := users.LoginUserRequest{
		Username: "mock_username",
		Password: "mockPassword@1234",
	}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/users/login", "POST", lur)

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}

	repo.On("Find", "mock_username", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = mu.Password
	})

	uh := users.UserHandler{
		Repo:         repo,
		AuthProvider: ap,
	}

	uh.Login(ctx)

	var actualResponse common.ErrorResponse
	decodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "Find", "mock_username", mock.AnythingOfType("*users.User"))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_Login_ShouldThrowErrorForInvalidUsername(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Invalid Credentials"}

	lur := users.LoginUserRequest{
		Username: "mock_username",
		Password: "mockPassword@123",
	}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/users/login", "POST", lur)

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}

	repo.On("Find", "mock_username", mock.AnythingOfType("*users.User")).Return(
		errors.New("Cannot find user with given username"),
	)

	uh := users.UserHandler{
		Repo:         repo,
		AuthProvider: ap,
	}

	uh.Login(ctx)

	var actualResponse common.ErrorResponse
	decodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "Find", "mock_username", mock.AnythingOfType("*users.User"))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_Login_ShouldThrowErrorForInvalidRequestBody(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Invalid request body"}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/users/login", "POST", nil)

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}

	uh := users.UserHandler{
		Repo:         repo,
		AuthProvider: ap,
	}

	uh.Login(ctx)

	var actualResponse common.ErrorResponse
	decodeJSONResponse(t, rec, &actualResponse)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_FetchUser_ShouldFetchUserDetailsSuccessfully(t *testing.T) {
	expectedResponse := users.UserResponse{Id: "mock_id", Username: "mock_username", Email: "test@email.com"}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/users/me", "GET", nil)

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}

	repo.On("FindById", "mock_id", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = mu.Password
	})

	// Mocking TokenAuthMiddleware
	ctx.Set("user_id", "mock_id")

	uh := users.UserHandler{
		Repo:         repo,
		AuthProvider: ap,
	}

	uh.FetchUser(ctx)

	var actualResponse users.UserResponse
	decodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "FindById", "mock_id", mock.AnythingOfType("*users.User"))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_Login_ShouldThrowErrorForInvalidId(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Something went wrong. Try again."}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/users/me", "GET", nil)

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}

	repo.On("FindById", "invalid_id", mock.AnythingOfType("*users.User")).Return(
		errors.New("Cannot find user with given username"),
	)

	// Mocking TokenAuthMiddleware
	ctx.Set("user_id", "invalid_id")

	uh := users.UserHandler{
		Repo:         repo,
		AuthProvider: ap,
	}

	uh.FetchUser(ctx)

	var actualResponse common.ErrorResponse
	decodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "FindById", "invalid_id", mock.AnythingOfType("*users.User"))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_ChangePassword_ShouldChangePasswordSuccessfully(t *testing.T) {
	cpr := users.ChangePasswordRequest{
		CurrentPassword: "mockPassword@123",
		NewPassword:     "mockPassword@1234",
	}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/users/me/password", "PATCH", cpr)

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}

	repo.On("FindById", "mock_id", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = mu.Password
	})

	repo.On("Update", mock.AnythingOfType("*users.User")).Return(nil)

	// Mocking TokenAuthMiddleware
	ctx.Set("user_id", "mock_id")

	uh := users.UserHandler{
		Repo:         repo,
		AuthProvider: ap,
	}

	uh.ChangePassword(ctx)

	repo.AssertCalled(t, "FindById", "mock_id", mock.AnythingOfType("*users.User"))
	repo.AssertCalled(t, "Update", mock.AnythingOfType("*users.User"))

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUserHandler_ChangePassword_ShouldThrowErrorIfCurrentPasswordInvalid(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Current password does not match"}
	cpr := users.ChangePasswordRequest{
		CurrentPassword: "wrongPassword",
		NewPassword:     "mockPassword@1234",
	}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/users/me/password", "PATCH", cpr)

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}

	repo.On("FindById", "mock_id", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = mu.Password
	})

	// Mocking TokenAuthMiddleware
	ctx.Set("user_id", "mock_id")

	uh := users.UserHandler{
		Repo:         repo,
		AuthProvider: ap,
	}

	uh.ChangePassword(ctx)

	var actualResponse common.ErrorResponse
	decodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "FindById", "mock_id", mock.AnythingOfType("*users.User"))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_ChangePassword_ShouldThrowErrorIfCurrentPasswordAndNewPasswordSame(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "New password and current password should not be the same"}
	cpr := users.ChangePasswordRequest{
		CurrentPassword: "mockPassword@123",
		NewPassword:     "mockPassword@123",
	}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/users/me/password", "PATCH", cpr)

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}

	repo.On("FindById", "mock_id", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = mu.Password
	})

	// Mocking TokenAuthMiddleware
	ctx.Set("user_id", "mock_id")

	uh := users.UserHandler{
		Repo:         repo,
		AuthProvider: ap,
	}

	uh.ChangePassword(ctx)

	var actualResponse common.ErrorResponse
	decodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "FindById", "mock_id", mock.AnythingOfType("*users.User"))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_ChangePassword_ShouldThrowErrorForInvalidBody(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Invalid Request body"}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/users/me/password", "PATCH", nil)

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}

	uh := users.UserHandler{
		Repo:         repo,
		AuthProvider: ap,
	}

	uh.ChangePassword(ctx)

	var actualResponse common.ErrorResponse
	decodeJSONResponse(t, rec, &actualResponse)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_ChangePassword_ShouldThrowErrorForInvalidUserId(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Something went wrong. Try again."}
	cpr := users.ChangePasswordRequest{
		CurrentPassword: "wrongPassword",
		NewPassword:     "mockPassword@1234",
	}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/users/me/password", "PATCH", cpr)

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}

	repo.On("FindById", "invalid_id", mock.AnythingOfType("*users.User")).Return(
		errors.New("Cannot find user with given username"),
	)

	// Mocking TokenAuthMiddleware
	ctx.Set("user_id", "invalid_id")

	uh := users.UserHandler{
		Repo:         repo,
		AuthProvider: ap,
	}

	uh.ChangePassword(ctx)

	var actualResponse common.ErrorResponse
	decodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "FindById", "invalid_id", mock.AnythingOfType("*users.User"))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_ChangePassword_ShouldThrowErrorForUpdateFailure(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Something went wrong. Try again."}
	cpr := users.ChangePasswordRequest{
		CurrentPassword: "mockPassword@123",
		NewPassword:     "mockPassword@1234",
	}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/users/me/password", "PATCH", cpr)

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}

	repo.On("FindById", "mock_id", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = mu.Password
	})

	repo.On("Update", mock.AnythingOfType("*users.User")).Return(errors.New("Mock Error"))

	// Mocking TokenAuthMiddleware
	ctx.Set("user_id", "mock_id")

	uh := users.UserHandler{
		Repo:         repo,
		AuthProvider: ap,
	}

	uh.ChangePassword(ctx)

	var actualResponse common.ErrorResponse
	decodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "FindById", "mock_id", mock.AnythingOfType("*users.User"))
	repo.AssertCalled(t, "Update", mock.AnythingOfType("*users.User"))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func mockUser() *users.User {
	return &users.User{
		Id:        "mock_id",
		Username:  "mock_username",
		Password:  "mockPassword@123",
		Email:     "test@email.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func prepareContextAndResponseRecorder(t *testing.T, url string, method string, reqBody any) (ctx *gin.Context, rec *httptest.ResponseRecorder) {
	rec = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(rec)

	var req *http.Request

	if reqBody != nil {
		jsonBytes, err := json.Marshal(reqBody)
		assert.NoError(t, err) // json.Marshal error
		buffer := bytes.NewBuffer(jsonBytes)
		req = httptest.NewRequest(method, url, buffer)
	} else {
		req = httptest.NewRequest(method, url, nil)
	}
	ctx.Request = req
	return
}

func decodeJSONResponse(t *testing.T, rec *httptest.ResponseRecorder, obj any) {
	if err := json.NewDecoder(rec.Body).Decode(obj); err != nil {
		t.Failed()
	}
}
