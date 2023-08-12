package vaults

import (
	"github.com/adarsh-a-tw/passwordly/middleware"
	"github.com/adarsh-a-tw/passwordly/users"
	"github.com/adarsh-a-tw/passwordly/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Authenticated Routes
	rg := r.Group("/api/v1/vaults")
	rg.Use(middleware.TokenAuthMiddleware)

	vaultsRepo := &VaultRepositoryImpl{
		Db: db,
	}
	userRepo := &users.UserRepositoryImpl{
		Db: db,
	}
	secretRepo := &SecretRepositoryImpl{Db: db}

	ep, _ := utils.NewEncryptionProvider()

	vh := VaultHandler{
		Ep:         ep,
		Repo:       vaultsRepo,
		UserRepo:   userRepo,
		SecretRepo: secretRepo,
	}

	sh := SecretHandler{
		Ep:        ep,
		Repo:      secretRepo,
		UserRepo:  userRepo,
		VaultRepo: vaultsRepo,
	}

	rg.POST("", vh.CreateVault)
	rg.GET("", vh.FetchVaults)

	rg.GET("/:id", vh.FetchVaultDetails)
	rg.PATCH("/:id", vh.UpdateVault)
	rg.DELETE("/:id", vh.DeleteVault)

	rg.POST("/:id/secrets", sh.CreateSecret)
}
