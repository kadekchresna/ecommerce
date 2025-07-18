package dao

import (
	"time"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/model"
)

type InboxDAO struct {
	UUID       uuid.UUID             `gorm:"type:uuid;primaryKey"`
	Metadata   string                `gorm:"type:jsonb;not null;default:'{}'"`
	Response   string                `gorm:"type:jsonb;not null"`
	Status     model.InboxStatusType `gorm:"type:inbox_status_type;not null"`
	Type       string                `gorm:"type:varchar;not null;default:''"`
	Action     string                `gorm:"type:varchar;not null;default:''"`
	Reference  string                `gorm:"type:varchar;not null;default:''"`
	RetryCount int                   `gorm:"column:retry_count;not null"`
	CreatedAt  time.Time             `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt  time.Time             `gorm:"type:timestamptz;not null;default:now()"`
}

func (InboxDAO) TableName() string {
	return "inbox"
}
