package jwt

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	USER_UUID_KEY     = "user_uuid"
	USER_FULLNAME_KEY = "user_fullname"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		appTokenHeader := c.Request().Header.Get("X-App-Token")
		if authHeader == "" && appTokenHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
		}

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token format"})
			}

			claims, err := ParseAccessToken(os.Getenv("APP_JWT_SECRET"), parts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}

			c.Set(USER_UUID_KEY, claims.UserUUID)
			c.Set(USER_FULLNAME_KEY, claims.UserFullName)
		}

		if appTokenHeader != "" {
			parts := strings.Split(appTokenHeader, " ")

			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token format"})
			}

			if os.Getenv("APP_STATIC_TOKEN") != parts[1] {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}

		}

		return next(c)
	}
}
