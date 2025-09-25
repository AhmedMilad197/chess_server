package utils

import (
	"chess_server/config"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func CreateToken(data map[string]interface{}, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{}
	for k, v := range data {
		claims[k] = v
	}
	claims["exp"] = time.Now().Add(duration).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.Config.JWTSecret)
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return config.Config.JWTSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
