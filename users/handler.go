package users

import (
	"net/http"

	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/adarsh-a-tw/passwordly/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	Repo         UserRepository
	AuthProvider utils.AuthProvider
}

func (uh *UserHandler) Create(ctx *gin.Context) {
	var cur CreateUserRequest
	if err := ctx.ShouldBindJSON(&cur); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid request body" + err.Error()})
		return
	}

	u := User{
		Id:       uuid.NewString(),
		Username: cur.Username,
		Email:    cur.Email,
		Password: cur.Password,
	}

	if err := uh.Repo.Create(&u); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	ctx.JSON(http.StatusCreated, UserResponse{
		Id:       u.Id,
		Username: u.Username,
		Email:    u.Email,
	})
}

func (uh *UserHandler) Login(ctx *gin.Context) {
	var lur LoginUserRequest
	if err := ctx.ShouldBindJSON(&lur); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid request body"})
		return
	}

	var u User
	if err := uh.Repo.Find(lur.Username, &u); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid Credentials"})
		return
	}

	if u.Password == lur.Password {
		tokenStr, err := uh.AuthProvider.GenerateToken(u.Id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
			return
		}
		ctx.JSON(http.StatusOK, LoginUserSuccessResponse{Token: tokenStr})
		return
	}

	ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid Credentials"})
}

func (uh *UserHandler) FetchUser(ctx *gin.Context) {
	id := ctx.GetString("user_id")
	var u User

	if err := uh.Repo.FindById(id, &u); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	ctx.JSON(
		http.StatusOK,
		UserResponse{
			Id:       u.Id,
			Username: u.Username,
			Email:    u.Email,
		},
	)

}

func (uh *UserHandler) ChangePassword(ctx *gin.Context) {
	id := ctx.GetString("user_id")
	var cpr ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&cpr); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid Request body"})
		return
	}

	var u User
	if err := uh.Repo.FindById(id, &u); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	if cpr.CurrentPassword != u.Password {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Current password does not match"})
		return
	}

	if cpr.NewPassword == u.Password {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "New password and current password should not be the same"})
		return
	}

	u.Password = cpr.NewPassword

	if err := uh.Repo.Update(&u); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
