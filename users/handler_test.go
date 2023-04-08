package users_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/adarsh-a-tw/passwordly/users"
	user_mocks "github.com/adarsh-a-tw/passwordly/users/mocks"
	utils_mocks "github.com/adarsh-a-tw/passwordly/utils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TODO: Need to check if there is a better way to do this
func init() {
	user_mocks.RegisterMockValidations()
}

func TestUserHandler_Create_ShouldCreateUserSuccessfully(t *testing.T) {
	cur := users.CreateUserRequest{
		Username: "mock_username",
		Password: "P@ssword123",
		Email:    "test@email.com",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users", "POST", cur)

	repo := &user_mocks.UserRepository{}
	hasher := &utils_mocks.PasswordHasher{}

	repo.On("UsernameAlreadyExists", "mock_username").Return(false, nil)
	repo.On("EmailAlreadyExists", "test@email.com").Return(false, nil)
	repo.On("Create", mock.AnythingOfType("*users.User")).Return(nil)

	hasher.On("HashPassword", "P@ssword123").Return("HashedPassword")

	uh := users.UserHandler{
		Repo:           repo,
		PasswordHasher: hasher,
	}

	uh.Create(ctx)

	var actualResponse users.UserResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "UsernameAlreadyExists", "mock_username")
	repo.AssertCalled(t, "EmailAlreadyExists", "test@email.com")
	repo.AssertCalled(t, "Create", mock.AnythingOfType("*users.User"))

	hasher.AssertCalled(t, "HashPassword", "P@ssword123")

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, cur.Username, actualResponse.Username)
	assert.Equal(t, cur.Email, actualResponse.Email)
}

