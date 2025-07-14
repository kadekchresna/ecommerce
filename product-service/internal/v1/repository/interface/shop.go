package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/product-service/internal/v1/model"
)

type IShopRepository interface {
	GetShop(ctx context.Context, shopUUID uuid.UUID) (*model.Shop, error)
}
