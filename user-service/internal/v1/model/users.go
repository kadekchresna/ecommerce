package model

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	UUID      uuid.UUID `json:"uuid"`
	FullName  string    `json:"fullname"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`

	UsersAuth UsersAuth `json:"user_auth"`
}

type UsersAuth struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	Salt        string `json:"salt"`
}