func TestUserHandler_Create_ShouldNotCreateUserWithAlreadyExistingUsername(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Username already exists. Try another."}
	cur := users.CreateUserRequest{
		Username: "mock_username",
		Password: "P@ssword123",
		Email:    "test@email.com",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users", "POST", cur)

	repo := &user_mocks.UserRepository{}

	repo.On("UsernameAlreadyExists", "mock_username").Return(true, nil)

	uh := users.UserHandler{
		Repo: repo,
	}

	uh.Create(ctx)

	var actualResponse common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "UsernameAlreadyExists", "mock_username")

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_Create_ShouldNotCreateUserWithAlreadyExistingEmail(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Email already exists. Try another."}
	cur := users.CreateUserRequest{
		Username: "mock_username",
		Password: "P@ssword123",
		Email:    "test@email.com",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users", "POST", cur)

	repo := &user_mocks.UserRepository{}

	repo.On("UsernameAlreadyExists", "mock_username").Return(false, nil)
	repo.On("EmailAlreadyExists", "test@email.com").Return(true, nil)

	uh := users.UserHandler{
		Repo: repo,
	}

	uh.Create(ctx)

	var actualResponse common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "UsernameAlreadyExists", "mock_username")
	repo.AssertCalled(t, "EmailAlreadyExists", "test@email.com")

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_Create_ShouldNotCreateUserWithInvalidRequestBody(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Invalid Request body"}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users", "POST", nil)

	repo := &user_mocks.UserRepository{}

	uh := users.UserHandler{
		Repo: repo,
	}

	uh.Create(ctx)

	var actualResponse common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_Create_ShouldThrowInternalServerErrorIfUsernameAlreadyExistsMethodFails(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Something went wrong. Try again."}
	cur := users.CreateUserRequest{
		Username: "mock_username",
		Password: "P@ssword123",
		Email:    "test@email.com",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users", "POST", cur)

	repo := &user_mocks.UserRepository{}

	repo.On("UsernameAlreadyExists", "mock_username").Return(false, errors.New("MOCK_ERROR"))

	uh := users.UserHandler{
		Repo: repo,
	}

	uh.Create(ctx)

	var actualResponse common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "UsernameAlreadyExists", "mock_username")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_Create_ShouldThrowInternalServerErrorIfEmailAlreadyExistsMethodFails(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Something went wrong. Try again."}
	cur := users.CreateUserRequest{
		Username: "mock_username",
		Password: "P@ssword123",
		Email:    "test@email.com",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users", "POST", cur)

	repo := &user_mocks.UserRepository{}

	repo.On("UsernameAlreadyExists", "mock_username").Return(false, nil)
	repo.On("EmailAlreadyExists", "test@email.com").Return(false, errors.New("MOCK_ERROR"))

	uh := users.UserHandler{
		Repo: repo,
	}

	uh.Create(ctx)

	var actualResponse common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "UsernameAlreadyExists", "mock_username")
	repo.AssertCalled(t, "EmailAlreadyExists", "test@email.com")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_Create_ShouldThrowInternalServerErrorIfCreateMethodFails(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Something went wrong. Try again."}
	cur := users.CreateUserRequest{
		Username: "mock_username",
		Password: "P@ssword123",
		Email:    "test@email.com",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users", "POST", cur)

	repo := &user_mocks.UserRepository{}
	hasher := &utils_mocks.PasswordHasher{}

	repo.On("UsernameAlreadyExists", "mock_username").Return(false, nil)
	repo.On("EmailAlreadyExists", "test@email.com").Return(false, nil)
	repo.On("Create", mock.AnythingOfType("*users.User")).Return(errors.New("MOCK_ERROR"))

	hasher.On("HashPassword", "P@ssword123").Return("HashedPassword")

	uh := users.UserHandler{
		Repo:           repo,
		PasswordHasher: hasher,
	}

	uh.Create(ctx)

	var actualResponse common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "UsernameAlreadyExists", "mock_username")
	repo.AssertCalled(t, "EmailAlreadyExists", "test@email.com")
	repo.AssertCalled(t, "Create", mock.AnythingOfType("*users.User"))

	hasher.AssertCalled(t, "HashPassword", "P@ssword123")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_Login_ShouldLoginUserSuccessfully(t *testing.T) {
	mockAPIToken := "MOCK_API_TOKEN"
	expectedResponse := users.LoginUserSuccessResponse{Token: mockAPIToken}

	lur := users.LoginUserRequest{
		Username: "mock_username",
		Password: "mockPassword@123",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users/login", "POST", lur)

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}
	hasher := &utils_mocks.PasswordHasher{}

	repo.On("Find", "mock_username", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = "HashedPassword"
	})
	ap.On("GenerateToken", "mock_id").Return(mockAPIToken, nil).Once()
	hasher.On("ComparePassword", "mockPassword@123", "HashedPassword").Return(true)

	uh := users.UserHandler{
		Repo:           repo,
		AuthProvider:   ap,
		PasswordHasher: hasher,
	}

	uh.Login(ctx)

	var actualResponse users.LoginUserSuccessResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "Find", "mock_username", mock.AnythingOfType("*users.User"))
	ap.AssertCalled(t, "GenerateToken", "mock_id")
	hasher.AssertCalled(t, "ComparePassword", "mockPassword@123", "HashedPassword")

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

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users/login", "POST", lur)

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}
	hasher := &utils_mocks.PasswordHasher{}

	repo.On("Find", "mock_username", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = "HashedPassword"
	})
	ap.On("GenerateToken", "mock_id").Return(mockAPIToken, nil).Return(
		"", errors.New("Something went wrong. Try again."),
	)
	hasher.On("ComparePassword", "mockPassword@123", "HashedPassword").Return(true)

	uh := users.UserHandler{
		Repo:           repo,
		AuthProvider:   ap,
		PasswordHasher: hasher,
	}

	uh.Login(ctx)

	var actualResponse common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "Find", "mock_username", mock.AnythingOfType("*users.User"))
	ap.AssertCalled(t, "GenerateToken", "mock_id")
	hasher.AssertCalled(t, "ComparePassword", "mockPassword@123", "HashedPassword")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_Login_ShouldThrowErrorForInvalidPassword(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Invalid Credentials"}

	lur := users.LoginUserRequest{
		Username: "mock_username",
		Password: "mockPassword@123",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users/login", "POST", lur)

	repo := &user_mocks.UserRepository{}
	hasher := &utils_mocks.PasswordHasher{}

	repo.On("Find", "mock_username", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = "HashedPassword"
	})
	hasher.On("ComparePassword", "mockPassword@123", "HashedPassword").Return(false)

	uh := users.UserHandler{
		Repo:           repo,
		PasswordHasher: hasher,
	}

	uh.Login(ctx)

	var actualResponse common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "Find", "mock_username", mock.AnythingOfType("*users.User"))
	hasher.AssertCalled(t, "ComparePassword", "mockPassword@123", "HashedPassword")

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_Login_ShouldThrowErrorForInvalidUsername(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Invalid Credentials"}

	lur := users.LoginUserRequest{
		Username: "mock_username",
		Password: "mockPassword@123",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users/login", "POST", lur)

	repo := &user_mocks.UserRepository{}

	repo.On("Find", "mock_username", mock.AnythingOfType("*users.User")).Return(
		errors.New("Cannot find user with given username"),
	)

	uh := users.UserHandler{
		Repo: repo,
	}

	uh.Login(ctx)

	var actualResponse common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "Find", "mock_username", mock.AnythingOfType("*users.User"))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_Login_ShouldThrowErrorForInvalidRequestBody(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Invalid request body"}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users/login", "POST", nil)

	uh := users.UserHandler{}

	uh.Login(ctx)

	var actualResponse common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_FetchUser_ShouldFetchUserDetailsSuccessfully(t *testing.T) {
	expectedResponse := users.UserResponse{Id: "mock_id", Username: "mock_username", Email: "test@email.com"}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users/me", "GET", nil)

	repo := &user_mocks.UserRepository{}

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
		Repo: repo,
	}

	uh.FetchUser(ctx)

	var actualResponse users.UserResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "FindById", "mock_id", mock.AnythingOfType("*users.User"))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_Login_ShouldThrowErrorForInvalidId(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Something went wrong. Try again."}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users/me", "GET", nil)

	repo := &user_mocks.UserRepository{}

	repo.On("FindById", "invalid_id", mock.AnythingOfType("*users.User")).Return(
		errors.New("Cannot find user with given username"),
	)

	// Mocking TokenAuthMiddleware
	ctx.Set("user_id", "invalid_id")

	uh := users.UserHandler{
		Repo: repo,
	}

	uh.FetchUser(ctx)

	var actualResponse common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "FindById", "invalid_id", mock.AnythingOfType("*users.User"))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_ChangePassword_ShouldChangePasswordSuccessfully(t *testing.T) {
	cpr := users.ChangePasswordRequest{
		CurrentPassword: "mockPassword@123",
		NewPassword:     "mockPassword@1234",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users/me/password", "PATCH", cpr)

	repo := &user_mocks.UserRepository{}
	hasher := &utils_mocks.PasswordHasher{}

	repo.On("FindById", "mock_id", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = "HashedPassword"
	})

	repo.On("Update", mock.AnythingOfType("*users.User")).Return(nil)
	hasher.On("ComparePassword", "mockPassword@123", "HashedPassword").Return(true)
	hasher.On("HashPassword", "mockPassword@1234").Return("HashedPassword2")

	// Mocking TokenAuthMiddleware
	ctx.Set("user_id", "mock_id")

	uh := users.UserHandler{
		Repo:           repo,
		PasswordHasher: hasher,
	}

	uh.ChangePassword(ctx)

	repo.AssertCalled(t, "FindById", "mock_id", mock.AnythingOfType("*users.User"))
	repo.AssertCalled(t, "Update", mock.AnythingOfType("*users.User"))

	hasher.AssertCalled(t, "ComparePassword", "mockPassword@123", "HashedPassword")
	hasher.AssertCalled(t, "HashPassword", "mockPassword@1234")

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUserHandler_ChangePassword_ShouldThrowErrorIfCurrentPasswordInvalid(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Current password does not match"}
	cpr := users.ChangePasswordRequest{
		CurrentPassword: "wrongPassword",
		NewPassword:     "mockPassword@1234",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users/me/password", "PATCH", cpr)

	repo := &user_mocks.UserRepository{}
	hasher := &utils_mocks.PasswordHasher{}

	repo.On("FindById", "mock_id", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = "HashedPassword"
	})
	hasher.On("ComparePassword", "wrongPassword", "HashedPassword").Return(false)

	// Mocking TokenAuthMiddleware
	ctx.Set("user_id", "mock_id")

	uh := users.UserHandler{
		Repo:           repo,
		PasswordHasher: hasher,
	}

	uh.ChangePassword(ctx)

	var actualResponse common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "FindById", "mock_id", mock.AnythingOfType("*users.User"))

	hasher.AssertCalled(t, "ComparePassword", "wrongPassword", "HashedPassword")

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_ChangePassword_ShouldThrowErrorIfCurrentPasswordAndNewPasswordSame(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "New password and current password should not be the same"}
	cpr := users.ChangePasswordRequest{
		CurrentPassword: "mockPassword@123",
		NewPassword:     "mockPassword@123",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users/me/password", "PATCH", cpr)

	repo := &user_mocks.UserRepository{}
	hasher := &utils_mocks.PasswordHasher{}

	repo.On("FindById", "mock_id", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = "HashedPassword"
	})
	hasher.On("ComparePassword", "mockPassword@123", "HashedPassword").Return(true)

	// Mocking TokenAuthMiddleware
	ctx.Set("user_id", "mock_id")

	uh := users.UserHandler{
		Repo:           repo,
		PasswordHasher: hasher,
	}

	uh.ChangePassword(ctx)

	var actualResponse common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "FindById", "mock_id", mock.AnythingOfType("*users.User"))

	hasher.AssertCalled(t, "ComparePassword", "mockPassword@123", "HashedPassword")

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_ChangePassword_ShouldThrowErrorForInvalidBody(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Invalid Request body"}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users/me/password", "PATCH", nil)

	uh := users.UserHandler{}

	uh.ChangePassword(ctx)

	var actualResponse common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_ChangePassword_ShouldThrowErrorForInvalidUserId(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Something went wrong. Try again."}
	cpr := users.ChangePasswordRequest{
		CurrentPassword: "wrongPassword",
		NewPassword:     "mockPassword@1234",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users/me/password", "PATCH", cpr)

	repo := &user_mocks.UserRepository{}

	repo.On("FindById", "invalid_id", mock.AnythingOfType("*users.User")).Return(
		errors.New("Cannot find user with given username"),
	)

	// Mocking TokenAuthMiddleware
	ctx.Set("user_id", "invalid_id")

	uh := users.UserHandler{
		Repo: repo,
	}

	uh.ChangePassword(ctx)

	var actualResponse common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

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

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users/me/password", "PATCH", cpr)

	repo := &user_mocks.UserRepository{}
	hasher := &utils_mocks.PasswordHasher{}

	repo.On("FindById", "mock_id", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = "HashedPassword"
	})

	repo.On("Update", mock.AnythingOfType("*users.User")).Return(errors.New("Mock Error"))
	hasher.On("ComparePassword", "mockPassword@123", "HashedPassword").Return(true)
	hasher.On("HashPassword", "mockPassword@1234").Return("HashedPassword2")

	// Mocking TokenAuthMiddleware
	ctx.Set("user_id", "mock_id")

	uh := users.UserHandler{
		Repo:           repo,
		PasswordHasher: hasher,
	}

	uh.ChangePassword(ctx)

	var actualResponse common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &actualResponse)

	repo.AssertCalled(t, "FindById", "mock_id", mock.AnythingOfType("*users.User"))
	repo.AssertCalled(t, "Update", mock.AnythingOfType("*users.User"))
	hasher.AssertCalled(t, "ComparePassword", "mockPassword@123", "HashedPassword")
	hasher.AssertCalled(t, "HashPassword", "mockPassword@1234")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUserHandler_UpdateUser_ShouldUpdateUserSuccessfully(t *testing.T) {
	uurs := []users.UpdateUserRequest{
		{
			Email: "test-update@email.com",
		},
		{
			Username: "test-update-username",
		},
	}

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}
	uh := users.UserHandler{
		Repo:         repo,
		AuthProvider: ap,
	}

	repo.On("FindById", "mock_id", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = mu.Password
	})

	repo.On("EmailAlreadyExists", mock.AnythingOfType("string")).Return(false, nil)
	repo.On("UsernameAlreadyExists", mock.AnythingOfType("string")).Return(false, nil)

	repo.On("Update", mock.AnythingOfType("*users.User")).Return(nil)

	for _, uur := range uurs {
		ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users/me", "PATCH", uur)

		// Mocking TokenAuthMiddleware
		ctx.Set("user_id", "mock_id")

		uh.UpdateUser(ctx)

		repo.AssertCalled(t, "FindById", "mock_id", mock.AnythingOfType("*users.User"))
		repo.AssertCalled(t, "EmailAlreadyExists", mock.AnythingOfType("string"))
		repo.AssertCalled(t, "UsernameAlreadyExists", mock.AnythingOfType("string"))
		repo.AssertCalled(t, "Update", mock.AnythingOfType("*users.User"))

		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestUserHandler_UpdateUser_ShouldNotUpdateUserForExistingUsername(t *testing.T) {
	uur := users.UpdateUserRequest{
		Username: "test-update-username",
	}

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}
	uh := users.UserHandler{
		Repo:         repo,
		AuthProvider: ap,
	}

	repo.On("FindById", "mock_id", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = mu.Password
	})

	repo.On("UsernameAlreadyExists", mock.AnythingOfType("string")).Return(true, nil)

	repo.On("Update", mock.AnythingOfType("*users.User")).Return(nil)

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users/me", "PATCH", uur)

	// Mocking TokenAuthMiddleware
	ctx.Set("user_id", "mock_id")

	uh.UpdateUser(ctx)

	repo.AssertCalled(t, "FindById", "mock_id", mock.AnythingOfType("*users.User"))
	repo.AssertCalled(t, "UsernameAlreadyExists", mock.AnythingOfType("string"))
	repo.AssertNotCalled(t, "EmailAlreadyExists")
	repo.AssertNotCalled(t, "Update")

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Username already exists. Try another.")
}

