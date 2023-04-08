package secrets_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/adarsh-a-tw/passwordly/secrets"
	s "github.com/adarsh-a-tw/passwordly/secrets"
	sm "github.com/adarsh-a-tw/passwordly/secrets/mocks"
	"github.com/adarsh-a-tw/passwordly/users"
	um "github.com/adarsh-a-tw/passwordly/users/mocks"
	"github.com/adarsh-a-tw/passwordly/vaults"
	vm "github.com/adarsh-a-tw/passwordly/vaults/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockUser users.User
var mockVault vaults.Vault

func init() {
	mockUser = users.User{
		Id: "test-user-id",
	}
	mockVault = vaults.Vault{
		Id:        "mock-vault",
		UserRefer: mockUser.Id,
		User:      mockUser,
	}
	secrets.RegisterValidations()
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

	var resp secrets.SecretResponse
	common.DecodeJSONResponse(t, rec, &resp)

	assert.Equal(t, csr.Name, resp.Name)
	assert.Equal(t, csr.Username, resp.Username)
	assert.Equal(t, csr.Password, resp.Password)
	assert.Equal(t, csr.Type, resp.Type)
}
