package dao

import (
	"github.com/google/uuid"
)

type Shop struct {
	UUID     uuid.UUID `json:"uuid"`
	Code     string    `json:"code"`
	UserUUID uuid.UUID `json:"user_uuid"`
	Name     string    `json:"name"`
	Desc     string    `json:"desc"`
}

type ShopResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    Shop   `json:"data"`
}