func TestUserHandler_UpdateUser_ShouldNotUpdateUserForExistingEmail(t *testing.T) {
	uur := users.UpdateUserRequest{
		Email: "test-update@email.com",
	}

	repo := &user_mocks.UserRepository{}
	ap := &utils_mocks.AuthProvider{}
	uh := users.UserHandler{
		Repo:         repo,
		AuthProvider: ap,
	}

	repo.On("FindById", "mock_id", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = mu.Password
	})

	repo.On("UsernameAlreadyExists", mock.AnythingOfType("string")).Return(false, nil)
	repo.On("EmailAlreadyExists", mock.AnythingOfType("string")).Return(true, nil)

	repo.On("Update", mock.AnythingOfType("*users.User")).Return(nil)

	ctx, rec := common.PrepareContextAndResponseRecorder(t, "/api/v1/users/me", "PATCH", uur)

	// Mocking TokenAuthMiddleware
	ctx.Set("user_id", "mock_id")

	uh.UpdateUser(ctx)

	repo.AssertCalled(t, "FindById", "mock_id", mock.AnythingOfType("*users.User"))
	repo.AssertCalled(t, "UsernameAlreadyExists", mock.AnythingOfType("string"))
	repo.AssertCalled(t, "EmailAlreadyExists", mock.AnythingOfType("string"))
	repo.AssertNotCalled(t, "Update")

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Email already exists. Try another.")
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
