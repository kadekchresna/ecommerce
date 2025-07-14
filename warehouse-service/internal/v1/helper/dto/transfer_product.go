package dto

import "github.com/google/uuid"

type TransferProductRequest struct {
	ProductUUID         uuid.UUID `json:"product_uuid"`
	TargetWarehouseUUID uuid.UUID `json:"target_warehouse_uuid"`
}
