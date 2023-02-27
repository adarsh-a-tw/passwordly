package users

import (
	"fmt"
	"net/http"

	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/adarsh-a-tw/passwordly/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	Repo           UserRepository
	AuthProvider   utils.AuthProvider
	PasswordHasher utils.PasswordHasher
}

func (uh *UserHandler) Create(ctx *gin.Context) {
	var cur CreateUserRequest
	if err := ctx.ShouldBindJSON(&cur); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid Request body"})
		return
	}

	if exists, err := uh.Repo.UsernameAlreadyExists(cur.Username); exists {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Username already exists. Try another."})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	if exists, err := uh.Repo.EmailAlreadyExists(cur.Email); exists {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Email already exists. Try another."})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	hashedPassword := uh.PasswordHasher.HashPassword(cur.Password)

	u := User{
		Id:       uuid.NewString(),
		Username: cur.Username,
		Email:    cur.Email,
		Password: hashedPassword,
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

	fmt.Println(uh.PasswordHasher.ComparePassword(lur.Password, u.Password))

	if uh.PasswordHasher.ComparePassword(lur.Password, u.Password) {
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

	if !uh.PasswordHasher.ComparePassword(cpr.CurrentPassword, u.Password) {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Current password does not match"})
		return
	}

	if cpr.NewPassword == cpr.CurrentPassword {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "New password and current password should not be the same"})
		return
	}

	u.Password = uh.PasswordHasher.HashPassword(cpr.NewPassword)

	if err := uh.Repo.Update(&u); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	id := ctx.GetString("user_id")
	var uur UpdateUserRequest
	if err := ctx.ShouldBindJSON(&uur); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid Request body"})
		return
	}

	var u User
	if err := uh.Repo.FindById(id, &u); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	if exists, err := uh.Repo.UsernameAlreadyExists(uur.Username); exists {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Username already exists. Try another."})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	if exists, err := uh.Repo.EmailAlreadyExists(uur.Email); exists {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Email already exists. Try another."})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	if u.Username != "" {
		u.Username = uur.Username
	}
	if u.Email != "" {
		u.Email = uur.Email
	}

	err := uh.Repo.Update(&u)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.InternalServerError())
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
