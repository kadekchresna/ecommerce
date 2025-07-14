package dto

import (
	"time"

	"github.com/kadekchresna/ecommerce/order-service/internal/v1/model"
)

type ProcessInboxRequest struct {
	Status    model.InboxStatusType
	Statuses  []model.InboxStatusType
	Type      string
	OlderThan time.Time
	Limit     int
	Offset    int
}
