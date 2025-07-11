package repository

import (
	"errors"

	"github.com/kadekchresna/go-boilerplate/internal/v1/model"
	"github.com/kadekchresna/go-boilerplate/internal/v1/repository/dao"
	repository_interface "github.com/kadekchresna/go-boilerplate/internal/v1/repository/interface"
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

func (r *usersRepository) Get() (*model.Users, error) {

	u := dao.UsersDAO{}
	if err := r.db.Model(dao.UsersDAO{}).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	res := u.ToModel()

	return &res, nil
}
