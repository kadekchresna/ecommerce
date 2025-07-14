package repository_interface

import (
	"context"

	"github.com/kadekchresna/ecommerce/user-service/internal/v1/helper/dto"
	"github.com/kadekchresna/ecommerce/user-service/internal/v1/model"
)

type IUsersRepository interface {
	GetByEmailOrPhoneNumber(ctx context.Context, request *dto.LoginUserRequest) (*model.Users, error)
}
