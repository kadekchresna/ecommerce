package model

import "github.com/google/uuid"

type ProductStock struct {
	ProductUUID     uuid.UUID `json:"product_uuid"`
	WarehouseUUID   uuid.UUID `json:"warehouse_uuid"`
	WarehouseName   string    `json:"warehouse_name"`
	ShopUUID        uuid.UUID `json:"warehouse_shop_uuid"`
	Status          string    `json:"status"`
	ReserveQuantity int       `json:"reserve_quantity"`
	StartQuantity   int       `json:"start_quantity"`
}
