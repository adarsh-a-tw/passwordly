package secrets_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/adarsh-a-tw/passwordly/common"
	s "github.com/adarsh-a-tw/passwordly/secrets"
	sm "github.com/adarsh-a-tw/passwordly/secrets/mocks"
	"github.com/adarsh-a-tw/passwordly/users"
	um "github.com/adarsh-a-tw/passwordly/users/mocks"
	"github.com/adarsh-a-tw/passwordly/vaults"
	vm "github.com/adarsh-a-tw/passwordly/vaults/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var mockUser users.User
var mockUser2 users.User
var mockVault vaults.Vault
var mockVault2 vaults.Vault

func init() {
	mockUser = users.User{
		Id: "test-user-id",
	}
	mockUser2 = users.User{
		Id: "test-user-id-2",
	}
	mockVault = vaults.Vault{
		Id:        "mock-vault",
		UserRefer: mockUser.Id,
		User:      mockUser,
	}
	mockVault2 = vaults.Vault{
		Id:        "mock-vault-2",
		UserRefer: mockUser2.Id,
		User:      mockUser2,
	}
	s.RegisterValidations()
}

func TestSecretHandler_CreateSecret_ShouldCreateSecretOfTypeCredential(t *testing.T) {
	csr := s.CreateSecretRequest{
		Name:     "test-secret",
		Type:     s.TypeCredential,
		Username: "test",
		Password: "test",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, fmt.Sprintf("/api/v1/vaults/%s/secrets", mockVault.Id), "POST", csr)
	ctx.Set("user_id", mockUser.Id)
	ctx.AddParam("id", mockVault.Id)

	msr := sm.SecretRepository{}
	msr.On(
		"CreateCredential",
		mock.AnythingOfType("*secrets.Secret"),
		mock.AnythingOfType("*secrets.Credential"),
	).Return(nil)

	mvr := vm.VaultRepository{}
	mvr.On("FetchById", mockVault.Id, mock.AnythingOfType("*vaults.Vault")).Run(func(args mock.Arguments) {
		v := args.Get(1).(*vaults.Vault)
		v.Id = mockVault.Id
		v.UserRefer = mockUser.Id
		v.User = mockUser
	}).Return(nil)

	mur := um.UserRepository{}
	mur.On("FindById", mockUser.Id, mock.AnythingOfType("*users.User")).Run(func(args mock.Arguments) {
		u := args.Get(1).(*users.User)
		u.Id = mockUser.Id
	}).Return(nil)

	h := s.SecretHandler{
		Repo:      &msr,
		VaultRepo: &mvr,
		UserRepo:  &mur,
	}

	h.CreateSecret(ctx)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp s.SecretResponse
	common.DecodeJSONResponse(t, rec, &resp)

	assert.Equal(t, csr.Name, resp.Name)
	assert.Equal(t, csr.Username, resp.Username)
	assert.Equal(t, csr.Password, resp.Password)
	assert.Equal(t, csr.Type, resp.Type)
}

func TestSecretHandler_CreateSecret_ShouldFailForEmptyRequestbody(t *testing.T) {
	ctx, rec := common.PrepareContextAndResponseRecorder(t, fmt.Sprintf("/api/v1/vaults/%s/secrets", mockVault.Id), "POST", nil)

	h := s.SecretHandler{}

	h.CreateSecret(ctx)

	var resp common.ErrorResponse
	common.DecodeJSONResponse(t, rec, &resp)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, resp.Message, "Invalid Request body")
}

func TestSecretHandler_CreateSecret_ShouldFailForInvalidRequestBodyForTypeCredential(t *testing.T) {
	csr := s.CreateSecretRequest{
		Name: "test-secret",
		Type: s.TypeCredential,
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, fmt.Sprintf("/api/v1/vaults/%s/secrets", mockVault.Id), "POST", csr)
	ctx.Set("user_id", mockUser.Id)
	ctx.AddParam("id", mockVault.Id)

	mvr := vm.VaultRepository{}
	mvr.On("FetchById", mockVault.Id, mock.AnythingOfType("*vaults.Vault")).Run(func(args mock.Arguments) {
		v := args.Get(1).(*vaults.Vault)
		v.Id = mockVault.Id
		v.UserRefer = mockUser.Id
		v.User = mockUser
	}).Return(nil)

	mur := um.UserRepository{}
	mur.On("FindById", mockUser.Id, mock.AnythingOfType("*users.User")).Run(func(args mock.Arguments) {
		u := args.Get(1).(*users.User)
		u.Id = mockUser.Id
	}).Return(nil)

	h := s.SecretHandler{
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

	csr := s.CreateSecretRequest{
		Name:     "test-secret",
		Type:     s.TypeCredential,
		Username: "test",
		Password: "test",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, fmt.Sprintf("/api/v1/vaults/%s/secrets", "mock-id"), "POST", csr)
	ctx.Set("user_id", mockUser.Id)
	ctx.AddParam("id", mockVault.Id)

	mvr := vm.VaultRepository{}
	mvr.On("FetchById", mockVault.Id, mock.AnythingOfType("*vaults.Vault")).Return(gorm.ErrRecordNotFound)

	mur := um.UserRepository{}
	mur.On("FindById", mockUser.Id, mock.AnythingOfType("*users.User")).Run(func(args mock.Arguments) {
		u := args.Get(1).(*users.User)
		u.Id = mockUser.Id
	}).Return(nil)

	h := s.SecretHandler{
		VaultRepo: &mvr,
		UserRepo:  &mur,
	}

	h.CreateSecret(ctx)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestSecretHandler_CreateSecret_ShouldFailForInvalidVaultOwner(t *testing.T) {

	csr := s.CreateSecretRequest{
		Name:     "test-secret",
		Type:     s.TypeCredential,
		Username: "test",
		Password: "test",
	}

	ctx, rec := common.PrepareContextAndResponseRecorder(t, fmt.Sprintf("/api/v1/vaults/%s/secrets", mockVault2.Id), "POST", csr)
	ctx.Set("user_id", mockUser.Id)
	ctx.AddParam("id", mockVault2.Id)

	mvr := vm.VaultRepository{}
	mvr.On("FetchById", mockVault2.Id, mock.AnythingOfType("*vaults.Vault")).Run(func(args mock.Arguments) {
		v := args.Get(1).(*vaults.Vault)
		v.Id = mockVault2.Id
		v.UserRefer = mockUser2.Id
		v.User = mockUser2
	}).Return(nil)

	mur := um.UserRepository{}
	mur.On("FindById", mockUser.Id, mock.AnythingOfType("*users.User")).Run(func(args mock.Arguments) {
		u := args.Get(1).(*users.User)
		u.Id = mockUser.Id
	}).Return(nil)

	h := s.SecretHandler{
		VaultRepo: &mvr,
		UserRepo:  &mur,
	}

	h.CreateSecret(ctx)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
