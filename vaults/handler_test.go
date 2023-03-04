package vaults_test

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
	"github.com/adarsh-a-tw/passwordly/vaults"
	vaults_mocks "github.com/adarsh-a-tw/passwordly/vaults/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVaultHandler_CreateVault_ShouldCreateVaultSuccessfully(t *testing.T) {
	cvr := vaults.CreateVaultRequest{
		Name: "Mock Vault",
	}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/vaults", "POST", cvr)

	repo := &vaults_mocks.VaultRepository{}
	userRepo := &user_mocks.UserRepository{}

	userRepo.On("FindById", "mock_user_id", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = mu.Password
	})
	repo.On("Create", mock.AnythingOfType("*vaults.Vault")).Return(nil)

	vh := vaults.VaultHandler{
		Repo:     repo,
		UserRepo: userRepo,
	}

	ctx.Set("user_id", "mock_user_id")

	vh.CreateVault(ctx)

	var actualResponse vaults.VaultResponse
	decodeJSONResponse(t, rec, &actualResponse)

	userRepo.AssertCalled(t, "FindById", "mock_user_id", mock.AnythingOfType("*users.User"))
	repo.AssertCalled(t, "Create", mock.AnythingOfType("*vaults.Vault"))

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, cvr.Name, actualResponse.Name)
}

func TestVaultHandler_CreateVault_ShouldThrowInternalServerErrorIfFindByIdMethodFails(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Something went wrong. Try again."}
	cvr := vaults.CreateVaultRequest{
		Name: "Mock Vault",
	}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/vaults", "POST", cvr)

	userRepo := &user_mocks.UserRepository{}

	userRepo.On("FindById", "mock_user_id", mock.AnythingOfType("*users.User")).Return(errors.New("MOCK_ERROR"))

	vh := vaults.VaultHandler{
		UserRepo: userRepo,
	}

	ctx.Set("user_id", "mock_user_id")

	vh.CreateVault(ctx)

	var actualResponse common.ErrorResponse
	decodeJSONResponse(t, rec, &actualResponse)

	userRepo.AssertCalled(t, "FindById", "mock_user_id", mock.AnythingOfType("*users.User"))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestVaultHandler_CreateVault_ShouldThrowErrorForInvalidRequestBody(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Invalid request body"}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/vaults", "POST", nil)

	vh := vaults.VaultHandler{}

	ctx.Set("user_id", "mock_user_id")

	vh.CreateVault(ctx)

	var actualResponse common.ErrorResponse
	decodeJSONResponse(t, rec, &actualResponse)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestVaultHandler_CreateVault_ShouldThrowInternalServerErrorIfCreateMethodFails(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Something went wrong. Try again."}
	cvr := vaults.CreateVaultRequest{
		Name: "Mock Vault",
	}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/vaults", "POST", cvr)

	repo := &vaults_mocks.VaultRepository{}
	userRepo := &user_mocks.UserRepository{}

	userRepo.On("FindById", "mock_user_id", mock.AnythingOfType("*users.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*users.User)
		mu := mockUser()
		arg.Id = mu.Id
		arg.Username = mu.Username
		arg.Email = mu.Email
		arg.Password = mu.Password
	})
	repo.On("Create", mock.AnythingOfType("*vaults.Vault")).Return(errors.New("MOCK_ERROR"))

	userRepo.On("FindById", "mock_user_id", mock.AnythingOfType("*users.User")).Return(errors.New("MOCK_ERROR"))

	vh := vaults.VaultHandler{
		Repo:     repo,
		UserRepo: userRepo,
	}

	ctx.Set("user_id", "mock_user_id")

	vh.CreateVault(ctx)

	var actualResponse common.ErrorResponse
	decodeJSONResponse(t, rec, &actualResponse)

	userRepo.AssertCalled(t, "FindById", "mock_user_id", mock.AnythingOfType("*users.User"))
	repo.AssertCalled(t, "Create", mock.AnythingOfType("*vaults.Vault"))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func mockUser() *users.User {
	return &users.User{
		Id:        "mock_user_id",
		Username:  "mock_username",
		Password:  "HashedPassword",
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
