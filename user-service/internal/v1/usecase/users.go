package usecase

import (
	"github.com/kadekchresna/go-boilerplate/internal/v1/model"
	repository_interface "github.com/kadekchresna/go-boilerplate/internal/v1/repository/interface"
	usecase_interface "github.com/kadekchresna/go-boilerplate/internal/v1/usecase/interface"
)

type usersUsecase struct {
	UsersRepository repository_interface.IUsersRepository
}

func NewUsersUsecase(
	UsersRepository repository_interface.IUsersRepository,
) usecase_interface.IUsersUsecase {
	return &usersUsecase{
		UsersRepository: UsersRepository,
	}
}

func (u *usersUsecase) Get() (*model.Users, error) {
	return u.UsersRepository.Get()
}
