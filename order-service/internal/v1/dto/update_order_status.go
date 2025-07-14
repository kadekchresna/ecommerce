package dto

import (
	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/model"
)

type UpdateOrderStatusRequest struct {
	OrderUUID    uuid.UUID
	NewStatus    model.OrderStatusType
	Metadata     model.OrderMetadata
	EventType    string
	EventPayload any
	InboxUUID    uuid.UUID
}
