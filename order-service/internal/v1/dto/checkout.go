package dto

import (
	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/model"
)

type Order struct {
	OrderDetails []OrderDetails `json:"order_details"`
}

type OrderDetails struct {
	ProductUUID uuid.UUID `json:"product_uuid"`
	Quantity    int       `json:"quantity"`
}

type CheckoutRequest struct {
	Order
	UserUUID uuid.UUID `json:"-"`
}

type CreateCheckoutRequest struct {
	Order        model.Order
	OrderDetails []model.OrderDetail
	EventType    string
	EventPayload model.OutboxOrderCreatedMetaRequest
	UserUUID     uuid.UUID
}

type CheckoutResponse struct {
}
