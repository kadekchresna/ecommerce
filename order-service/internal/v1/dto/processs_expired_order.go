package dto

import (
	"time"

	"github.com/kadekchresna/ecommerce/order-service/internal/v1/model"
)

type ProcessExpiredOrderRequest struct {
	Status    model.OrderStatusType
	Statuses  []model.OrderStatusType
	Type      string
	OlderThan time.Time
	Limit     int
	Offset    int
}
