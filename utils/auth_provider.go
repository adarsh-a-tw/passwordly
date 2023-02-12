package utils

import (
	"fmt"
	"time"

	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/golang-jwt/jwt/v4"
)

type AuthProvider interface {
	GenerateToken(uid string) (tokenStr string, err error)
	VerifyToken(tokenStr string) (uid string, err error)
}

type AuthProviderImpl struct{}

func (ap *AuthProviderImpl) GenerateToken(uid string) (tokenStr string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": uid,
		"exp":     jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	})

	return token.SignedString(common.JWTSecretKey)
}

func (ap *AuthProviderImpl) VerifyToken(tokenStr string) (uid string, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return common.JWTSecretKey, nil
	})

	if err != nil {
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uid = fmt.Sprintf("%s", claims["user_id"])
		return
	}

	err = fmt.Errorf("token is invalid")
	return
}
