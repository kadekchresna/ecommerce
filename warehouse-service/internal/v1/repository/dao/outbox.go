package dao

import (
	"time"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/model"
)

type OutboxDAO struct {
	UUID       uuid.UUID              `gorm:"type:uuid;primaryKey"`
	Metadata   string                 `gorm:"type:jsonb;not null"`
	Response   string                 `gorm:"type:jsonb;not null"`
	Status     model.OutboxStatusType `gorm:"type:outbox_status_type;not null"`
	Action     string                 `gorm:"type:varchar;not null;default:''"`
	Type       string                 `gorm:"type:varchar;not null;default:''"`
	Reference  string                 `gorm:"type:varchar;not null;default:''"`
	RetryCount int                    `gorm:"column:retry_count;not null"`
	CreatedAt  time.Time              `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt  time.Time              `gorm:"type:timestamptz;not null;default:now()"`
}

func (OutboxDAO) TableName() string {
	return "outbox"
}
