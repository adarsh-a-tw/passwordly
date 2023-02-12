package users

import (
	"net/http"

	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/adarsh-a-tw/passwordly/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	ur UserRepository
	ap utils.AuthProvider
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

func (uh *UserHandler) Login(ctx *gin.Context) {
	var lur LoginUserRequest
	if err := ctx.ShouldBindJSON(&lur); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid request body"})
		return
	}

	var u User
	if err := uh.ur.Find(lur.Username, &u); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid Credentials"})
		return
	}

	if u.Password == lur.Password {
		tokenStr, err := uh.ap.GenerateToken(u.Id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: "Something went wrong. Try again."})
			return
		}
		ctx.JSON(http.StatusOK, LoginUserSuccessResponse{Token: tokenStr})
		return
	}

	ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid Credentials"})
}
