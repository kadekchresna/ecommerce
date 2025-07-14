package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/product-service/helper/logger"
	"github.com/kadekchresna/ecommerce/product-service/internal/v1/dto"
	"github.com/kadekchresna/ecommerce/product-service/internal/v1/model"
	"github.com/kadekchresna/ecommerce/product-service/internal/v1/repository/dao"
	repository_interface "github.com/kadekchresna/ecommerce/product-service/internal/v1/repository/interface"
	"gorm.io/gorm"
)

type productsRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository_interface.IProductsRepository {
	return &productsRepository{
		db: db,
	}
}

func (r *productsRepository) GetProductsPaginate(ctx context.Context, request *dto.GetProductsPaginateRequest) (*dto.GetProductsPaginateResponse, error) {

	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit <= 0 {
		request.Limit = 10
	}

	offset := (request.Page - 1) * request.Limit
	var total int64 = 0

	countQ := r.db.WithContext(ctx).
		Model(&dao.ProductsDAO{})

	if len(request.Search) > 0 {
		countQ = countQ.Where("LOWER(title) LIKE ?", fmt.Sprintf("%%%s%%", request.Search))
	}

	if err := countQ.
		Count(&total).Error; err != nil {
		err = fmt.Errorf("error Query Select :: productsRepository.GetProductsPaginate().Count() %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return nil, err
	}

	totalPages := int((total + int64(request.Limit) - 1) / int64(request.Limit))
	if total == 0 {
		return &dto.GetProductsPaginateResponse{
			Products:   nil,
			Page:       request.Page,
			Limit:      request.Limit,
			Total:      0,
			TotalPages: totalPages,
		}, nil
	}

	products := []dao.ProductsDAO{}

	db := r.db.WithContext(ctx).Model(dao.ProductsDAO{}).
		Order("created_at DESC").
		Offset(offset).
		Limit(request.Limit)

	if len(request.Search) > 0 {
		db = db.Where("LOWER(title) LIKE ?", fmt.Sprintf("%%%s%%", request.Search))
	}

	if err := db.Find(&products).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.GetProductsPaginateResponse{
				Products:   nil,
				Page:       request.Page,
				Limit:      request.Limit,
				Total:      0,
				TotalPages: totalPages,
			}, nil
		}

		err = fmt.Errorf("error Query Select :: productsRepository.GetProductsPaginate().Find() %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return nil, err
	}

	productsWithStock := make([]dto.ProductWithStock, 0, len(products))

	for _, p := range products {
		productsWithStock = append(productsWithStock, dto.ProductWithStock{
			Products: model.Products{
				UUID:        p.UUID,
				Title:       p.Title,
				Desc:        p.Desc,
				TopImageURL: p.TopImageURL,
				Price:       p.Price,
				Code:        p.Code,
				CreatedAt:   p.CreatedAt,
				UpdatedAt:   p.UpdatedAt,
				CreatedBy:   p.CreatedBy,
				UpdatedBy:   p.UpdatedBy,
			},
		})
	}

	return &dto.GetProductsPaginateResponse{
		Products:   productsWithStock,
		Page:       request.Page,
		Limit:      request.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (r *productsRepository) GetProduct(ctx context.Context, productUUID uuid.UUID) (*model.Products, error) {
	p := dao.ProductsDAO{}
	if err := r.db.WithContext(ctx).Where("uuid = ?", productUUID).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("error Query Select :: productsRepository.GetProduct(). product is not found", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())

			return nil, nil
		}

		err = fmt.Errorf("error Query Select :: productsRepository.GetProduct(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return nil, err
	}

	return &model.Products{
		UUID:  p.UUID,
		Title: p.Title,
		Desc:  p.Desc,
		Price: p.Price,
	}, nil
}
