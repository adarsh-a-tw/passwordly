package middleware

import (
	"net/http"
	"strings"

	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/adarsh-a-tw/passwordly/utils"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware(ctx *gin.Context) {
	tokenStr := extractTokenFromHeader(ctx.Request.Header.Get("authorization"))
	ap := utils.AuthProviderImpl{}
	if uid, err := ap.VerifyAccessToken(tokenStr); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.ErrorResponse{Message: "Invalid Token"})
	} else {
		ctx.Set("user_id", uid)
		ctx.Next()
	}
}

func extractTokenFromHeader(headerString string) string {
	splitToken := strings.Split(headerString, "Bearer ")
	return strings.Join(splitToken, "")
}
