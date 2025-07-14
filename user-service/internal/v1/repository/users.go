package repository

import (
	"context"

	"github.com/kadekchresna/ecommerce/user-service/internal/v1/helper/dto"
	"github.com/kadekchresna/ecommerce/user-service/internal/v1/model"
	"github.com/kadekchresna/ecommerce/user-service/internal/v1/repository/dao"
	repository_interface "github.com/kadekchresna/ecommerce/user-service/internal/v1/repository/interface"
	"gorm.io/gorm"
)

type usersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) repository_interface.IUsersRepository {
	return &usersRepository{
		db: db,
	}
}

func (r *usersRepository) GetByEmailOrPhoneNumber(ctx context.Context, request *dto.LoginUserRequest) (*model.Users, error) {
	var user dao.UsersAuthDAO
	err := r.db.WithContext(ctx).Preload("User").Where("email = ? OR phone_number = ?", request.LoginCredential, request.LoginCredential).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	ud := model.Users{
		UUID:     user.User.UUID,
		FullName: user.User.Fullname,
		Code:     user.User.Code,
		UsersAuth: model.UsersAuth{
			Email:    user.Email,
			Password: user.Password,
			Salt:     user.Salt,
		},
	}

	return &ud, nil
}
