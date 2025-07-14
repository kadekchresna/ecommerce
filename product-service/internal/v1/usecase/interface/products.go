package usecase_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/product-service/internal/v1/dto"
	"github.com/kadekchresna/ecommerce/product-service/internal/v1/model"
)

type IProductsUsecase interface {
	GetProductsPaginate(ctx context.Context, request *dto.GetProductsPaginateRequest) (*dto.GetProductsPaginateResponse, error)
	GetProduct(ctx context.Context, productUUID uuid.UUID) (*model.Products, error)
}
