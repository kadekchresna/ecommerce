package model

import (
	"time"

	"github.com/google/uuid"
)

type Products struct {
	UUID        uuid.UUID `json:"uuid"`
	Title       string    `json:"title"`
	Desc        string    `json:"desc"`
	TopImageURL string    `json:"top_image_url"`
	Price       float64   `json:"price"`
	Code        string    `json:"code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   uuid.UUID `json:"created_by"`
	UpdatedBy   uuid.UUID `json:"updated_by"`
}
