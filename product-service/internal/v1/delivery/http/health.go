package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewHealthHandler(
	v1 *echo.Group,
	DB *gorm.DB,
) *healthHandler {
	h := &healthHandler{
		DB: DB,
	}

	v1.GET("/healthz", h.Healthz)

	return h

}

type healthHandler struct {
	v1 *echo.Group
	DB *gorm.DB
}

type HealthStatus struct {
	Postgres string `json:"postgres"`
	Redis    string `json:"redis"`
	Kafka    string `json:"kafka"`
}

func (h *healthHandler) Healthz(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
	defer cancel()

	status := HealthStatus{
		Postgres: "ok",
	}

	if err := PingGormDB(ctx, h.DB); err != nil {
		status.Postgres = "error: " + err.Error()
	}

	code := http.StatusOK
	if status.Postgres != "ok" {
		code = http.StatusServiceUnavailable
	}

	return c.JSON(code, status)
}

func PingGormDB(ctx context.Context, db *gorm.DB) error {

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctxTimeout); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	return nil
}
