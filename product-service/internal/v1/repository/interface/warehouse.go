package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/product-service/internal/v1/model"
)

type IWarehouseRepository interface {
	GetProductStock(ctx context.Context, productUUID uuid.UUID) (*model.ProductStock, error)
}
