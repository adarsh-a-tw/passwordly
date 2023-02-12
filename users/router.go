package users

import (
	"github.com/adarsh-a-tw/passwordly/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Unauthenticated Routes
	rg := r.Group("/v1/users")

	uh := UserHandler{
		ur: &UserRepositoryImpl{
			db: db,
		},
		ap: &utils.AuthProviderImpl{},
	}

	rg.POST("", uh.Create)
	rg.POST("/login", uh.Login)
}
