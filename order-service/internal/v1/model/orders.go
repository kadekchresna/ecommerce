package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatusType string

const (
	OrderStatusCreated        OrderStatusType = "created"
	OrderStatusInPayment      OrderStatusType = "in-payment"
	OrderStatusCancelled      OrderStatusType = "cancelled"
	OrderStatusCompleted      OrderStatusType = "completed"
	OrderStatusExpired        OrderStatusType = "expired"
	OrderStatusFailed         OrderStatusType = "failed"
	OrderStatusReservingStock OrderStatusType = "reserving-stock"
)

type OrderMetadata struct {
	Reason string `json:"reason"`
}

type Order struct {
	UUID        uuid.UUID       `json:"uuid"`
	Code        string          `json:"code"`
	Metadata    string          `json:"metadata"`
	UserUUID    uuid.UUID       `json:"user_uuid"`
	TotalAmount float64         `json:"total_amount"`
	ExpiredAt   time.Time       `json:"expired_at"`
	Status      OrderStatusType `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`
}

type OrderDetail struct {
	UUID         uuid.UUID `json:"uuid"`
	ProductUUID  uuid.UUID `json:"product_uuid"`
	ProductTitle string    `json:"product_title"`
	ProductPrice float64   `json:"product_price"`
	Quantity     int       `json:"quantity"`
	SubTotal     float64   `json:"sub_total"`
	OrderUUID    uuid.UUID `json:"order_uuid"`
}
