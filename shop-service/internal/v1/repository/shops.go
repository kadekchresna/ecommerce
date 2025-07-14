package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/shop-service/internal/v1/model"
	"github.com/kadekchresna/ecommerce/shop-service/internal/v1/repository/dao"
	repository_interface "github.com/kadekchresna/ecommerce/shop-service/internal/v1/repository/interface"
	"gorm.io/gorm"
)

type shopsRepository struct {
	db *gorm.DB
}

func NewShopsRepository(db *gorm.DB) repository_interface.IShopsRepository {
	return &shopsRepository{
		db: db,
	}
}

func (r *shopsRepository) GetShops(ctx context.Context, shopUUID uuid.UUID) (*model.Shops, error) {

	s := dao.ShopsDAO{}
	if err := r.db.WithContext(ctx).Model(dao.ShopsDAO{}).First(&s, shopUUID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &model.Shops{
		UUID: s.UUID,
		Code: s.Code,
		Name: s.Name,
		Desc: s.Desc,
	}, nil
}
