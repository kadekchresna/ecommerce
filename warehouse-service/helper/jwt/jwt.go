package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserUUID     uuid.UUID `json:"user_id"`
	UserFullName string    `json:"user_fullname"`
	jwt.RegisteredClaims
}

func GenerateToken(jwtSecret string, userUUID uuid.UUID, userFullname string) (string, string, error) {
	accessClaims := Claims{
		UserUUID:     userUUID,
		UserFullName: userFullname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := at.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := Claims{
		UserUUID:     userUUID,
		UserFullName: userFullname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := rt.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ParseAccessToken(jwtSecret string, tokenStr string) (*Claims, error) {
	return parseToken(jwtSecret, tokenStr)
}

func ParseRefreshToken(jwtSecret string, tokenStr string) (*Claims, error) {
	return parseToken(jwtSecret, tokenStr)
}

func parseToken(jwtSecret string, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
