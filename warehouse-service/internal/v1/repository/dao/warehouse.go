package dao

import (
	"time"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/model"
)

type WarehouseDAO struct {
	UUID     uuid.UUID                 `gorm:"type:uuid;primaryKey"`
	Name     string                    `gorm:"type:varchar;not null"`
	Code     string                    `gorm:"type:varchar;not null"`
	Desc     string                    `gorm:"type:text;not null"`
	ShopUUID uuid.UUID                 `gorm:"type:uuid;not null"`
	Status   model.WarehouseStatusType `gorm:"type:warehouse_status_type;not null"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null"`
	UpdatedBy uuid.UUID `gorm:"type:uuid;not null"`
}

func (WarehouseDAO) TableName() string {
	return "warehouses"
}

type WarehouseStockDAO struct {
	UUID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	WarehouseUUID uuid.UUID `gorm:"type:uuid;not null"`
	ProductUUID   uuid.UUID `gorm:"type:uuid;not null"`

	StartQuantity   int `gorm:"type:int;not null;default:0"`
	ReserveQuantity int `gorm:"type:int;not null;default:0"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null"`
	UpdatedBy uuid.UUID `gorm:"type:uuid;not null"`

	Warehouse WarehouseDAO `gorm:"foreignKey:UUID;references:WarehouseUUID"`
}

func (WarehouseStockDAO) TableName() string {
	return "warehouses_stock"
}
