package dao

import (
	"time"

	"github.com/google/uuid"
	"github.com/kadekchresna/go-boilerplate/internal/v1/model"
)

type UsersDAO struct {
	UUID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	FullName  string    `gorm:"type:varchar;not null;default:''"`
	Code      string    `gorm:"type:varchar;not null;default:'';"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null"`
	UpdatedBy uuid.UUID `gorm:"type:uuid;not null"`
}

func (UsersDAO) TableName() string {
	return "users"
}

func (u UsersDAO) ToModel() model.Users {
	return model.Users{
		UUID:      u.UUID,
		FullName:  u.FullName,
		Code:      u.Code,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		CreatedBy: u.CreatedBy,
		UpdatedBy: u.UpdatedBy,
	}
}

type UsersAuthDAO struct {
	UUID        uuid.UUID `gorm:"type:uuid;default:primaryKey"`
	Email       string    `gorm:"type:varchar;not null;default:''"`
	PhoneNumber string    `gorm:"type:varchar;not null;default:''"`
	Password    string    `gorm:"type:varchar;not null;default:''"`
	Salt        string    `gorm:"type:varchar;not null;default:''"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UserUUID    uuid.UUID `gorm:"type:uuid;not null;index:user_auth_user_uuid_idx"`
	CreatedBy   uuid.UUID `gorm:"type:uuid;not null"`
	UpdatedBy   uuid.UUID `gorm:"type:uuid;not null"`
}

func (UsersAuthDAO) TableName() string {
	return "users_auth"
}
