package model

import (
	"time"

	"github.com/google/uuid"
)

type WarehouseStatusType string

const (
	WarehouseStatusActive   WarehouseStatusType = "active"
	WarehouseStatusInactive WarehouseStatusType = "inactive"
)

type Users struct {
	UUID      uuid.UUID `json:"uuid"`
	FullName  string    `json:"fullname"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`
}

type Warehouse struct {
	UUID     uuid.UUID           `json:"uuid"`
	Name     string              `json:"name"`
	Code     string              `json:"code"`
	Desc     string              `json:"desc"`
	ShopUUID uuid.UUID           `json:"shop_uuid"`
	Status   WarehouseStatusType `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`

	WarehouseStock WarehouseStock `json:"warehouse"`
}

type WarehouseStock struct {
	UUID          uuid.UUID `json:"uuid"`
	WarehouseUUID uuid.UUID `json:"warehouse_uuid"`
	ProductUUID   uuid.UUID `json:"product_uuid"`

	StartQuantity   int `json:"start_quantity"`
	ReserveQuantity int `json:"reserve_quantity"`

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

type Order struct {
	UUID uuid.UUID `json:"uuid"`
}
