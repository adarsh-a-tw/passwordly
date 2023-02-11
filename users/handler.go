package users

import (
	"net/http"

	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	ur UserRepository
}

func (uh *UserHandler) Create(ctx *gin.Context) {
	var cur CreateUserRequest
	if err := ctx.ShouldBindJSON(&cur); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid request body"})
		return
	}

	u := User{
		Id:       uuid.NewString(),
		Username: cur.Username,
		Email:    cur.Email,
		Password: cur.Password,
	}

	if err := uh.ur.Create(&u); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, UserResponse{
		Id:        u.Id,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	})
}
