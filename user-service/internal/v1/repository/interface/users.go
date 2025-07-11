package repository_interface

import "github.com/kadekchresna/go-boilerplate/internal/v1/model"

type IUsersRepository interface {
	Get() (*model.Users, error)
}
