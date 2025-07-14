package model

import (
	"time"

	"github.com/google/uuid"
)

type Shops struct {
	UUID     uuid.UUID `json:"uuid"`
	Code     string    `json:"code"`
	UserUUID uuid.UUID `json:"user_uuid"`
	Name     string    `json:"name"`
	Desc     string    `json:"desc"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`
}
