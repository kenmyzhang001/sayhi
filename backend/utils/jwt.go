package utils

import (
	"errors"
	"sayhi/backend/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// getJWTSecret 获取JWT密钥
func getJWTSecret() []byte {
	if config.AppConfig != nil {
		return []byte(config.AppConfig.JWT.Secret)
	}
	return []byte("sayhi-secret-key-change-in-production")
}

// Claims JWT声明
type Claims struct {
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(username string, userID int64) (string, error) {
	nowTime := time.Now()
	expireHours := 24
	if config.AppConfig != nil {
		expireHours = config.AppConfig.JWT.ExpireTime
	}
	expireTime := nowTime.Add(time.Duration(expireHours) * time.Hour)

	claims := Claims{
		Username: username,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			NotBefore: jwt.NewNumericDate(nowTime),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(getJWTSecret())
	return token, err
}

// ParseToken 解析JWT token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

