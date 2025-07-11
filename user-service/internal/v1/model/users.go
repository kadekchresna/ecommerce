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
}
