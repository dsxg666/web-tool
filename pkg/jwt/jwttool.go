package jwt

import (
	"errors"
	"time"

	"github.com/dsxg666/web-tool/global"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewClaims(userId, username string) CustomClaims {
	return CustomClaims{
		UserId:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * global.JwtTokenSetting.ExpirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "tool.lankaiyun.com",
		},
	}
}

func NewJwtToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(global.JwtTokenSetting.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseJwtToken(jwtToken string) (*CustomClaims, bool, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.JwtTokenSetting.SecretKey), nil
	})

	if err != nil {
		return nil, false, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, true, nil
	} else {
		return nil, false, errors.New("invalid token")
	}
}
