package vaults

import (
	"net/http"

	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/adarsh-a-tw/passwordly/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VaultHandler struct {
	Repo     VaultRepository
	UserRepo users.UserRepository
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
		Id:   v.Id,
		Name: v.Name,
	})
}
