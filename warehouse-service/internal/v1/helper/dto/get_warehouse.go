package dto

import (
	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/model"
)

type GetWarehouseRequest struct {
	UUID uuid.UUID `json:"uuid"`
}

type GetWarehouseResponse struct {
	Warehouse *model.Warehouse `json:"warehouse"`
}
