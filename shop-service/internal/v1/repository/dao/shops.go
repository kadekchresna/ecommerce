package dao

import (
	"time"

	"github.com/google/uuid"
)

type ShopsDAO struct {
	UUID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code     string    `gorm:"type:varchar;not null;default:'';uniqueIndex:shops_code_idx"`
	UserUUID uuid.UUID `gorm:"type:uuid;not null"`
	Name     string    `gorm:"type:varchar;not null;default:''"`
	Desc     string    `gorm:"type:text;not null;default:''"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null"`
	UpdatedBy uuid.UUID `gorm:"type:uuid;not null"`
}

func (ShopsDAO) TableName() string {
	return "shops"
}
