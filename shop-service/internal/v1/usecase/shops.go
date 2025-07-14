package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/shop-service/internal/v1/model"
	repository_interface "github.com/kadekchresna/ecommerce/shop-service/internal/v1/repository/interface"
	usecase_interface "github.com/kadekchresna/ecommerce/shop-service/internal/v1/usecase/interface"
)

type shopsUsecase struct {
	ShopsRepository repository_interface.IShopsRepository
}

func NewShopsUsecase(
	ShopsRepository repository_interface.IShopsRepository,
) usecase_interface.IShopsUsecase {
	return &shopsUsecase{
		ShopsRepository: ShopsRepository,
	}
}

func (u *shopsUsecase) GetShops(ctx context.Context, shopUUID uuid.UUID) (*model.Shops, error) {
	return u.ShopsRepository.GetShops(ctx, shopUUID)
}
