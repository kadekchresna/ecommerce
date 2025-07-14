package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/product-service/helper/logger"
	"github.com/kadekchresna/ecommerce/product-service/internal/v1/dto"
	"github.com/kadekchresna/ecommerce/product-service/internal/v1/model"
	repository_interface "github.com/kadekchresna/ecommerce/product-service/internal/v1/repository/interface"
	usecase_interface "github.com/kadekchresna/ecommerce/product-service/internal/v1/usecase/interface"
)

type productsUsecase struct {
	ProductsRepository  repository_interface.IProductsRepository
	ShopRepository      repository_interface.IShopRepository
	WarehouseRepository repository_interface.IWarehouseRepository
}

func NewProductsUsecase(
	ProductsRepository repository_interface.IProductsRepository,
	ShopRepository repository_interface.IShopRepository,
	WarehouseRepository repository_interface.IWarehouseRepository,
) usecase_interface.IProductsUsecase {
	return &productsUsecase{
		ProductsRepository:  ProductsRepository,
		ShopRepository:      ShopRepository,
		WarehouseRepository: WarehouseRepository,
	}
}

func (u *productsUsecase) GetProductsPaginate(ctx context.Context, request *dto.GetProductsPaginateRequest) (*dto.GetProductsPaginateResponse, error) {

	res, err := u.ProductsRepository.GetProductsPaginate(ctx, request)
	if err != nil {
		return nil, err
	}

	for i := range res.Products {
		productStock, err := u.WarehouseRepository.GetProductStock(ctx, res.Products[i].UUID)
		if err != nil {
			err := fmt.Errorf("error get product stock :: WarehouseRepository.GetProductStock(). %s", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())
			return nil, err
		}

		shop, err := u.ShopRepository.GetShop(ctx, productStock.ShopUUID)
		if err != nil {
			err := fmt.Errorf("error get shop :: ShopRepository.GetShop(). %s", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())
			return nil, err
		}

		res.Products[i].AvailableStock = productStock.StartQuantity - productStock.ReserveQuantity
		res.Products[i].WarehouseStatus = productStock.Status
		res.Products[i].WarehouseName = productStock.WarehouseName
		res.Products[i].ShopName = shop.Name
	}

	return res, nil
}

func (u *productsUsecase) GetProduct(ctx context.Context, productUUID uuid.UUID) (*model.Products, error) {
	return u.ProductsRepository.GetProduct(ctx, productUUID)
}
