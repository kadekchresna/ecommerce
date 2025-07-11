package usecase_interface

import "github.com/kadekchresna/go-boilerplate/internal/v1/model"

type IUsersUsecase interface {
	Get() (*model.Users, error)
}
