package vaults

import (
	"net/http"

	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/adarsh-a-tw/passwordly/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SecretHandler struct {
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

	s := Secret{
		Id:    uuid.NewString(),
		Name:  csr.Name,
		Type:  string(csr.Type),
		Vault: v,
	}

	switch csr.Type {
	case TypeCredential:
		sh.handleCreateCredential(ctx, &csr, &s)
	default:
	}
}

// private methods
func (sh *SecretHandler) handleCreateCredential(ctx *gin.Context, csr *CreateSecretRequest, s *Secret) {
	if csr.Username == "" || csr.Password == "" {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid Request body"})
		return
	}
	c := Credential{
		Id:       uuid.NewString(),
		Username: csr.Username,
		Password: csr.Password,
		Secret:   *s,
	}
	if err := sh.Repo.CreateCredential(s, &c); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}
	ctx.JSON(
		http.StatusCreated,
		SecretResponse{
			Id:        s.Id,
			Name:      s.Name,
			Type:      s.SecretType(),
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
			Username:  c.Username,
			Password:  c.Password,
		},
	)
}
