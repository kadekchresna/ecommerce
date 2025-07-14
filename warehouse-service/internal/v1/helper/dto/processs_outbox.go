package dto

import (
	"time"

	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/model"
)

type ProcessOutboxRequest struct {
	Status    model.OutboxStatusType
	Statuses  []model.OutboxStatusType
	Type      string
	OlderThan time.Time
	Limit     int
	Offset    int
}
