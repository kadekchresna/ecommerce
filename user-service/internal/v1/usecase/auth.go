package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kadekchresna/ecommerce/user-service/config"
	"github.com/kadekchresna/ecommerce/user-service/helper/jwt"
	"github.com/kadekchresna/ecommerce/user-service/helper/password"
	"github.com/kadekchresna/ecommerce/user-service/internal/v1/helper/dto"
	repository_interface "github.com/kadekchresna/ecommerce/user-service/internal/v1/repository/interface"
	usecase_interface "github.com/kadekchresna/ecommerce/user-service/internal/v1/usecase/interface"
)

type authUsecase struct {
	config          config.Config
	UsersRepository repository_interface.IUsersRepository
}

func NewAuthUsecase(
	config config.Config,
	UsersRepository repository_interface.IUsersRepository,
) usecase_interface.IAuthUsecase {
	return &authUsecase{
		config:          config,
		UsersRepository: UsersRepository,
	}
}

func (u *authUsecase) Login(ctx context.Context, request *dto.LoginUserRequest) (*dto.LoginUserResponse, error) {

	request.LoginCredential = strings.TrimSpace(request.LoginCredential)
	if len(request.LoginCredential) == 0 {
		return nil, errors.New("Email or Phone Number is required")
	}

	request.Password = strings.TrimSpace(request.Password)
	if len(request.Password) == 0 {
		return nil, errors.New("Password is required")
	}

	user, err := u.UsersRepository.GetByEmailOrPhoneNumber(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("Login failed. %s", err.Error())
	}

	if user == nil {
		return nil, fmt.Errorf("User with credential %s is not found", request.LoginCredential)
	}

	if !password.ComparePasswordWithHash(request.Password, user.UsersAuth.Salt, user.UsersAuth.Password) {
		return nil, errors.New("Password is invalid")
	}

	accessToken, refreshToken, err := jwt.GenerateToken(u.config.AppJWTSecret, user.UUID, user.FullName)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate token user, %s", err.Error())
	}

	return &dto.LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
