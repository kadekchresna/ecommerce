package dao

import (
	"time"

	"github.com/google/uuid"
)

type ProductsDAO struct {
	UUID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title       string    `gorm:"type:varchar;not null;default:''"`
	Desc        string    `gorm:"type:text;not null;default:''"`
	TopImageURL string    `gorm:"type:varchar;not null;default:''"`
	Price       float64   `gorm:"type:float8;not null;default:0.0"`
	Code        string    `gorm:"type:varchar;not null;default:'';index:products_code_idx"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null"`
	UpdatedBy uuid.UUID `gorm:"type:uuid;not null"`
}

func (ProductsDAO) TableName() string {
	return "products"
}
