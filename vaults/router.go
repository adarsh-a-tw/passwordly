package vaults

import (
	"github.com/adarsh-a-tw/passwordly/middleware"
	"github.com/adarsh-a-tw/passwordly/users"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Authenticated Routes
	rg := r.Group("/api/v1/vaults")
	rg.Use(middleware.TokenAuthMiddleware)

	repo := &VaultRepositoryImpl{
		Db: db,
	}
	userRepo := &users.UserRepositoryImpl{
		Db: db,
	}

	vh := VaultHandler{
		Repo:     repo,
		UserRepo: userRepo,
	}

	rg.POST("", vh.CreateVault)
	rg.GET("", vh.FetchVaults)

	rg.PATCH("/:id", vh.UpdateVault)
	rg.DELETE("/:id", vh.DeleteVault)
}
