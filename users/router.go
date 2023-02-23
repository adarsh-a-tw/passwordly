package users

import (
	"github.com/adarsh-a-tw/passwordly/middleware"
	"github.com/adarsh-a-tw/passwordly/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Unauthenticated Routes
	urg := r.Group("/api/v1/users")

	// Authenticated Routes
	rg := r.Group("/api/v1/users")
	rg.Use(middleware.TokenAuthMiddleware)

	uh := UserHandler{
		Repo: &UserRepositoryImpl{
			db: db,
		},
		AuthProvider: &utils.AuthProviderImpl{},
	}

	urg.POST("", uh.Create)
	urg.POST("/login", uh.Login)

	rg.GET("/me", uh.FetchUser)
}
