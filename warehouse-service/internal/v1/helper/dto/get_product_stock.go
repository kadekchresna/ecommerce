package dto

import (
	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/model"
)

type GetProductStockRequest struct {
	ProductUUID   uuid.UUID `json:"product_uuid"`
	WarehouseUUID uuid.UUID `json:"warehouse_uuid"`
}

type GetProductStockResponse struct {
	ProductUUID     uuid.UUID                 `json:"product_uuid"`
	WarehouseUUID   uuid.UUID                 `json:"warehouse_uuid"`
	WarehouseName   string                    `json:"warehouse_name"`
	ShopUUID        uuid.UUID                 `json:"warehouse_shop_uuid"`
	Status          model.WarehouseStatusType `json:"status"`
	ReserveQuantity int                       `json:"reserve_quantity"`
	StartQuantity   int                       `json:"start_quantity"`
}
