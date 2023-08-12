package vaults

import (
	"net/http"

	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/adarsh-a-tw/passwordly/users"
	"github.com/adarsh-a-tw/passwordly/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SecretHandler struct {
	Ep        utils.EncryptionProvider
	Repo      SecretRepository
	VaultRepo VaultRepository
	UserRepo  users.UserRepository
}

func (sh *SecretHandler) CreateSecret(ctx *gin.Context) {
	var csr CreateSecretRequest
	if err := ctx.ShouldBindJSON(&csr); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid Request body"})
		return
	}

	userId := ctx.GetString("user_id")
	var u users.User

	if err := sh.UserRepo.FindById(userId, &u); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	vaultId := ctx.Param("id")

	valid, err := ValidateVaultOwner(sh.VaultRepo, vaultId, userId)

	if err != nil {
		handleGormError(ctx, err)
		return
	}

	if !valid {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	var v Vault
	if err := sh.VaultRepo.FetchById(vaultId, &v); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	switch csr.Type {
	case TypeCredential:
		sh.handleCreateCredential(ctx, &csr, &v)
	default:
	}
}

// private methods
func (sh *SecretHandler) handleCreateCredential(ctx *gin.Context, csr *CreateSecretRequest, v *Vault) {
	if csr.Username == "" || csr.Password == "" {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid Request body"})
		return
	}
	ep, err := sh.Ep.Encrypt(csr.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}
	c := Credential{
		Id:       uuid.NewString(),
		Name:     csr.Name,
		Username: csr.Username,
		Password: ep,
		Vault:    *v,
	}
	if err := sh.Repo.CreateCredential(&c); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}
	c.Password = csr.Password
	sr := SecretResponse{}
	sr.load(c)

	ctx.JSON(
		http.StatusCreated,
		sr,
	)
}
