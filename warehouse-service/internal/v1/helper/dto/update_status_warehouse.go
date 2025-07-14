package dto

import (
	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/model"
)

type UpdateStatusWarehouseRequest struct {
	WarehouseUUID uuid.UUID                 `json:"warehouse_uuid"`
	Status        model.WarehouseStatusType `json:"status"`
	UserUUID      uuid.UUID                 `json:"-"`
}
