package model

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserToken struct {
	Id   int64
	Name string
	jwt.RegisteredClaims
}

// 签名密钥
const signKey = "ydsy"

func GetJwt(id int64, name string) (string, error) {
	if id < 0 || name == "" {
		return "", errors.New("参数错误")
	}
	token := &UserToken{
		Id:   id,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ydsy",
			Subject:   "ydsy",
			Audience:  jwt.ClaimStrings{"Android", "IOS", "H5"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "Test-1",
		},
	}
	tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, token).SignedString([]byte(signKey))
	return tokenStr, err
}

func CheckJwt(tokenStr string) (*UserToken, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("校验失败,Token不合格")
	}

	claims, ok := token.Claims.(*UserToken)
	if !ok {
		return nil, errors.New("token转义失败")
	}
	return claims, nil
}

// 退出登录用到的jwt
func BlacklistToken(context *gin.Context) {
	tokenStr := ""
	context.Set("token", tokenStr)
}
