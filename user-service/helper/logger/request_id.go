package logger

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type key string

const RequestIDKey key = "request_id"

// In your middleware
func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Request().Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx := context.WithValue(c.Request().Context(), RequestIDKey, requestID)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
