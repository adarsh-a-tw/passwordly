package secrets

import (
	"github.com/adarsh-a-tw/passwordly/middleware"
	"github.com/adarsh-a-tw/passwordly/users"
	"github.com/adarsh-a-tw/passwordly/vaults"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Authenticated Routes
	rg := r.Group("/api/v1/vaults/:id/secrets")
	rg.Use(middleware.TokenAuthMiddleware)

	repo := &SecretRepositoryImpl{
		Db: db,
	}
	userRepo := &users.UserRepositoryImpl{
		Db: db,
	}
	vaultRepo := &vaults.VaultRepositoryImpl{
		Db: db,
	}

	sh := SecretHandler{
		Repo:      repo,
		UserRepo:  userRepo,
		VaultRepo: vaultRepo,
	}

	rg.POST("", sh.CreateSecret)
}
