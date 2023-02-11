package users

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	rg := r.Group("/v1/users")

	uh := UserHandler{
		ur: &UserRepositoryImpl{
			db: db,
		},
	}

	rg.POST("", uh.Create)
}
