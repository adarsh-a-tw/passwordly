package vaults_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	"gorm.io/gorm"
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

func TestVaultHandler_FetchVaults_ShouldFetchVaultsSuccessfully(t *testing.T) {
	mv := mockVaults()
	mockVault1 := (*mv)[0]
	mockVault2 := (*mv)[1]

	expectedResponse := vaults.VaultListResponse{
		Vaults: []vaults.VaultResponse{
			{
				Id:   mockVault1.Id,
				Name: mockVault1.Name,
			},
			{
				Id:   mockVault2.Id,
				Name: mockVault2.Name,
			},
		},
	}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/vaults", "GET", nil)

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
	repo.On("FetchByUserId", "mock_user_id", mock.AnythingOfType("*[]vaults.Vault")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*[]vaults.Vault)
		(*arg) = append((*arg), mockVault1)
		(*arg) = append((*arg), mockVault2)
	})

	vh := vaults.VaultHandler{
		Repo:     repo,
		UserRepo: userRepo,
	}

	ctx.Set("user_id", "mock_user_id")

	vh.FetchVaults(ctx)

	var actualResponse vaults.VaultListResponse
	decodeJSONResponse(t, rec, &actualResponse)

	userRepo.AssertCalled(t, "FindById", "mock_user_id", mock.AnythingOfType("*users.User"))
	repo.AssertCalled(t, "FetchByUserId", "mock_user_id", mock.AnythingOfType("*[]vaults.Vault"))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestVaultHandler_FetchVaults_ShouldThrowInternalServerErrorIfFindByIdMethodFails(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Something went wrong. Try again."}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/vaults", "GET", nil)

	userRepo := &user_mocks.UserRepository{}

	userRepo.On("FindById", "mock_user_id", mock.AnythingOfType("*users.User")).Return(errors.New("MOCK_ERROR"))

	vh := vaults.VaultHandler{
		UserRepo: userRepo,
	}

	ctx.Set("user_id", "mock_user_id")

	vh.FetchVaults(ctx)

	var actualResponse common.ErrorResponse
	decodeJSONResponse(t, rec, &actualResponse)

	userRepo.AssertCalled(t, "FindById", "mock_user_id", mock.AnythingOfType("*users.User"))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestVaultHandler_FetchVaults_ShouldThrowInternalServerErrorIfFetchByUserIdMethodFails(t *testing.T) {
	expectedResponse := common.ErrorResponse{Message: "Something went wrong. Try again."}

	ctx, rec := prepareContextAndResponseRecorder(t, "/api/v1/vaults", "GET", nil)

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
	repo.On("FetchByUserId", "mock_user_id", mock.AnythingOfType("*[]vaults.Vault")).Return(errors.New("MOCK_ERROR"))

	vh := vaults.VaultHandler{
		Repo:     repo,
		UserRepo: userRepo,
	}

	ctx.Set("user_id", "mock_user_id")

	vh.FetchVaults(ctx)

	var actualResponse common.ErrorResponse
	decodeJSONResponse(t, rec, &actualResponse)

	userRepo.AssertCalled(t, "FindById", "mock_user_id", mock.AnythingOfType("*users.User"))
	repo.AssertCalled(t, "FetchByUserId", "mock_user_id", mock.AnythingOfType("*[]vaults.Vault"))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestVaultHandler_UpdateVault_ShouldUpdateVaultSuccessfully(t *testing.T) {

	uvr := vaults.UpdateVaultRequest{
		Name: "test-vault-renamed",
	}
	existingVault := (*mockVaults())[0]
	userId := "mock_user_id"

	ctx, rec := prepareContextAndResponseRecorder(t, fmt.Sprintf("/api/v1/vaults/%s", existingVault.Id), "PATCH", uvr)
	repo := &vaults_mocks.VaultRepository{}
	ctx.Set("user_id", userId)
	ctx.AddParam("id", existingVault.Id)

	repo.On("FetchById", existingVault.Id, mock.AnythingOfType("*vaults.Vault")).Return(nil).Run(func(args mock.Arguments) {
		vault := args.Get(1).(*vaults.Vault)
		vault.Id = existingVault.Id
		vault.Name = existingVault.Name
		vault.UserRefer = userId
	})

	repo.On("Update", mock.AnythingOfType("*vaults.Vault")).Return(nil)

	vh := vaults.VaultHandler{
		Repo:     repo,
		UserRepo: &user_mocks.UserRepository{},
	}

	vh.UpdateVault(ctx)

	repo.AssertNumberOfCalls(t, "FetchById", 2)
	repo.AssertNumberOfCalls(t, "Update", 1)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, `{"message":"Vault updated successfully"}`, rec.Body.String())
}

