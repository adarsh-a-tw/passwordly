package vaults_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/adarsh-a-tw/passwordly/users"
	um "github.com/adarsh-a-tw/passwordly/users/mocks"
	v "github.com/adarsh-a-tw/passwordly/vaults"
	vm "github.com/adarsh-a-tw/passwordly/vaults/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var mockUser1 users.User
var mockUser2 users.User
var mockVault v.Vault
var mockVault2 v.Vault

func init() {
	mockUser1 = users.User{
		Id: "test-user-id",
	}
	mockUser2 = users.User{
		Id: "test-user-id-2",
	}
	mockVault = v.Vault{
		Id:        "mock-vault",
		UserRefer: mockUser1.Id,
		User:      mockUser1,
	}
	mockVault2 = v.Vault{
		Id:        "mock-vault-2",
		UserRefer: mockUser2.Id,
		User:      mockUser2,
	}
	v.RegisterValidations()
}

func TestSecretHandler_CreateSecret_ShouldCreateSecretOfTypeCredential(t *testing.T) {
	csr := v.CreateSecretRequest{
		Name:     "test-secret",
		Type:     v.TypeCredential,
		Username: "test",
		Password: "test",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, fmt.Sprintf("/api/v1/vaults/%s/secrets", mockVault.Id), "POST", csr)
	ctx.Set("user_id", mockUser1.Id)
	ctx.AddParam("id", mockVault.Id)

	msr := vm.SecretRepository{}
	msr.On(
		"CreateCredential",
		mock.AnythingOfType("*vaults.Credential"),
	).Return(nil)

	mvr := vm.VaultRepository{}
	mvr.On("FetchById", mockVault.Id, mock.AnythingOfType("*vaults.Vault")).Run(func(args mock.Arguments) {
		v := args.Get(1).(*v.Vault)
		v.Id = mockVault.Id
		v.UserRefer = mockUser1.Id
		v.User = mockUser1
	}).Return(nil)

	mur := um.UserRepository{}
	mur.On("FindById", mockUser1.Id, mock.AnythingOfType("*users.User")).Run(func(args mock.Arguments) {
		u := args.Get(1).(*users.User)
		u.Id = mockUser1.Id
	}).Return(nil)

	h := v.SecretHandler{
		Repo:      &msr,
		VaultRepo: &mvr,
		UserRepo:  &mur,
	}

	h.CreateSecret(ctx)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp v.SecretResponse
	common.DecodeJSONResponse(t, rec, &resp)

	assert.Equal(t, csr.Name, resp.Name)
	assert.Equal(t, csr.Username, resp.Username)
	assert.Equal(t, csr.Password, resp.Password)
	assert.Equal(t, csr.Type, resp.Type)
}

func TestSecretHandler_CreateSecret_ShouldFailForEmptyRequestbody(t *testing.T) {
	ctx, rec := common.PrepareContextAndResponseRecorder(t, fmt.Sprintf("/api/v1/vaults/%s/secrets", mockVault.Id), "POST", nil)

	h := v.SecretHandler{}

	h.CreateSecret(ctx)

	var resp common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &resp)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, resp.Message, "Invalid Request body")
}

func TestSecretHandler_CreateSecret_ShouldFailForInvalidRequestBodyForTypeCredential(t *testing.T) {
	csr := v.CreateSecretRequest{
		Name: "test-secret",
		Type: v.TypeCredential,
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, fmt.Sprintf("/api/v1/vaults/%s/secrets", mockVault.Id), "POST", csr)
	ctx.Set("user_id", mockUser1.Id)
	ctx.AddParam("id", mockVault.Id)

	mvr := vm.VaultRepository{}
	mvr.On("FetchById", mockVault.Id, mock.AnythingOfType("*vaults.Vault")).Run(func(args mock.Arguments) {
		v := args.Get(1).(*v.Vault)
		v.Id = mockVault.Id
		v.UserRefer = mockUser1.Id
		v.User = mockUser1
	}).Return(nil)

	mur := um.UserRepository{}
	mur.On("FindById", mockUser1.Id, mock.AnythingOfType("*users.User")).Run(func(args mock.Arguments) {
		u := args.Get(1).(*users.User)
		u.Id = mockUser1.Id
	}).Return(nil)

	h := v.SecretHandler{
		VaultRepo: &mvr,
		UserRepo:  &mur,
	}

	h.CreateSecret(ctx)

	var resp common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &resp)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, resp.Message, "Invalid Request body")
}

func TestSecretHandler_CreateSecret_ShouldFailForInvalidVaultId(t *testing.T) {

	csr := v.CreateSecretRequest{
		Name:     "test-secret",
		Type:     v.TypeCredential,
		Username: "test",
		Password: "test",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, fmt.Sprintf("/api/v1/vaults/%s/secrets", "mock-id"), "POST", csr)
	ctx.Set("user_id", mockUser1.Id)
	ctx.AddParam("id", mockVault.Id)

	mvr := vm.VaultRepository{}
	mvr.On("FetchById", mockVault.Id, mock.AnythingOfType("*vaults.Vault")).Return(gorm.ErrRecordNotFound)

	mur := um.UserRepository{}
	mur.On("FindById", mockUser1.Id, mock.AnythingOfType("*users.User")).Run(func(args mock.Arguments) {
		u := args.Get(1).(*users.User)
		u.Id = mockUser1.Id
	}).Return(nil)

	h := v.SecretHandler{
		VaultRepo: &mvr,
		UserRepo:  &mur,
	}

	h.CreateSecret(ctx)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestSecretHandler_CreateSecret_ShouldFailForInvalidVaultOwner(t *testing.T) {

	csr := v.CreateSecretRequest{
		Name:     "test-secret",
		Type:     v.TypeCredential,
		Username: "test",
		Password: "test",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, fmt.Sprintf("/api/v1/vaults/%s/secrets", mockVault2.Id), "POST", csr)
	ctx.Set("user_id", mockUser1.Id)
	ctx.AddParam("id", mockVault2.Id)

	mvr := vm.VaultRepository{}
	mvr.On("FetchById", mockVault2.Id, mock.AnythingOfType("*vaults.Vault")).Run(func(args mock.Arguments) {
		v := args.Get(1).(*v.Vault)
		v.Id = mockVault2.Id
		v.UserRefer = mockUser2.Id
		v.User = mockUser2
	}).Return(nil)

	mur := um.UserRepository{}
	mur.On("FindById", mockUser1.Id, mock.AnythingOfType("*users.User")).Run(func(args mock.Arguments) {
		u := args.Get(1).(*users.User)
		u.Id = mockUser1.Id
	}).Return(nil)

	h := v.SecretHandler{
		VaultRepo: &mvr,
		UserRepo:  &mur,
	}

	h.CreateSecret(ctx)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
