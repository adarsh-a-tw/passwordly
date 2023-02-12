package users

import (
	"github.com/adarsh-a-tw/passwordly/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Unauthenticated Routes
	rg := r.Group("/api/v1/users")

	uh := UserHandler{
		Repo: &UserRepositoryImpl{
			db: db,
		},
		AuthProvider: &utils.AuthProviderImpl{},
	}

	rg.POST("", uh.Create)
	rg.POST("/login", uh.Login)
}
