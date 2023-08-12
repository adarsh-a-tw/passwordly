package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/golang-jwt/jwt/v4"
)

type AuthProvider interface {
	GenerateTokenPair(uid string) (tokenPair AuthTokenPair, err error)
	GenerateAccessToken(refreshToken string) (accessToken string, err error)
	VerifyAccessToken(accessToken string) (uid string, err error)
}

type AuthTokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AuthTokenType string

const (
	AccessToken  AuthTokenType = "ACCESS"
	RefreshToken AuthTokenType = "REFRESH"
)

type parsedAuthToken struct {
	uid       string
	tokenType AuthTokenType
}

type AuthProviderImpl struct{}

func (ap *AuthProviderImpl) GenerateTokenPair(uid string) (tokenPair AuthTokenPair, err error) {
	authTokenPair := AuthTokenPair{}

	if accessToken, err := generateJwtTokenString(uid, AccessToken, time.Now().Add(10*time.Minute)); err != nil {
		return authTokenPair, err
	} else {
		authTokenPair.AccessToken = accessToken
	}

	if refreshToken, err := generateJwtTokenString(uid, RefreshToken, time.Now().Add(24*time.Hour)); err != nil {
		return authTokenPair, err
	} else {
		authTokenPair.RefreshToken = refreshToken
	}

	return authTokenPair, nil
}

func (ap *AuthProviderImpl) GenerateAccessToken(refreshToken string) (accessToken string, err error) {
	var pat parsedAuthToken
	if pat, err = parseJwtTokenString(refreshToken); err != nil {
		return
	}

	if pat.tokenType == AccessToken {
		err = errors.New("Invalid AuthToken Type")
		return
	}

	accessToken, err = generateJwtTokenString(pat.uid, AccessToken, time.Now().Add(10*time.Minute))
	return
}

func (ap *AuthProviderImpl) VerifyAccessToken(accessToken string) (uid string, err error) {
	var pat parsedAuthToken
	if pat, err = parseJwtTokenString(accessToken); err != nil {
		return
	}

	if pat.tokenType == RefreshToken {
		err = errors.New("Invalid AuthToken Type")
		return
	}

	uid = pat.uid
	return
}

func generateJwtTokenString(uid string, tokenType AuthTokenType, ttl time.Time) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": uid,
		"type":    tokenType,
		"exp":     jwt.NewNumericDate(ttl),
	}).SignedString([]byte(common.Cfg.JwtSecretKey))
}

func parseJwtTokenString(tokenStr string) (pat parsedAuthToken, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(common.Cfg.JwtSecretKey), nil
	})

	if err != nil {
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		pat.uid = fmt.Sprintf("%s", claims["user_id"])
		pat.tokenType, err = getAuthType(fmt.Sprintf("%s", claims["type"]))
		return
	}

	err = errors.New("Invalid AuthToken")
	return
}

func getAuthType(authTypeString string) (authType AuthTokenType, err error) {
	switch authTypeString {
	case string(AccessToken):
		authType = AccessToken
	case string(RefreshToken):
		authType = RefreshToken
	default:
		err = errors.New("Invalid AuthToken Type")
	}
	return
}
