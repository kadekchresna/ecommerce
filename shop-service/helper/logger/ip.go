package logger

import (
	"context"

	"github.com/labstack/echo/v4"
)

const ClientIPKey key = "client_ip"

func ClientIPMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()

		ctx := context.WithValue(c.Request().Context(), ClientIPKey, ip)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