func TestVaultHandler_UpdateVault_ShouldNotUpdateSuccessfullyWhenVaultOwnerIsNotRequester(t *testing.T) {

	uvr := vaults.UpdateVaultRequest{
		Name: "test-vault-renamed",
	}
	existingVault := (*mockVaults())[0]
	userId := "mock_user_id"

	ctx, rec := prepareContextAndResponseRecorder(t, fmt.Sprintf("/api/v1/vaults/%s", existingVault.Id), "PATCH", uvr)
	repo := &vaults_mocks.VaultRepository{}
	ctx.Set("user_id", userId)
	ctx.AddParam("id", existingVault.Id)

	repo.On("FetchById", existingVault.Id, mock.AnythingOfType("*vaults.Vault")).Return(nil).Run(func(args mock.Arguments) {
		vault := args.Get(1).(*vaults.Vault)
		vault.Id = existingVault.Id
		vault.Name = existingVault.Name
		vault.UserRefer = "mock_different_user_id"
	})

	vh := vaults.VaultHandler{
		Repo:     repo,
		UserRepo: &user_mocks.UserRepository{},
	}

	vh.UpdateVault(ctx)

	repo.AssertNumberOfCalls(t, "FetchById", 1)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestVaultHandler_UpdateVault_ShouldNotUpdateSuccessfullyIfFetchByIdMethodFails(t *testing.T) {
	uvr := vaults.UpdateVaultRequest{
		Name: "test-vault-renamed",
	}
	existingVault := (*mockVaults())[0]
	userId := "mock_user_id"
	ctx, rec := prepareContextAndResponseRecorder(t, fmt.Sprintf("/api/v1/vaults/%s", existingVault.Id), "PATCH", uvr)
	repo := &vaults_mocks.VaultRepository{}
	ctx.Set("user_id", userId)
	ctx.AddParam("id", existingVault.Id)
	repo.On("FetchById", existingVault.Id, mock.AnythingOfType("*vaults.Vault")).Return(errors.New("mock FetchById fail"))

	vh := vaults.VaultHandler{
		Repo:     repo,
		UserRepo: &user_mocks.UserRepository{},
	}
	vh.UpdateVault(ctx)

	assert.Equal(t, `{"message":"Something went wrong. Try again."}`, rec.Body.String())
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestVaultHandler_UpdateVault_ShouldNotUpdateSuccessfullyIfFetchByIdReturnsNotFoundError(t *testing.T) {
	uvr := vaults.UpdateVaultRequest{
		Name: "test-vault-renamed",
	}
	existingVault := (*mockVaults())[0]
	userId := "mock_user_id"
	ctx, rec := prepareContextAndResponseRecorder(t, fmt.Sprintf("/api/v1/vaults/%s", existingVault.Id), "PATCH", uvr)
	repo := &vaults_mocks.VaultRepository{}
	ctx.Set("user_id", userId)
	ctx.AddParam("id", existingVault.Id)
	repo.On("FetchById", existingVault.Id, mock.AnythingOfType("*vaults.Vault")).Return(gorm.ErrRecordNotFound)

	vh := vaults.VaultHandler{
		Repo:     repo,
		UserRepo: &user_mocks.UserRepository{},
	}
	vh.UpdateVault(ctx)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestVaultHandler_UpdateVault_ShouldNotUpdateSuccessfullyIfUpdateMethodFails(t *testing.T) {

	uvr := vaults.UpdateVaultRequest{
		Name: "test-vault-renamed",
	}
	existingVault := (*mockVaults())[0]
	userId := "mock_user_id"

	ctx, rec := prepareContextAndResponseRecorder(t, fmt.Sprintf("/api/v1/vaults/%s", existingVault.Id), "PATCH", uvr)
	repo := &vaults_mocks.VaultRepository{}
	ctx.Set("user_id", userId)
	ctx.AddParam("id", existingVault.Id)

	repo.On("FetchById", existingVault.Id, mock.AnythingOfType("*vaults.Vault")).Return(nil).Run(func(args mock.Arguments) {
		vault := args.Get(1).(*vaults.Vault)
		vault.Id = existingVault.Id
		vault.Name = existingVault.Name
		vault.UserRefer = userId
	})

	repo.On("Update", mock.AnythingOfType("*vaults.Vault")).Return(errors.New("mock error"))

	vh := vaults.VaultHandler{
		Repo:     repo,
		UserRepo: &user_mocks.UserRepository{},
	}
	vh.UpdateVault(ctx)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
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

func mockVaults() *[]vaults.Vault {
	return &[]vaults.Vault{
		{
			Id:        "mock_vault_id_1",
			Name:      "Mock Vault 1",
			User:      *mockUser(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Id:        "mock_vault_id_2",
			Name:      "Mock Vault 2",
			User:      *mockUser(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
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
