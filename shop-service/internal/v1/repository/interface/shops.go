package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/shop-service/internal/v1/model"
)

type IShopsRepository interface {
	GetShops(ctx context.Context, shopUUID uuid.UUID) (*model.Shops, error)
}
