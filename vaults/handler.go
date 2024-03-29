package vaults

import (
	"errors"
	"net/http"

	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/adarsh-a-tw/passwordly/users"
	"github.com/adarsh-a-tw/passwordly/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VaultHandler struct {
	Ep         utils.EncryptionProvider
	Repo       VaultRepository
	UserRepo   users.UserRepository
	SecretRepo SecretRepository
}

func (vh *VaultHandler) CreateVault(ctx *gin.Context) {
	var cvr CreateVaultRequest
	if err := ctx.ShouldBindJSON(&cvr); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid request body"})
		return
	}

	id := ctx.GetString("user_id")
	var u users.User

	if err := vh.UserRepo.FindById(id, &u); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	v := Vault{
		Id:   uuid.NewString(),
		Name: cvr.Name,
		User: u,
	}

	if err := vh.Repo.Create(&v); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	ctx.JSON(http.StatusCreated, VaultResponse{
		Id:        v.Id,
		Name:      v.Name,
		Secrets:   []SecretResponse{},
		CreatedAt: v.CreatedAt.Unix(),
		UpdatedAt: v.UpdatedAt.Unix(),
	})
}

func (vh *VaultHandler) FetchVaults(ctx *gin.Context) {
	id := ctx.GetString("user_id")
	var u users.User

	if err := vh.UserRepo.FindById(id, &u); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	var vaults []Vault
	if err := vh.Repo.FetchByUserId(u.Id, &vaults); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	response := VaultListResponse{}
	response.load(vaults)

	ctx.JSON(http.StatusOK, response)
}

func (vh *VaultHandler) FetchVaultDetails(ctx *gin.Context) {
	userId := ctx.GetString("user_id")
	var u users.User

	if err := vh.UserRepo.FindById(userId, &u); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	vaultId := ctx.Param("id")
	valid, err := ValidateVaultOwner(vh.Repo, vaultId, userId)

	if err != nil {
		handleGormError(ctx, err)
		return
	}

	if !valid {
		ctx.JSON(http.StatusUnauthorized, common.ErrorResponse{Message: "Unauthorized access"})
		return
	}

	var vault Vault
	if err := vh.Repo.FetchById(vaultId, &vault); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	var credentials []Credential
	if err = vh.SecretRepo.FindCredentials(&credentials, vaultId); err != nil {
		handleGormError(ctx, err)
		return
	}
	if err = vh.decryptCredentials(credentials); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	vr := VaultResponse{}
	vr.load(vault, credentials)

	ctx.JSON(http.StatusOK, vr)
}

func (vh *VaultHandler) UpdateVault(ctx *gin.Context) {
	userId := ctx.GetString("user_id")
	vaultId := ctx.Param("id")

	valid, err := ValidateVaultOwner(vh.Repo, vaultId, userId)

	if err != nil {
		handleGormError(ctx, err)
		return
	}

	if !valid {
		ctx.JSON(http.StatusUnauthorized, common.ErrorResponse{Message: "Unauthorized access"})
		return
	}

	var vault Vault

	if err := vh.Repo.FetchById(vaultId, &vault); err != nil {
		handleGormError(ctx, err)
		return
	}

	var uvr UpdateVaultRequest

	if err := ctx.ShouldBindJSON(&uvr); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: err.Error()})
		return
	}

	vault.Name = uvr.Name

	if err := vh.Repo.Update(&vault); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Vault updated successfully"})
}

func (vh *VaultHandler) DeleteVault(ctx *gin.Context) {
	userId := ctx.GetString("user_id")
	vaultId := ctx.Param("id")

	valid, err := ValidateVaultOwner(vh.Repo, vaultId, userId)

	if err != nil {
		handleGormError(ctx, err)
		return
	}

	if !valid {
		ctx.JSON(http.StatusUnauthorized, common.ErrorResponse{Message: "Unauthorized access"})
		return
	}

	var vault Vault

	if err := vh.Repo.FetchById(vaultId, &vault); err != nil {
		handleGormError(ctx, err)
		return
	}

	if err := vh.Repo.Delete(&vault); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Vault deleted successfully"})
}

func handleGormError(ctx *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.AbortWithStatus(http.StatusNotFound)
	} else {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
	}
}

func (vh *VaultHandler) decryptCredentials(creds []Credential) error {
	for i := range creds {
		pwd, err := vh.Ep.Decrypt(string(creds[i].Password))
		if err != nil {
			return err
		}
		creds[i].Password = []byte(pwd)
	}
	return nil
}
