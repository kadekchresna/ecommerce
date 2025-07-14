package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/model"
)

type IProductRepository interface {
	GetProduct(ctx context.Context, productUUID uuid.UUID) (*model.Products, error)
}
