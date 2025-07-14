package usecase_interface

import (
	"context"

	"github.com/kadekchresna/ecommerce/user-service/internal/v1/helper/dto"
)

type IAuthUsecase interface {
	Login(ctx context.Context, request *dto.LoginUserRequest) (*dto.LoginUserResponse, error)
}
